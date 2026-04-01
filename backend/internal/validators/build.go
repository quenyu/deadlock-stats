package validators

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
)

func ValidateHeroName(heroName string) error {
	heroName = strings.TrimSpace(heroName)

	if heroName == "" {
		return cErrors.ErrInvalidHeroName
	}

	if len(heroName) < 2 {
		return cErrors.ErrFieldTooShort
	}

	if len(heroName) > 50 {
		return cErrors.ErrFieldTooLong
	}

	heroNameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-&']+$`)
	if !heroNameRegex.MatchString(heroName) {
		return cErrors.ErrInvalidHeroName
	}

	return nil
}

func validateHeroNameTag(fl validator.FieldLevel) bool {
	heroName := fl.Field().String()
	return ValidateHeroName(heroName) == nil
}

func ValidateBuildTitle(title string) error {
	title = strings.TrimSpace(title)

	if title == "" {
		return cErrors.ErrInvalidBuildTitle
	}

	if len(title) < 5 {
		return cErrors.ErrFieldTooShort
	}

	if len(title) > 100 {
		return cErrors.ErrFieldTooLong
	}

	return nil
}

func validateBuildTitleTag(fl validator.FieldLevel) bool {
	title := fl.Field().String()
	return ValidateBuildTitle(title) == nil
}

func ValidateBuildDescription(description string) error {
	return ValidateDescription(description, 2000)
}

func ValidatePatchVersion(version string) error {
	version = strings.TrimSpace(version)

	if version == "" {
		return cErrors.ErrInvalidPatchVer
	}

	versionRegex := regexp.MustCompile(`^\d+\.\d+(\.\d+)?$`)
	if !versionRegex.MatchString(version) {
		return cErrors.ErrInvalidPatchVer
	}

	return nil
}

func validatePatchVersionTag(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	return ValidatePatchVersion(version) == nil
}

func ValidateItemID(itemID int) error {
	if itemID <= 0 {
		return cErrors.ErrInvalidItemID
	}
	return nil
}

func ValidateAbilitySlot(slot int) error {
	if slot < 1 || slot > 4 {
		return cErrors.ErrInvalidAbility
	}
	return nil
}
