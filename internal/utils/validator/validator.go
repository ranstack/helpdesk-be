package validator

import (
	"fmt"
	"helpdesk/internal/utils/errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

func (v *Validator) ToAppError() *errors.AppError {
	if v.Valid() {
		return nil
	}

	details := make(map[string]interface{})
	for field, msg := range v.Errors {
		details[field] = msg
	}

	return errors.Validation("Validation failed").WithDetails(details)
}

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinLength(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

func MaxLength(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func InRange(value string, min, max int) bool {
	length := utf8.RuneCountInString(value)
	return length >= min && length <= max
}

func ValidateEmail(value string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(value)
}

func ValidateString(v *Validator, field, value string, required bool, minLen, maxLen int) {
	if required {
		v.Check(Required(value), field, fmt.Sprintf("%s is required", field))
	}

	if value != "" {
		if minLen > 0 {
			v.Check(MinLength(value, minLen), field, fmt.Sprintf("%s must be at least %d characters long", field, minLen))
		}
		if maxLen > 0 {
			v.Check(MaxLength(value, maxLen), field, fmt.Sprintf("%s must not be more than %d characters long", field, maxLen))
		}
	}
}
