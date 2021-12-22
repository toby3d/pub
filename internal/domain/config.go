package domain

import "testing"

type Config struct {
	BaseURL string
}

func TestConfig(tb testing.TB) *Config {
	tb.Helper()

	return &Config{
		BaseURL: "https://example.com/",
	}
}
