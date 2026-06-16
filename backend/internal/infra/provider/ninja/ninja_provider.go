package ninja

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sentence.valc.cc/backend/internal/domain"
)

const defaultBaseURL = "https://api.api-ninjas.com/v2/quoteoftheday"

// Provider implements domain.QuoteProvider against api-ninjas.
type Provider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// Option configures the Provider.
type Option func(*Provider)

// WithBaseURL overrides the upstream URL (used in tests).
func WithBaseURL(url string) Option { return func(p *Provider) { p.baseURL = url } }

// WithHTTPClient overrides the HTTP client.
func WithHTTPClient(c *http.Client) Option { return func(p *Provider) { p.client = c } }

// New builds a Provider with the given API key.
func New(apiKey string, opts ...Option) *Provider {
	p := &Provider{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

type ninjaQuote struct {
	Quote      string   `json:"quote"`
	Author     string   `json:"author"`
	Work       string   `json:"work"`
	Categories []string `json:"categories"`
}

// QuoteOfTheDay fetches the quote of the day from api-ninjas.
func (p *Provider) QuoteOfTheDay(ctx context.Context) (*domain.Quote, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Api-Key", p.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call api-ninjas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("api-ninjas status %d: %s", resp.StatusCode, body)
	}

	var quotes []ninjaQuote
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(quotes) == 0 {
		return nil, fmt.Errorf("api-ninjas returned no quotes")
	}

	n := quotes[0]
	return &domain.Quote{
		Quote:      n.Quote,
		Author:     n.Author,
		Work:       n.Work,
		Categories: n.Categories,
	}, nil
}
