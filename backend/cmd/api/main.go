package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clauseguard/backend/internal/analysis"
	"clauseguard/backend/internal/chat"
	"clauseguard/backend/internal/config"
	"clauseguard/backend/internal/httpapi"
	"clauseguard/backend/internal/llm"
	"clauseguard/backend/internal/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("load configuration", "error", err)
		os.Exit(1)
	}

	var aiClient llm.Client
	switch cfg.LLMProvider {
	case "mock":
		aiClient = llm.NewMockClient()
	case "openrouter":
		aiClient = llm.NewOpenRouterClient()
	}

	var reportStore store.Store
	var closeStore func(context.Context) error
	if cfg.MongoURI == "" {
		reportStore = store.NewMemoryStore(cfg.ReportTTL)
	} else {
		startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		mongoStore, err := store.NewMongoStore(startupCtx, cfg.MongoURI, cfg.MongoDB, cfg.ReportTTL)
		cancel()
		if err != nil {
			logger.Error("MongoDB is unavailable; saved reports will be unavailable", "error", err)
		} else {
			reportStore = mongoStore
			closeStore = mongoStore.Close
		}
	}

	analysisService := analysis.NewService(aiClient)
	chatService := chat.NewService(aiClient)
	handler := httpapi.NewRouter(cfg, analysisService, chatService, reportStore, logger)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("API server listening", "address", server.Addr, "llm_provider", cfg.LLMProvider)
		serverErrors <- server.ListenAndServe()
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("API server failed", "error", err)
			os.Exit(1)
		}
	case <-signalCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
			_ = server.Close()
		}
		if closeStore != nil {
			if err := closeStore(shutdownCtx); err != nil {
				logger.Error("close report store", "error", err)
			}
		}
	}
}
