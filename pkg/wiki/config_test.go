package wiki

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("HPC_WIKI_BASE_URL", "https://wiki.example/doku.php")
	t.Setenv("HPC_WIKI_USERNAME", "tester")
	t.Setenv("HPC_WIKI_PASSWORD", "secret")
	t.Setenv("HPC_WIKI_TIMEOUT_MS", "2500")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if got := config.BaseURL.String(); got != "https://wiki.example/doku.php" {
		t.Fatalf("BaseURL = %q", got)
	}
	if config.Username != "tester" || config.Password != "secret" {
		t.Fatal("credentials were not loaded")
	}
	if config.Timeout != 2500*time.Millisecond {
		t.Fatalf("Timeout = %v", config.Timeout)
	}
}

func TestLoadConfigRejectsPartialCredentials(t *testing.T) {
	t.Setenv("HPC_WIKI_USERNAME", "tester")
	t.Setenv("HPC_WIKI_PASSWORD", "")
	if _, err := LoadConfig(); err == nil {
		t.Fatal("LoadConfig() succeeded with partial credentials")
	}
}
