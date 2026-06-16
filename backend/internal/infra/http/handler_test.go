package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sentence.valc.cc/backend/internal/domain"
	apihttp "sentence.valc.cc/backend/internal/infra/http"
)

type stubQOTD struct {
	quote *domain.Quote
	err   error
}

func (s stubQOTD) Execute(_ context.Context) (*domain.Quote, error) { return s.quote, s.err }

type stubReact struct {
	err      error
	gotID    int64
	gotValue int
}

func (s *stubReact) Execute(_ context.Context, id int64, reaction int) error {
	s.gotID = id
	s.gotValue = reaction
	return s.err
}

func newServer(qotd apihttp.QuoteOfTheDayService, react apihttp.ReactionService) http.Handler {
	return apihttp.NewRouter(apihttp.NewHandler(qotd, react))
}

func TestQuoteOfTheDay_OK(t *testing.T) {
	q := &domain.Quote{ID: 7, Quote: "x", Author: "a", Work: "w", Categories: []string{"life"}, LikeQuantity: 2, DislikeQuantity: 1}
	srv := newServer(stubQOTD{quote: q}, &stubReact{})

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/quote-of-the-day", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["id"].(float64) != 7 || body["quote"] != "x" {
		t.Errorf("unexpected body %v", body)
	}
	if body["like_quantity"].(float64) != 2 || body["deslike_quantity"].(float64) != 1 {
		t.Errorf("reaction counters wrong: %v", body)
	}
}

func TestQuoteOfTheDay_ProviderError(t *testing.T) {
	srv := newServer(stubQOTD{err: errors.New("boom")}, &stubReact{})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/quote-of-the-day", nil))
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rec.Code)
	}
}

func TestReact_OK(t *testing.T) {
	react := &stubReact{}
	srv := newServer(stubQOTD{}, react)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/quotes/9/reactions", strings.NewReader(`{"reaction":1}`))
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if react.gotID != 9 || react.gotValue != 1 {
		t.Errorf("usecase got id=%d value=%d", react.gotID, react.gotValue)
	}
}

func TestReact_InvalidReaction(t *testing.T) {
	srv := newServer(stubQOTD{}, &stubReact{err: domain.ErrInvalidReaction})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/quotes/9/reactions", strings.NewReader(`{"reaction":7}`))
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestReact_NotFound(t *testing.T) {
	srv := newServer(stubQOTD{}, &stubReact{err: domain.ErrQuoteNotFound})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/quotes/9/reactions", strings.NewReader(`{"reaction":1}`))
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestReact_BadID(t *testing.T) {
	srv := newServer(stubQOTD{}, &stubReact{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/quotes/abc/reactions", strings.NewReader(`{"reaction":1}`))
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestReact_MalformedBody(t *testing.T) {
	srv := newServer(stubQOTD{}, &stubReact{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/quotes/9/reactions", strings.NewReader(`not-json`))
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
