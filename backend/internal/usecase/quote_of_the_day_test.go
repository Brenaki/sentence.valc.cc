package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"sentence.valc.cc/backend/internal/domain"
	"sentence.valc.cc/backend/internal/usecase"
)

func sampleQuote() *domain.Quote {
	return &domain.Quote{
		Quote:      "Be the change, you seek from society",
		Author:     "Kandarp Gandhi",
		Work:       "Buddhist Banker",
		Categories: []string{"inspirational", "wisdom"},
	}
}

func TestGetQuoteOfTheDay_ReturnsFromRepoWhenExists(t *testing.T) {
	repo := newFakeRepo()
	today := time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC)
	existing := sampleQuote()
	existing.ID = 42
	repo.byDate[dayKey(today)] = existing

	provider := &fakeProvider{quote: sampleQuote()}
	uc := usecase.NewGetQuoteOfTheDay(repo, provider, fixedClock{t: today})

	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("expected cached quote ID 42, got %d", got.ID)
	}
	if provider.called != 0 {
		t.Errorf("provider must not be called when cache hit, called=%d", provider.called)
	}
	if len(repo.saved) != 0 {
		t.Errorf("repo.Save must not be called on cache hit")
	}
}

func TestGetQuoteOfTheDay_FetchesAndSavesWhenMissing(t *testing.T) {
	repo := newFakeRepo()
	today := time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC)
	provider := &fakeProvider{quote: sampleQuote()}
	uc := usecase.NewGetQuoteOfTheDay(repo, provider, fixedClock{t: today})

	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.called != 1 {
		t.Errorf("provider should be called once, called=%d", provider.called)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("quote should be saved once, saved=%d", len(repo.saved))
	}
	if got.ID == 0 {
		t.Errorf("returned quote should have persisted ID")
	}
	if got.Author != "Kandarp Gandhi" {
		t.Errorf("unexpected author %q", got.Author)
	}
}

func TestGetQuoteOfTheDay_ProviderErrorPropagates(t *testing.T) {
	repo := newFakeRepo()
	today := time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC)
	provider := &fakeProvider{err: errors.New("boom")}
	uc := usecase.NewGetQuoteOfTheDay(repo, provider, fixedClock{t: today})

	_, err := uc.Execute(context.Background())
	if err == nil {
		t.Fatal("expected error from provider")
	}
	if len(repo.saved) != 0 {
		t.Errorf("nothing should be saved when provider fails")
	}
}

func TestGetQuoteOfTheDay_RepoFindErrorPropagates(t *testing.T) {
	repo := newFakeRepo()
	repo.findErr = errors.New("db down")
	today := time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC)
	provider := &fakeProvider{quote: sampleQuote()}
	uc := usecase.NewGetQuoteOfTheDay(repo, provider, fixedClock{t: today})

	_, err := uc.Execute(context.Background())
	if err == nil {
		t.Fatal("expected error when repo lookup fails for non-NotFound reason")
	}
	if provider.called != 0 {
		t.Errorf("provider must not be called on unexpected repo error")
	}
}
