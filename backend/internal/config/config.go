package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port             string
	FrontendOrigin   string
	LLMProvider      string
	OpenRouterKey    string
	ModelID          string
	ModelFallbacks   []string
	LLMTimeout       time.Duration
	MaxClauses       int
	DailyBudget      int
	RateLimitPerHour int
	MongoURI         string
	MongoDB          string
	ReportTTL        time.Duration
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
	timeoutSeconds, err := positiveIntEnv("LLM_TIMEOUT_SECONDS", 45)
	if err != nil {
		return Config{}, err
	}
	maxClauses, err := positiveIntEnv("MAX_CLAUSES", 40)
	if err != nil {
		return Config{}, err
	}
	dailyBudget, err := positiveIntEnv("LLM_DAILY_REQUEST_BUDGET", 40)
	if err != nil {
		return Config{}, err
	}
	rateLimitPerHour, err := positiveIntEnv("AI_RATE_LIMIT_PER_HOUR", 10)
	if err != nil {
		return Config{}, err
	}

	apiKey := strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY"))
	modelID := strings.TrimSpace(os.Getenv("MODEL_ID"))
	if provider == "openrouter" {
		if apiKey == "" {
			return Config{}, fmt.Errorf("OPENROUTER_API_KEY is required when LLM_PROVIDER=openrouter")
		}
		if modelID == "" {
			return Config{}, fmt.Errorf("MODEL_ID is required when LLM_PROVIDER=openrouter")
		}
	}

	return Config{
		Port:             envOr("PORT", "8080"),
		FrontendOrigin:   strings.TrimSpace(os.Getenv("FRONTEND_ORIGIN")),
		LLMProvider:      provider,
		OpenRouterKey:    apiKey,
		ModelID:          modelID,
		ModelFallbacks:   splitCSV(os.Getenv("MODEL_FALLBACKS")),
		LLMTimeout:       time.Duration(timeoutSeconds) * time.Second,
		MaxClauses:       maxClauses,
		DailyBudget:      dailyBudget,
		RateLimitPerHour: rateLimitPerHour,
		MongoURI:         strings.TrimSpace(os.Getenv("MONGODB_URI")),
		MongoDB:          envOr("MONGODB_DB", "clauseguard"),
		ReportTTL:        time.Duration(ttlDays) * 24 * time.Hour,
	}, nil
}

func positiveIntEnv(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return parsed, nil
}

func splitCSV(value string) []string {
	var entries []string
	for _, entry := range strings.Split(value, ",") {
		if entry = strings.TrimSpace(entry); entry != "" {
			entries = append(entries, entry)
		}
	}
	return entries
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
