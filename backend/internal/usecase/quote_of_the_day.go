package usecase

import (
	"context"
	"errors"
	"fmt"

	"sentence.valc.cc/backend/internal/domain"
)

// GetQuoteOfTheDay returns today's quote, fetching from the provider and
// persisting it only when the repository has no record for the current day.
type GetQuoteOfTheDay struct {
	repo     domain.QuoteRepository
	provider domain.QuoteProvider
	clock    domain.Clock
}

// NewGetQuoteOfTheDay wires the use case with its dependencies.
func NewGetQuoteOfTheDay(repo domain.QuoteRepository, provider domain.QuoteProvider, clock domain.Clock) *GetQuoteOfTheDay {
	return &GetQuoteOfTheDay{repo: repo, provider: provider, clock: clock}
}

// Execute resolves the quote of the day.
func (uc *GetQuoteOfTheDay) Execute(ctx context.Context) (*domain.Quote, error) {
	today := uc.clock.Now()

	q, err := uc.repo.FindByDate(ctx, today)
	if err == nil {
		return q, nil
	}
	if !errors.Is(err, domain.ErrQuoteNotFound) {
		return nil, fmt.Errorf("lookup quote of the day: %w", err)
	}

	fetched, err := uc.provider.QuoteOfTheDay(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch quote from provider: %w", err)
	}

	saved, err := uc.repo.Save(ctx, fetched)
	if err != nil {
		return nil, fmt.Errorf("persist quote: %w", err)
	}
	return saved, nil
}
