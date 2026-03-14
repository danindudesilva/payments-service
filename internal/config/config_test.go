package config

import "testing"

func TestHTTPAddress(t *testing.T) {
	cfg := Config{HTTPPort: "9999"}

	if got, want := cfg.HTTPAddress(), ":9999"; got != want {
		t.Fatalf("HTTPAddress() = %q, want %q", got, want)
	}
}
