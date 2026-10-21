package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"clauseguard/backend/internal/analysis"
	"clauseguard/backend/internal/chat"
	"clauseguard/backend/internal/contracts"
	"clauseguard/backend/internal/llm"
	"clauseguard/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

const (
	maxRequestBytes  = 1 << 20
	maxContractRunes = 100_000
)

type Handlers struct {
	analysis    *analysis.Service
	chatService *chat.Service
	reports     store.Store
	llmProvider string
	maxClauses  int
	limits      *requestLimiter
	logger      *slog.Logger
}

func (h *Handlers) health(w http.ResponseWriter, _ *http.Request) {
	storeStatus := "unavailable"
	if h.reports != nil {
		switch h.reports.(type) {
		case *store.MemoryStore:
			storeStatus = "memory"
		default:
			storeStatus = "mongodb"
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":       "ok",
		"llm_provider": h.llmProvider,
		"store":        storeStatus,
	})
}

func (h *Handlers) analyze(w http.ResponseWriter, r *http.Request) {
	var request analysis.AnalyzeRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if utf8.RuneCountInString(request.ContractText) > maxContractRunes {
		writeError(w, http.StatusBadRequest, "contract_text must not exceed 100,000 characters.")
		return
	}
	if strings.TrimSpace(request.ContractText) == "" {
		writeError(w, http.StatusBadRequest, "contract_text must not be empty.")
		return
	}
	clauses := contracts.Split(request.ContractText)
	if len(clauses) == 0 {
		writeError(w, http.StatusBadRequest, "contract_text must contain at least one clause.")
		return
	}
	if len(clauses) > h.maxClauses {
		writeError(w, http.StatusBadRequest, (&llm.TooManyClausesError{Max: h.maxClauses}).Error())
		return
	}
	if h.llmProvider == "openrouter" {
		if !h.allowAIRequest(w, r) {
			return
		}
	}

	report, err := h.analysis.Analyze(r.Context(), request.ContractText)
	if err != nil {
		if errors.Is(err, analysis.ErrEmptyContract) {
			writeError(w, http.StatusBadRequest, "contract_text must contain at least one clause.")
			return
		}
		if h.writeLLMError(w, err) {
			return
		}
		h.logger.Error("analyze contract", "error", err)
		writeError(w, http.StatusInternalServerError, "Could not analyze this contract.")
		return
	}
	if h.reports != nil {
		if err := h.reports.SaveReport(r.Context(), report); err != nil {
			h.logger.Error("save report", "error", err)
			report.ID = ""
		}
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handlers) getReport(w http.ResponseWriter, r *http.Request) {
	if h.reports == nil {
		writeError(w, http.StatusNotFound, "This report is missing or has expired.")
		return
	}
	report, err := h.reports.GetReport(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, store.ErrReportNotFound) {
			writeError(w, http.StatusNotFound, "This report is missing or has expired.")
			return
		}
		h.logger.Error("get report", "error", err)
		writeError(w, http.StatusInternalServerError, "Could not load this report.")
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handlers) chat(w http.ResponseWriter, r *http.Request) {
	var request chat.Request
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(request.Question) == "" {
		writeError(w, http.StatusBadRequest, "question must not be empty.")
		return
	}
	if utf8.RuneCountInString(request.Question) > 1000 {
		writeError(w, http.StatusBadRequest, "question must not exceed 1,000 characters.")
		return
	}
	if h.llmProvider == "openrouter" {
		if !h.allowAIRequest(w, r) {
			return
		}
	}
	response, err := h.chatService.Ask(r.Context(), request)
	if err != nil {
		if h.writeLLMError(w, err) {
			return
		}
		h.logger.Error("answer contract question", "error", err)
		writeError(w, http.StatusInternalServerError, "Could not answer this question.")
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handlers) allowAIRequest(w http.ResponseWriter, r *http.Request) bool {
	if h.limits == nil {
		return true
	}
	switch h.limits.allow(clientIP(r), time.Now().UTC()) {
	case limitDailyBudget:
		writeError(w, http.StatusTooManyRequests, "Demo limit reached for today. Please try again tomorrow.")
		return false
	case limitHourlyIP:
		writeError(w, http.StatusTooManyRequests, "You’ve reached the demo’s AI request limit. Please try again in an hour.")
		return false
	default:
		return true
	}
}

func (h *Handlers) writeLLMError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, llm.ErrProviderRateLimited) {
		writeError(w, http.StatusTooManyRequests, "The AI service is busy right now. Please try again shortly.")
		return true
	}
	if errors.Is(err, llm.ErrProviderUnavailable) {
		writeError(w, http.StatusBadGateway, "AI analysis is temporarily unavailable. Please try again later.")
		return true
	}
	if errors.Is(err, llm.ErrQuestionTooLong) {
		writeError(w, http.StatusBadRequest, "question must not exceed 1,000 characters.")
		return true
	}
	var tooManyClauses *llm.TooManyClausesError
	if errors.As(err, &tooManyClauses) {
		writeError(w, http.StatusBadRequest, tooManyClauses.Error())
		return true
	}
	return false
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body must contain valid JSON.")
		}
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return errors.New("request body is too large.")
		}
		return errors.New("request body must contain valid JSON.")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object.")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
