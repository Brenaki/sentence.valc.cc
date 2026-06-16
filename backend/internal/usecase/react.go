package usecase

import (
	"context"
	"fmt"

	"sentence.valc.cc/backend/internal/domain"
)

// Reaction values accepted from clients.
const (
	ReactionDislike = 0
	ReactionLike    = 1
)

// ReactToQuote records a like or dislike against a quote.
type ReactToQuote struct {
	repo domain.QuoteRepository
}

// NewReactToQuote wires the use case.
func NewReactToQuote(repo domain.QuoteRepository) *ReactToQuote {
	return &ReactToQuote{repo: repo}
}

// Execute validates the reaction and delegates persistence to the repository.
func (uc *ReactToQuote) Execute(ctx context.Context, quoteID int64, reaction int) error {
	var like bool
	switch reaction {
	case ReactionLike:
		like = true
	case ReactionDislike:
		like = false
	default:
		return domain.ErrInvalidReaction
	}

	if err := uc.repo.AddReaction(ctx, quoteID, like); err != nil {
		return fmt.Errorf("add reaction: %w", err)
	}
	return nil
}
