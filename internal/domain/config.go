package domain

import (
	"net/url"
	"testing"
)

type (
	// Config represent a global micropub instance configuration.
	Config struct {
		HTTP     ConfigHTTP `envPrefix:"HTTP_"`
		MediaDir string     `env:"MEDIA_DIR" envDefault:"media"`
	}

	// ConfigHTTP represents HTTP configs which used for instance serving
	// and Location header responses.
	ConfigHTTP struct {
		Bind  string `env:"BIND" envDefault:":3000"`
		Host  string `env:"HOST" envDefault:"localhost:3000"`
		Proto string `env:"PROTO" envDefault:"http"`
	}
)

// TestConfig returns a valid Config for tests.
func TestConfig(tb testing.TB) *Config {
	tb.Helper()

	return &Config{
		HTTP: ConfigHTTP{
			Bind:  ":3000",
			Host:  "example.com",
			Proto: "https",
		},
		MediaDir: "media",
	}
}

// BaseURL returns root *url.URL based on provided proto and host.
func (c ConfigHTTP) BaseURL() *url.URL {
	return &url.URL{
		Scheme: c.Proto,
		Host:   c.Host,
		Path:   "/",
	}
}
