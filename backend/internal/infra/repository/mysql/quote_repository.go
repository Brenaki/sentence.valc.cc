package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"sentence.valc.cc/backend/internal/domain"
)

// QuoteRepository implements domain.QuoteRepository on top of MySQL.
type QuoteRepository struct {
	db *sql.DB
}

// NewQuoteRepository builds the repository.
func NewQuoteRepository(db *sql.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

const selectColumns = "id, quote, author, work, categories, like_quantity, deslike_quantity, created_at, updated_at"

// FindByDate returns the quote created on the given calendar day.
func (r *QuoteRepository) FindByDate(ctx context.Context, day time.Time) (*domain.Quote, error) {
	const q = "SELECT " + selectColumns + " FROM frases WHERE DATE(created_at) = ? ORDER BY id DESC LIMIT 1"
	row := r.db.QueryRowContext(ctx, q, day.Format("2006-01-02"))
	return scanQuote(row)
}

// FindByID returns a quote by id.
func (r *QuoteRepository) FindByID(ctx context.Context, id int64) (*domain.Quote, error) {
	const q = "SELECT " + selectColumns + " FROM frases WHERE id = ?"
	row := r.db.QueryRowContext(ctx, q, id)
	return scanQuote(row)
}

// Save inserts a new quote and reads it back with DB-managed timestamps.
func (r *QuoteRepository) Save(ctx context.Context, q *domain.Quote) (*domain.Quote, error) {
	cats, err := json.Marshal(normalizeCategories(q.Categories))
	if err != nil {
		return nil, fmt.Errorf("marshal categories: %w", err)
	}
	const ins = `INSERT INTO frases (quote, author, work, categories, like_quantity, deslike_quantity)
	             VALUES (?, ?, ?, ?, 0, 0)`
	res, err := r.db.ExecContext(ctx, ins, q.Quote, q.Author, q.Work, string(cats))
	if err != nil {
		return nil, fmt.Errorf("insert quote: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}
	return r.FindByID(ctx, id)
}

// AddReaction increments the like or dislike counter for a quote.
func (r *QuoteRepository) AddReaction(ctx context.Context, id int64, like bool) error {
	column := "deslike_quantity"
	if like {
		column = "like_quantity"
	}
	q := fmt.Sprintf("UPDATE frases SET %s = %s + 1, updated_at = NOW() WHERE id = ?", column, column)
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("update reaction: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrQuoteNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanQuote(row rowScanner) (*domain.Quote, error) {
	var (
		q    domain.Quote
		cats string
	)
	err := row.Scan(&q.ID, &q.Quote, &q.Author, &q.Work, &cats, &q.LikeQuantity, &q.DislikeQuantity, &q.CreatedAt, &q.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrQuoteNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan quote: %w", err)
	}
	if cats != "" {
		if err := json.Unmarshal([]byte(cats), &q.Categories); err != nil {
			return nil, fmt.Errorf("unmarshal categories: %w", err)
		}
	}
	return &q, nil
}

func normalizeCategories(c []string) []string {
	if c == nil {
		return []string{}
	}
	return c
}
