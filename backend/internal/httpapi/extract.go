package httpapi

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

const (
	maxUploadBytes  = 5 << 20
	maxDocxXMLBytes = 5 << 20
	maxPDFTextBytes = maxContractRunes * utf8.UTFMax
)

type extractFailure struct {
	status  int
	message string
}

func (e *extractFailure) Error() string {
	return e.message
}

type extractResponse struct {
	Text       string `json:"text"`
	Filename   string `json:"filename"`
	Characters int    `json:"characters"`
}

func (h *Handlers) extract(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

	filename, contents, err := readMultipartFile(r)
	if err != nil {
		writeExtractError(w, err)
		return
	}

	text, err := extractDocumentText(filename, contents)
	if err != nil {
		writeExtractError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, extractResponse{
		Text:       text,
		Filename:   filename,
		Characters: utf8.RuneCountInString(text),
	})
}

func readMultipartFile(r *http.Request) (string, []byte, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !strings.EqualFold(mediaType, "multipart/form-data") {
		return "", nil, &extractFailure{
			status:  http.StatusBadRequest,
			message: "Upload must use multipart/form-data with a file field.",
		}
	}

	reader, err := r.MultipartReader()
	if err != nil {
		return "", nil, &extractFailure{
			status:  http.StatusBadRequest,
			message: "Could not read the uploaded file.",
		}
	}

	var filename string
	var contents []byte
	fileFound := false

	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", nil, err
		}

		if part.FormName() != "file" {
			_, copyErr := io.Copy(io.Discard, part)
			_ = part.Close()
			if copyErr != nil {
				return "", nil, copyErr
			}
			continue
		}

		if fileFound {
			_ = part.Close()
			return "", nil, &extractFailure{
				status:  http.StatusBadRequest,
				message: "Upload one file at a time.",
			}
		}
		fileFound = true
		filename = part.FileName()
		if filename == "" {
			_ = part.Close()
			return "", nil, &extractFailure{
				status:  http.StatusBadRequest,
				message: "Choose a file to upload.",
			}
		}

		contents, err = io.ReadAll(part)
		_ = part.Close()
		if err != nil {
			return "", nil, err
		}
	}

	if !fileFound {
		return "", nil, &extractFailure{
			status:  http.StatusBadRequest,
			message: "The multipart request must include a file field.",
		}
	}

	return filename, contents, nil
}

func extractDocumentText(filename string, contents []byte) (string, error) {
	extension := strings.ToLower(filepath.Ext(filename))
	switch extension {
	case ".txt", ".docx", ".pdf":
	default:
		return "", &extractFailure{
			status:  http.StatusUnsupportedMediaType,
			message: "Only .txt, .docx, and .pdf files are supported.",
		}
	}

	if len(contents) == 0 {
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "The uploaded file is empty.",
		}
	}

	var text string
	var err error
	switch extension {
	case ".txt":
		if !utf8.Valid(contents) {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "The TXT file is not valid UTF-8 text.",
			}
		}
		text = strings.TrimPrefix(string(contents), "\uFEFF")
	case ".docx":
		if !hasZIPSignature(contents) {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "The file contents do not match the .docx extension.",
			}
		}
		text, err = extractDOCXText(contents)
	case ".pdf":
		if !bytes.HasPrefix(contents, []byte("%PDF-")) {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "The file contents do not match the .pdf extension.",
			}
		}
		text, err = extractPDFText(contents)
	}
	if err != nil {
		return "", err
	}

	text = strings.TrimSpace(text)
	if text == "" {
		if extension == ".pdf" {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "No selectable text found. Scanned PDFs are not supported yet.",
			}
		}
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "No text was found in the uploaded file.",
		}
	}
	if !utf8.ValidString(text) {
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "The uploaded file did not contain readable text.",
		}
	}
	if utf8.RuneCountInString(text) > maxContractRunes {
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "Extracted text must not exceed 100,000 characters.",
		}
	}

	return text, nil
}

func hasZIPSignature(contents []byte) bool {
	if len(contents) < 4 || contents[0] != 'P' || contents[1] != 'K' {
		return false
	}
	return (contents[2] == 3 && contents[3] == 4) ||
		(contents[2] == 5 && contents[3] == 6) ||
		(contents[2] == 7 && contents[3] == 8)
}

func extractDOCXText(contents []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(contents), int64(len(contents)))
	if err != nil {
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "Could not read the DOCX file.",
		}
	}

	for _, file := range reader.File {
		if file.Name != "word/document.xml" {
			continue
		}
		if file.UncompressedSize64 > maxDocxXMLBytes {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "The DOCX document content is too large to extract.",
			}
		}

		xmlFile, err := file.Open()
		if err != nil {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "Could not read the DOCX file.",
			}
		}
		xmlContents, readErr := io.ReadAll(io.LimitReader(xmlFile, maxDocxXMLBytes+1))
		closeErr := xmlFile.Close()
		if readErr != nil || closeErr != nil {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "Could not read the DOCX file.",
			}
		}
		if len(xmlContents) > maxDocxXMLBytes {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "The DOCX document content is too large to extract.",
			}
		}
		return extractWordDocumentText(xmlContents)
	}

	return "", &extractFailure{
		status:  http.StatusBadRequest,
		message: "The DOCX file does not contain word/document.xml.",
	}
}

func extractWordDocumentText(contents []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(contents))
	var text strings.Builder
	inParagraph := false
	inText := false
	paragraphCount := 0

	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", &extractFailure{
				status:  http.StatusBadRequest,
				message: "Could not read text from the DOCX file.",
			}
		}

		switch value := token.(type) {
		case xml.StartElement:
			switch value.Name.Local {
			case "p":
				if paragraphCount > 0 {
					text.WriteByte('\n')
				}
				paragraphCount++
				inParagraph = true
			case "t":
				inText = true
			case "tab":
				if inParagraph {
					text.WriteByte('\t')
				}
			case "br", "cr":
				if inParagraph {
					text.WriteByte('\n')
				}
			}
		case xml.EndElement:
			switch value.Name.Local {
			case "t":
				inText = false
			case "p":
				inParagraph = false
			}
		case xml.CharData:
			if inParagraph && inText {
				_, _ = text.Write(value)
			}
		}
	}

	return text.String(), nil
}

func extractPDFText(contents []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(contents), int64(len(contents)))
	if err != nil {
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "Could not read the PDF file.",
		}
	}

	textReader, err := reader.GetPlainText()
	if err != nil {
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "Could not extract text from the PDF file.",
		}
	}
	text, err := io.ReadAll(io.LimitReader(textReader, maxPDFTextBytes+1))
	if err != nil {
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "Could not extract text from the PDF file.",
		}
	}
	if len(text) > maxPDFTextBytes {
		return "", &extractFailure{
			status:  http.StatusBadRequest,
			message: "Extracted text must not exceed 100,000 characters.",
		}
	}
	return string(text), nil
}

func writeExtractError(w http.ResponseWriter, err error) {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		writeError(w, http.StatusRequestEntityTooLarge, "The upload must not exceed 5 MB.")
		return
	}

	var failure *extractFailure
	if errors.As(err, &failure) {
		writeError(w, failure.status, failure.message)
		return
	}

	writeError(w, http.StatusBadRequest, "Could not read the uploaded file.")
}
