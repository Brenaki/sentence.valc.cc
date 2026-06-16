package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"sentence.valc.cc/backend/internal/domain"
)

// QuoteOfTheDayService is the inbound port consumed by the handler.
type QuoteOfTheDayService interface {
	Execute(ctx context.Context) (*domain.Quote, error)
}

// ReactionService is the inbound port for recording reactions.
type ReactionService interface {
	Execute(ctx context.Context, quoteID int64, reaction int) error
}

// Handler exposes HTTP endpoints for the quote use cases.
type Handler struct {
	qotd  QuoteOfTheDayService
	react ReactionService
}

// NewHandler wires the handler with its services.
func NewHandler(qotd QuoteOfTheDayService, react ReactionService) *Handler {
	return &Handler{qotd: qotd, react: react}
}

// NewRouter builds the application router (Go 1.22+ method-aware mux).
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /quote-of-the-day", h.QuoteOfTheDay)
	mux.HandleFunc("POST /quotes/{id}/reactions", h.React)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

type quoteResponse struct {
	ID              int64    `json:"id"`
	Quote           string   `json:"quote"`
	Author          string   `json:"author"`
	Work            string   `json:"work"`
	Categories      []string `json:"categories"`
	LikeQuantity    int      `json:"like_quantity"`
	DislikeQuantity int      `json:"deslike_quantity"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

func toResponse(q *domain.Quote) quoteResponse {
	cats := q.Categories
	if cats == nil {
		cats = []string{}
	}
	return quoteResponse{
		ID:              q.ID,
		Quote:           q.Quote,
		Author:          q.Author,
		Work:            q.Work,
		Categories:      cats,
		LikeQuantity:    q.LikeQuantity,
		DislikeQuantity: q.DislikeQuantity,
		CreatedAt:       q.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       q.UpdatedAt.Format(time.RFC3339),
	}
}

// QuoteOfTheDay handles GET /quote-of-the-day.
func (h *Handler) QuoteOfTheDay(w http.ResponseWriter, r *http.Request) {
	q, err := h.qotd.Execute(r.Context())
	if err != nil {
		slog.Error("quote of the day", "err", err)
		writeError(w, http.StatusBadGateway, "could not resolve quote of the day")
		return
	}
	writeJSON(w, http.StatusOK, toResponse(q))
}

type reactionRequest struct {
	Reaction int `json:"reaction"`
}

// React handles POST /quotes/{id}/reactions.
func (h *Handler) React(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid quote id")
		return
	}

	var req reactionRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.react.Execute(r.Context(), id, req.Reaction)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, domain.ErrInvalidReaction):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrQuoteNotFound):
		writeError(w, http.StatusNotFound, "quote not found")
	default:
		slog.Error("react to quote", "err", err)
		writeError(w, http.StatusInternalServerError, "could not record reaction")
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
