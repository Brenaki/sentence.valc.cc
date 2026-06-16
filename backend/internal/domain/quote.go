package domain

import (
	"context"
	"errors"
	"time"
)

// Reaction sentinel errors and domain errors.
var (
	// ErrQuoteNotFound is returned by a QuoteRepository when no quote matches.
	ErrQuoteNotFound = errors.New("quote not found")
	// ErrInvalidReaction is returned when a reaction value is not 0 or 1.
	ErrInvalidReaction = errors.New("invalid reaction: must be 0 (dislike) or 1 (like)")
)

// Quote is the core domain entity, mapping to the `frases` table.
type Quote struct {
	ID              int64
	Quote           string
	Author          string
	Work            string
	Categories      []string
	LikeQuantity    int
	DislikeQuantity int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// QuoteRepository is the persistence port (driven adapter).
type QuoteRepository interface {
	// FindByDate returns the quote created on the given calendar day.
	// It returns ErrQuoteNotFound when none exists for that day.
	FindByDate(ctx context.Context, day time.Time) (*Quote, error)
	// FindByID returns a quote by its identifier, or ErrQuoteNotFound.
	FindByID(ctx context.Context, id int64) (*Quote, error)
	// Save persists a new quote and returns it with ID/timestamps populated.
	Save(ctx context.Context, q *Quote) (*Quote, error)
	// AddReaction increments the like (true) or dislike (false) counter.
	AddReaction(ctx context.Context, id int64, like bool) error
}

// QuoteProvider is the external quote source port (driven adapter).
type QuoteProvider interface {
	// QuoteOfTheDay fetches the current quote of the day from the source.
	QuoteOfTheDay(ctx context.Context) (*Quote, error)
}

// Clock abstracts time for testability (dependency inversion).
type Clock interface {
	Now() time.Time
}

// RealClock is the production Clock backed by time.Now.
type RealClock struct{}

// Now returns the current local time.
func (RealClock) Now() time.Time { return time.Now() }
