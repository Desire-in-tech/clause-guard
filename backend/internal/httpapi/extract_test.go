package httpapi

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractTXTUpload(t *testing.T) {
	response := uploadTestFile(t, "agreement.txt", []byte("Payment is due within 30 days."))

	if response.Text != "Payment is due within 30 days." {
		t.Fatalf("unexpected extracted text: %q", response.Text)
	}
	if response.Filename != "agreement.txt" {
		t.Fatalf("unexpected filename: %q", response.Filename)
	}
	if response.Characters != len(response.Text) {
		t.Fatalf("unexpected character count: got %d, want %d", response.Characters, len(response.Text))
	}
}

func TestExtractDOCXUpload(t *testing.T) {
	docx := makeTestDOCX(t, `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>First paragraph.</w:t></w:r></w:p><w:p><w:r><w:t>Second paragraph.</w:t></w:r></w:p></w:body></w:document>`)
	response := uploadTestFile(t, "agreement.docx", docx)

	if response.Text != "First paragraph.\nSecond paragraph." {
		t.Fatalf("unexpected extracted text: %q", response.Text)
	}
}

func TestExtractRejectsUnsupportedType(t *testing.T) {
	body := makeMultipartBody(t, "agreement.rtf", []byte("Some text."))
	request := httptest.NewRequest(http.MethodPost, "/api/extract", body)
	request.Header.Set("Content-Type", body.contentType)
	recorder := httptest.NewRecorder()

	(&Handlers{}).extract(recorder, request)

	if recorder.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("unexpected status: got %d, want %d", recorder.Code, http.StatusUnsupportedMediaType)
	}
	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response["error"] == "" {
		t.Fatal("expected a JSON error message")
	}
}

type multipartTestBody struct {
	*bytes.Buffer
	contentType string
}

func makeMultipartBody(t *testing.T, filename string, contents []byte) multipartTestBody {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create multipart file: %v", err)
	}
	if _, err := part.Write(contents); err != nil {
		t.Fatalf("write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	return multipartTestBody{
		Buffer:      &body,
		contentType: writer.FormDataContentType(),
	}
}

func uploadTestFile(t *testing.T, filename string, contents []byte) extractResponse {
	t.Helper()
	body := makeMultipartBody(t, filename, contents)
	request := httptest.NewRequest(http.MethodPost, "/api/extract", body)
	request.Header.Set("Content-Type", body.contentType)
	recorder := httptest.NewRecorder()

	(&Handlers{}).extract(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, body %s", recorder.Code, recorder.Body.String())
	}
	var response extractResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return response
}

func makeTestDOCX(t *testing.T, documentXML string) []byte {
	t.Helper()
	var contents bytes.Buffer
	writer := zip.NewWriter(&contents)
	document, err := writer.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create DOCX document part: %v", err)
	}
	if _, err := io.WriteString(document, documentXML); err != nil {
		t.Fatalf("write DOCX document part: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close DOCX archive: %v", err)
	}
	return contents.Bytes()
}
