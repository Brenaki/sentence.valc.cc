package ninja_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"sentence.valc.cc/backend/internal/infra/provider/ninja"
)

const sampleBody = `[{"quote":"Be the change, you seek from society","author":"Kandarp Gandhi","work":"Buddhist Banker","categories":["inspirational","philosophy"]}]`

func TestQuoteOfTheDay_ParsesResponse(t *testing.T) {
	var gotKey string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-Api-Key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleBody))
	}))
	defer ts.Close()

	p := ninja.New("secret-key", ninja.WithBaseURL(ts.URL), ninja.WithHTTPClient(ts.Client()))
	q, err := p.QuoteOfTheDay(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotKey != "secret-key" {
		t.Errorf("X-Api-Key header = %q", gotKey)
	}
	if q.Author != "Kandarp Gandhi" || q.Quote == "" {
		t.Errorf("bad parse: %+v", q)
	}
	if len(q.Categories) != 2 {
		t.Errorf("expected 2 categories, got %v", q.Categories)
	}
}

func TestQuoteOfTheDay_EmptyArray(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	p := ninja.New("k", ninja.WithBaseURL(ts.URL), ninja.WithHTTPClient(ts.Client()))
	if _, err := p.QuoteOfTheDay(context.Background()); err == nil {
		t.Fatal("expected error on empty array")
	}
}

func TestQuoteOfTheDay_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	p := ninja.New("k", ninja.WithBaseURL(ts.URL), ninja.WithHTTPClient(ts.Client()))
	if _, err := p.QuoteOfTheDay(context.Background()); err == nil {
		t.Fatal("expected error on non-200 status")
	}
}
