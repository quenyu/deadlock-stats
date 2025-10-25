package validators

import (
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
)

func ValidateCrosshairTitle(title string) error {
	if len(title) < 3 {
		return cErrors.ErrFieldTooShort
	}
	if len(title) > 100 {
		return cErrors.ErrFieldTooLong
	}
	return nil
}

func ValidateCrosshairDescription(description string) error {
	return ValidateDescription(description, 2000)
}

func ValidateOpacity(value float64) error {
	if value < 0 || value > 1 {
		return cErrors.ErrInvalidRequestBody
	}
	return nil
}

func ValidateIntRange(value, min, max int) error {
	if value < min || value > max {
		return cErrors.ErrInvalidRequestBody
	}
	return nil
}
