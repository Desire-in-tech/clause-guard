package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadDefaultsToMockWithDemoLimits(t *testing.T) {
	setConfigEnvironment(t)
	t.Setenv("LLM_PROVIDER", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.LLMProvider != "mock" {
		t.Fatalf("LLMProvider = %q, want mock", cfg.LLMProvider)
	}
	if cfg.LLMTimeout != 45*time.Second || cfg.MaxClauses != 40 ||
		cfg.DailyBudget != 40 || cfg.RateLimitPerHour != 10 {
		t.Fatalf("unexpected default limits: %+v", cfg)
	}
}

func TestOpenRouterRequiresKeyAndPrimaryModel(t *testing.T) {
	t.Run("missing key", func(t *testing.T) {
		setConfigEnvironment(t)
		t.Setenv("LLM_PROVIDER", "openrouter")
		t.Setenv("MODEL_ID", "primary")

		_, err := Load()
		if err == nil || !strings.Contains(err.Error(), "OPENROUTER_API_KEY is required") {
			t.Fatalf("Load() error = %v, want missing-key error", err)
		}
	})

	t.Run("missing model", func(t *testing.T) {
		setConfigEnvironment(t)
		t.Setenv("LLM_PROVIDER", "openrouter")
		t.Setenv("OPENROUTER_API_KEY", "test-only-key")

		_, err := Load()
		if err == nil || !strings.Contains(err.Error(), "MODEL_ID is required") {
			t.Fatalf("Load() error = %v, want missing-model error", err)
		}
	})
}

func TestLoadOpenRouterFallbacksAndOverrides(t *testing.T) {
	setConfigEnvironment(t)
	t.Setenv("LLM_PROVIDER", "openrouter")
	t.Setenv("OPENROUTER_API_KEY", "test-only-key")
	t.Setenv("MODEL_ID", "primary-model")
	t.Setenv("MODEL_FALLBACKS", " fallback-a, ,fallback-b ")
	t.Setenv("LLM_TIMEOUT_SECONDS", "30")
	t.Setenv("MAX_CLAUSES", "25")
	t.Setenv("LLM_DAILY_REQUEST_BUDGET", "18")
	t.Setenv("AI_RATE_LIMIT_PER_HOUR", "4")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if strings.Join(cfg.ModelFallbacks, ",") != "fallback-a,fallback-b" {
		t.Fatalf("ModelFallbacks = %v", cfg.ModelFallbacks)
	}
	if cfg.LLMTimeout != 30*time.Second || cfg.MaxClauses != 25 ||
		cfg.DailyBudget != 18 || cfg.RateLimitPerHour != 4 {
		t.Fatalf("unexpected configured limits: %+v", cfg)
	}
}

func setConfigEnvironment(t *testing.T) {
	t.Helper()
	for key, value := range map[string]string{
		"LLM_PROVIDER":             "mock",
		"OPENROUTER_API_KEY":       "",
		"MODEL_ID":                 "",
		"MODEL_FALLBACKS":          "",
		"LLM_TIMEOUT_SECONDS":      "",
		"MAX_CLAUSES":              "",
		"LLM_DAILY_REQUEST_BUDGET": "",
		"AI_RATE_LIMIT_PER_HOUR":   "",
		"REPORT_TTL_DAYS":          "7",
	} {
		t.Setenv(key, value)
	}
}
