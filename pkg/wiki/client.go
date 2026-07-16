package wiki

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: config.Timeout}
	}
	return &Client{config: config, httpClient: httpClient}
}

func NormalizePageID(id string) (string, error) {
	normalized := strings.Trim(strings.TrimSpace(id), ":")
	if normalized == "" {
		return "", fmt.Errorf("page id must not be empty")
	}
	if utf8.RuneCountInString(normalized) > 300 {
		return "", fmt.Errorf("page id is too long")
	}
	if strings.Contains(normalized, "..") || strings.ContainsAny(normalized, "/\\?#") {
		return "", fmt.Errorf("page id contains unsupported characters")
	}
	for _, r := range normalized {
		if r < 0x20 || r == 0x7f {
			return "", fmt.Errorf("page id contains unsupported characters")
		}
	}
	return normalized, nil
}

func (c *Client) GetPage(ctx context.Context, rawID string) (*Page, error) {
	id, err := NormalizePageID(rawID)
	if err != nil {
		return nil, err
	}
	pageURL := c.buildURL(url.Values{"id": {id}, "do": {"export_xhtmlbody"}})
	body, err := c.fetch(ctx, pageURL)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return ParsePage(body, id, pageURL)
}

func (c *Client) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("query must not be empty")
	}
	if utf8.RuneCountInString(query) > 300 {
		return nil, fmt.Errorf("query is too long")
	}
	if limit == 0 {
		limit = 10
	}
	if limit < 1 || limit > 50 {
		return nil, fmt.Errorf("limit must be an integer from 1 to 50")
	}
	searchURL := c.buildURL(url.Values{"id": {"index.html"}, "do": {"search"}, "q": {query}})
	body, err := c.fetch(ctx, searchURL)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return ParseSearch(body, c.config.BaseURL, limit)
}

func (c *Client) buildURL(values url.Values) *url.URL {
	result := *c.config.BaseURL
	query := result.Query()
	for key, list := range values {
		query.Del(key)
		for _, value := range list {
			query.Add(key, value)
		}
	}
	result.RawQuery = query.Encode()
	return &result
}

func (c *Client) fetch(ctx context.Context, target *url.URL) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, target.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create Wiki request: %w", err)
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml")
	request.Header.Set("User-Agent", "hpc-wiki-mcp/0.1.0")
	if c.config.Username != "" {
		request.SetBasicAuth(c.config.Username, c.config.Password)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request Wiki: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		response.Body.Close()
		return nil, fmt.Errorf("wiki request failed: HTTP %d", response.StatusCode)
	}
	return response.Body, nil
}
