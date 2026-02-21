package validation

import (
	"regexp"
	"strings"

	"github.com/safar/microservices-demo/gateway/internal/errors"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// ValidateEmail checks if an email is valid
func ValidateEmail(email string) *errors.ValidationError {
	if email == "" {
		return &errors.ValidationError{
			Field:   "email",
			Message: "Email is required",
		}
	}

	if !emailRegex.MatchString(email) {
		return &errors.ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		}
	}

	return nil
}

// ValidatePassword checks if a password meets requirements
func ValidatePassword(password string) *errors.ValidationError {
	if password == "" {
		return &errors.ValidationError{
			Field:   "password",
			Message: "Password is required",
		}
	}

	if len(password) < 8 {
		return &errors.ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		}
	}

	return nil
}

// ValidateRequired checks if a required field is provided
func ValidateRequired(field, value string) *errors.ValidationError {
	if strings.TrimSpace(value) == "" {
		return &errors.ValidationError{
			Field:   field,
			Message: field + " is required",
		}
	}
	return nil
}

// ValidateMinLength checks if a string meets minimum length
func ValidateMinLength(field, value string, minLength int) *errors.ValidationError {
	if len(value) < minLength {
		return &errors.ValidationError{
			Field:   field,
			Message: field + " must be at least " + string(rune(minLength)) + " characters long",
		}
	}
	return nil
}

// ValidateMaxLength checks if a string doesn't exceed maximum length
func ValidateMaxLength(field, value string, maxLength int) *errors.ValidationError {
	if len(value) > maxLength {
		return &errors.ValidationError{
			Field:   field,
			Message: field + " must not exceed " + string(rune(maxLength)) + " characters",
		}
	}
	return nil
}

// ValidatePositive checks if a number is positive
func ValidatePositive(field string, value int64) *errors.ValidationError {
	if value <= 0 {
		return &errors.ValidationError{
			Field:   field,
			Message: field + " must be positive",
		}
	}
	return nil
}

// ValidateRange checks if a number is within a range
func ValidateRange(field string, value, min, max int64) *errors.ValidationError {
	if value < min || value > max {
		return &errors.ValidationError{
			Field:   field,
			Message: field + " must be between " + string(rune(min)) + " and " + string(rune(max)),
		}
	}
	return nil
}

// Validate runs multiple validation checks and collects errors
func Validate(checks ...func() *errors.ValidationError) []errors.ValidationError {
	var validationErrors []errors.ValidationError

	for _, check := range checks {
		if err := check(); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}

	return validationErrors
}
