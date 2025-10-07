package errors

import "errors"

var (
	// --- Player-related errors ---
	ErrPlayerNotFound = errors.New("player not found")
	ErrInvalidSteamID = errors.New("invalid steam id")
	ErrInvalidQuery   = errors.New("invalid query parameter")
	ErrRateLimited    = errors.New("rate limited")
	ErrAPIUnavailable = errors.New("external api unavailable")

	// --- Auth-related errors ---
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrSessionExpired      = errors.New("session expired")
	ErrSteamAuthFailed     = errors.New("steam authentication failed")
	ErrJWTGenerationFailed = errors.New("failed to generate JWT token")
)
