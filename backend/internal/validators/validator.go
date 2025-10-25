package validators

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		registerCustomValidators()
	})

	return validate
}

func registerCustomValidators() {
	// Player validators
	validate.RegisterValidation("steamid", validateSteamIDTag)

	// Common validators
	validate.RegisterValidation("username", validateUsernameTag)

	// Build validators
	validate.RegisterValidation("heroname", validateHeroNameTag)
	validate.RegisterValidation("buildtitle", validateBuildTitleTag)
	validate.RegisterValidation("patchver", validatePatchVersionTag)

}

func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}

func ValidateVar(field interface{}, tag string) error {
	return GetValidator().Var(field, tag)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

func FormatValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	if err == nil {
		return validationErrors
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrors = append(validationErrors, ValidationError{
				Field:   fe.Field(),
				Message: getErrorMessage(fe),
				Tag:     fe.Tag(),
				Value:   fe.Param(),
			})
		}
	} else {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "general",
			Message: err.Error(),
			Tag:     "error",
		})
	}

	return validationErrors
}

func getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be no more than %s characters", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fe.Param())
	case "steamid":
		return fmt.Sprintf("%s must be a valid Steam ID", field)
	case "username":
		return fmt.Sprintf("%s must be a valid username (3-50 alphanumeric characters)", field)
	case "heroname":
		return fmt.Sprintf("%s must be a valid hero name", field)
	case "buildtitle":
		return fmt.Sprintf("%s must be a valid build title (5-100 characters)", field)
	case "patchver":
		return fmt.Sprintf("%s must be a valid patch version (e.g., 1.2.3)", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s validation failed on '%s'", field, tag)
	}
}
