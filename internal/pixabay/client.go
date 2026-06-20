package pixabay

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client calls the PixaBay REST API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithBaseURL overrides the PixaBay API base URL (mainly for tests).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL != "" {
			c.baseURL = baseURL
		}
	}
}

// NewClient creates a PixaBay API client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Search performs an image search against the PixaBay API.
func (c *Client) Search(ctx context.Context, params SearchParams) (SearchResult, error) {
	if strings.TrimSpace(params.Query) == "" {
		return SearchResult{}, ErrEmptyQuery
	}

	params = normalizeSearchParams(params)

	reqURL, err := c.buildSearchURL(params)
	if err != nil {
		return SearchResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return SearchResult{}, fmt.Errorf("pixabay: create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchResult{}, fmt.Errorf("pixabay: request failed: %w", err)
	}
	defer resp.Body.Close()

	rateLimit := parseRateLimit(resp.Header)

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return SearchResult{}, fmt.Errorf("pixabay: read response: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return SearchResult{RateLimit: rateLimit}, fmt.Errorf("%w: %s", ErrRateLimited, strings.TrimSpace(string(body)))
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return SearchResult{RateLimit: rateLimit}, &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(string(body)),
		}
	}

	var apiResp apiSearchResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return SearchResult{RateLimit: rateLimit}, fmt.Errorf("pixabay: decode response: %w", err)
	}

	return toSearchResult(apiResp, rateLimit), nil
}

func (c *Client) buildSearchURL(params SearchParams) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("pixabay: parse base url: %w", err)
	}

	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("q", params.Query)
	q.Set("order", params.Order)
	q.Set("page", strconv.Itoa(params.Page))
	q.Set("per_page", strconv.Itoa(params.PerPage))
	q.Set("image_type", params.ImageType)
	q.Set("safesearch", strconv.FormatBool(params.SafeSearch))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func parseRateLimit(header http.Header) RateLimit {
	return RateLimit{
		Limit:     parseHeaderInt(header.Get("X-RateLimit-Limit")),
		Remaining: parseHeaderInt(header.Get("X-RateLimit-Remaining")),
		Reset:     parseHeaderInt(header.Get("X-RateLimit-Reset")),
	}
}

func parseHeaderInt(value string) int {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return n
}
