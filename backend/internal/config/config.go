package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port           string
	FrontendOrigin string
	LLMProvider    string
	OpenRouterKey  string
	ModelID        string
	MongoURI       string
	MongoDB        string
	ReportTTL      time.Duration
}

func Load() (Config, error) {
	ttlDays := 7
	if value := strings.TrimSpace(os.Getenv("REPORT_TTL_DAYS")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return Config{}, fmt.Errorf("REPORT_TTL_DAYS must be a positive integer")
		}
		ttlDays = parsed
	}

	provider := strings.ToLower(envOr("LLM_PROVIDER", "mock"))
	if provider != "mock" && provider != "openrouter" {
		return Config{}, fmt.Errorf("LLM_PROVIDER must be mock or openrouter")
	}

	return Config{
		Port:           envOr("PORT", "8080"),
		FrontendOrigin: strings.TrimSpace(os.Getenv("FRONTEND_ORIGIN")),
		LLMProvider:    provider,
		OpenRouterKey:  strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY")),
		ModelID:        strings.TrimSpace(os.Getenv("MODEL_ID")),
		MongoURI:       strings.TrimSpace(os.Getenv("MONGODB_URI")),
		MongoDB:        envOr("MONGODB_DB", "clauseguard"),
		ReportTTL:      time.Duration(ttlDays) * 24 * time.Hour,
	}, nil
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
