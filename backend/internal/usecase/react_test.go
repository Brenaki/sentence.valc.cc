package usecase_test

import (
	"context"
	"errors"
	"testing"

	"sentence.valc.cc/backend/internal/domain"
	"sentence.valc.cc/backend/internal/usecase"
)

func TestReact_Like(t *testing.T) {
	repo := newFakeRepo()
	repo.byID[1] = sampleQuote()
	uc := usecase.NewReactToQuote(repo)

	if err := uc.Execute(context.Background(), 1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.reactions) != 1 || !repo.reactions[0].like {
		t.Fatalf("expected one like reaction, got %+v", repo.reactions)
	}
}

func TestReact_Dislike(t *testing.T) {
	repo := newFakeRepo()
	repo.byID[1] = sampleQuote()
	uc := usecase.NewReactToQuote(repo)

	if err := uc.Execute(context.Background(), 1, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.reactions) != 1 || repo.reactions[0].like {
		t.Fatalf("expected one dislike reaction, got %+v", repo.reactions)
	}
}

func TestReact_InvalidValue(t *testing.T) {
	repo := newFakeRepo()
	uc := usecase.NewReactToQuote(repo)

	err := uc.Execute(context.Background(), 1, 5)
	if !errors.Is(err, domain.ErrInvalidReaction) {
		t.Fatalf("expected ErrInvalidReaction, got %v", err)
	}
	if len(repo.reactions) != 0 {
		t.Errorf("no reaction should be recorded for invalid value")
	}
}

func TestReact_RepoErrorPropagates(t *testing.T) {
	repo := newFakeRepo()
	repo.reactErr = errors.New("db down")
	uc := usecase.NewReactToQuote(repo)

	if err := uc.Execute(context.Background(), 1, 1); err == nil {
		t.Fatal("expected repo error to propagate")
	}
}
