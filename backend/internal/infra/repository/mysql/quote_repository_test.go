package mysql_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"sentence.valc.cc/backend/internal/domain"
	repo "sentence.valc.cc/backend/internal/infra/repository/mysql"
)

func TestFindByDate_Found(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewQuoteRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "quote", "author", "work", "categories", "like_quantity", "deslike_quantity", "created_at", "updated_at"}).
		AddRow(1, "q", "a", "w", `["life","wisdom"]`, 3, 2, now, now)
	mock.ExpectQuery("SELECT (.+) FROM frases WHERE DATE\\(created_at\\)").
		WithArgs(now.Format("2006-01-02")).
		WillReturnRows(rows)

	q, err := r.FindByDate(context.Background(), now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.ID != 1 || len(q.Categories) != 2 || q.LikeQuantity != 3 {
		t.Errorf("bad row mapping: %+v", q)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestFindByDate_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewQuoteRepository(db)

	empty := sqlmock.NewRows([]string{"id", "quote", "author", "work", "categories", "like_quantity", "deslike_quantity", "created_at", "updated_at"})
	mock.ExpectQuery("SELECT (.+) FROM frases").WillReturnRows(empty)
	_, err := r.FindByDate(context.Background(), time.Now())
	if !errors.Is(err, domain.ErrQuoteNotFound) {
		t.Fatalf("expected ErrQuoteNotFound, got %v", err)
	}
}

func TestSave_InsertsAndReadsBack(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewQuoteRepository(db)

	now := time.Now()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO frases")).
		WithArgs("q", "a", "w", `["life"]`).
		WillReturnResult(sqlmock.NewResult(10, 1))
	rows := sqlmock.NewRows([]string{"id", "quote", "author", "work", "categories", "like_quantity", "deslike_quantity", "created_at", "updated_at"}).
		AddRow(10, "q", "a", "w", `["life"]`, 0, 0, now, now)
	mock.ExpectQuery("SELECT (.+) FROM frases WHERE id = ?").WithArgs(int64(10)).WillReturnRows(rows)

	got, err := r.Save(context.Background(), &domain.Quote{Quote: "q", Author: "a", Work: "w", Categories: []string{"life"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 10 {
		t.Errorf("expected id 10, got %d", got.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestAddReaction_Like(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewQuoteRepository(db)

	mock.ExpectExec("UPDATE frases SET like_quantity = like_quantity \\+ 1").
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := r.AddReaction(context.Background(), 5, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestAddReaction_Dislike(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewQuoteRepository(db)

	mock.ExpectExec("UPDATE frases SET deslike_quantity = deslike_quantity \\+ 1").
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := r.AddReaction(context.Background(), 5, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddReaction_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewQuoteRepository(db)

	mock.ExpectExec("UPDATE frases").WithArgs(int64(99)).WillReturnResult(sqlmock.NewResult(0, 0))
	if err := r.AddReaction(context.Background(), 99, true); !errors.Is(err, domain.ErrQuoteNotFound) {
		t.Fatalf("expected ErrQuoteNotFound, got %v", err)
	}
}
