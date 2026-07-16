package wiki

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultBaseURL = "http://www.hpc.is.uec.ac.jp/wiki/hpcwiki2/doku.php"

type Config struct {
	BaseURL  *url.URL
	Username string
	Password string
	Timeout  time.Duration
}

func LoadConfig() (*Config, error) {
	baseURL, err := url.Parse(envOrDefault("HPC_WIKI_BASE_URL", defaultBaseURL))
	if err != nil {
		return nil, fmt.Errorf("parse HPC_WIKI_BASE_URL: %w", err)
	}
	if baseURL.Scheme != "http" && baseURL.Scheme != "https" {
		return nil, fmt.Errorf("HPC_WIKI_BASE_URL must use http or https")
	}

	username := strings.TrimSpace(os.Getenv("HPC_WIKI_USERNAME"))
	password := os.Getenv("HPC_WIKI_PASSWORD")
	if (username == "") != (password == "") {
		return nil, fmt.Errorf("HPC_WIKI_USERNAME and HPC_WIKI_PASSWORD must be set together")
	}

	timeoutMS, err := strconv.Atoi(envOrDefault("HPC_WIKI_TIMEOUT_MS", "15000"))
	if err != nil || timeoutMS < 1000 || timeoutMS > 120000 {
		return nil, fmt.Errorf("HPC_WIKI_TIMEOUT_MS must be an integer from 1000 to 120000")
	}

	return &Config{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		Timeout:  time.Duration(timeoutMS) * time.Millisecond,
	}, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
