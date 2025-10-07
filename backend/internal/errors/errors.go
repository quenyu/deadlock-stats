package errors

import "errors"

var (
	// --- Player-related errors ---
	ErrPlayerNotFound    = errors.New("player not found")
	ErrInvalidSteamID    = errors.New("invalid steam id")
	ErrInvalidQuery      = errors.New("invalid query parameter")
	ErrRateLimited       = errors.New("rate limited")
	ErrAPIUnavailable    = errors.New("external api unavailable")
	ErrPlayerDataMissing = errors.New("player data missing")

	// --- Auth-related errors ---
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrSessionExpired      = errors.New("session expired")
	ErrSteamAuthFailed     = errors.New("steam authentication failed")
	ErrJWTGenerationFailed = errors.New("failed to generate JWT token")

	// --- Crosshair-related errors ---
	ErrCrosshairNotFound  = errors.New("crosshair not found")
	ErrInvalidCrosshairID = errors.New("invalid crosshair ID")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrCrosshairForbidden = errors.New("not allowed to modify this crosshair")
	ErrInvalidRequestBody = errors.New("invalid request body")

	// --- Match / Search-related errors ---
	ErrMatchNotFound   = errors.New("match not found")
	ErrInvalidSearch   = errors.New("invalid search type")
	ErrNoSearchResults = errors.New("no results found")

	// --- System / Internal errors ---
	ErrDatabaseError   = errors.New("database operation failed")
	ErrCacheError      = errors.New("cache operation failed")
	ErrUnknownInternal = errors.New("unknown internal error")
)
