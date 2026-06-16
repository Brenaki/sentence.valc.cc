package usecase_test

import (
	"context"
	"time"

	"sentence.valc.cc/backend/internal/domain"
)

// fakeRepo is an in-memory QuoteRepository for use-case tests.
type fakeRepo struct {
	byDate    map[string]*domain.Quote
	byID      map[int64]*domain.Quote
	saveErr   error
	findErr   error
	reactErr  error
	saved     []*domain.Quote
	reactions []reaction
	nextID    int64
}

type reaction struct {
	id   int64
	like bool
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		byDate: map[string]*domain.Quote{},
		byID:   map[int64]*domain.Quote{},
		nextID: 1,
	}
}

func dayKey(t time.Time) string { return t.Format("2006-01-02") }

func (f *fakeRepo) FindByDate(_ context.Context, day time.Time) (*domain.Quote, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	q, ok := f.byDate[dayKey(day)]
	if !ok {
		return nil, domain.ErrQuoteNotFound
	}
	return q, nil
}

func (f *fakeRepo) FindByID(_ context.Context, id int64) (*domain.Quote, error) {
	q, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrQuoteNotFound
	}
	return q, nil
}

func (f *fakeRepo) Save(_ context.Context, q *domain.Quote) (*domain.Quote, error) {
	if f.saveErr != nil {
		return nil, f.saveErr
	}
	saved := *q
	saved.ID = f.nextID
	f.nextID++
	now := time.Now()
	saved.CreatedAt = now
	saved.UpdatedAt = now
	f.byDate[dayKey(now)] = &saved
	f.byID[saved.ID] = &saved
	f.saved = append(f.saved, &saved)
	return &saved, nil
}

func (f *fakeRepo) AddReaction(_ context.Context, id int64, like bool) error {
	if f.reactErr != nil {
		return f.reactErr
	}
	f.reactions = append(f.reactions, reaction{id: id, like: like})
	if q, ok := f.byID[id]; ok {
		if like {
			q.LikeQuantity++
		} else {
			q.DislikeQuantity++
		}
	}
	return nil
}

// fakeProvider is a stub QuoteProvider.
type fakeProvider struct {
	quote  *domain.Quote
	err    error
	called int
}

func (p *fakeProvider) QuoteOfTheDay(_ context.Context) (*domain.Quote, error) {
	p.called++
	if p.err != nil {
		return nil, p.err
	}
	return p.quote, nil
}

// fixedClock returns a constant time.
type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }
