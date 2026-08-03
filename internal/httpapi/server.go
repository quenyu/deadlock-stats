package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/quenyu/deadlock-stats/internal/ghost"
)

type GhostBuilder interface {
	Build(ctx context.Context, accountID string, matchID int64) (*ghost.Report, error)
}

type Server struct {
	ghost  GhostBuilder
	logger *slog.Logger
	mux    *http.ServeMux
}

func New(ghostBuilder GhostBuilder, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}

	s := &Server{
		ghost:  ghostBuilder,
		logger: logger,
		mux:    http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.securityHeaders(s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /api/v1/players/{accountID}/matches/{matchID}/ghost", s.ghostMatch)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) ghostMatch(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("accountID")
	if !isDigits(accountID) || len(accountID) > 20 {
		writeError(w, http.StatusBadRequest, "invalid_account_id", "accountID must contain 1-20 digits")
		return
	}

	matchID, err := strconv.ParseInt(r.PathValue("matchID"), 10, 64)
	if err != nil || matchID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_match_id", "matchID must be a positive integer")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()

	report, err := s.ghost.Build(ctx, accountID, matchID)
	if err != nil {
		s.handleGhostError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handleGhostError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ghost.ErrTargetNotFound):
		writeError(w, http.StatusNotFound, "target_match_not_found", err.Error())
	case errors.Is(err, ghost.ErrReferenceNotFound):
		writeError(w, http.StatusUnprocessableEntity, "reference_match_not_found", err.Error())
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
		writeError(w, http.StatusGatewayTimeout, "upstream_timeout", "match history request timed out")
	default:
		s.logger.Error("ghost match failed", "error", err)
		writeError(w, http.StatusBadGateway, "upstream_error", "could not build ghost match report")
	}
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Default().Error("encode JSON response", "error", err)
	}
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
