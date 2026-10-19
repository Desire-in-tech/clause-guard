package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"clauseguard/backend/internal/analysis"
	"clauseguard/backend/internal/chat"
	"clauseguard/backend/internal/config"
	"clauseguard/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	cfg config.Config,
	analysisService *analysis.Service,
	chatService *chat.Service,
	reportStore store.Store,
	logger *slog.Logger,
) http.Handler {
	router := chi.NewRouter()
	handler := &Handlers{
		analysis:    analysisService,
		chatService: chatService,
		reports:     reportStore,
		llmProvider: cfg.LLMProvider,
		logger:      logger,
	}

	router.Use(cors(cfg.FrontendOrigin))
	router.Use(requestLogger(logger))
	router.Get("/api/health", handler.health)
	router.Post("/api/analyze", handler.analyze)
	router.Get("/api/reports/{id}", handler.getReport)
	router.Post("/api/chat", handler.chat)
	router.Post("/api/extract", handler.extract)
	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "Route not found.")
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed.")
	})
	return router
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)
			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", recorder.status,
				"duration", time.Since(started),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
