package validator

import (
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Custom validation rules
var (
	phoneRegex = regexp.MustCompile(`^[+]?[0-9]{10,15}$`)
	slugRegex  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

// Init initializes custom validators
func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators
		v.RegisterValidation("strong_password", strongPassword)
		v.RegisterValidation("phone", phoneNumber)
		v.RegisterValidation("slug", slug)
		v.RegisterValidation("no_sql_injection", noSQLInjection)

		// Register custom tag name func
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// strongPassword validates password strength
func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasMinLen  = len(password) >= 8
		hasMaxLen  = len(password) <= 50
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasMaxLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// phoneNumber validates phone numbers
func phoneNumber(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return phoneRegex.MatchString(phone)
}

// slug validates URL slugs
func slug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	return slugRegex.MatchString(slug)
}

// noSQLInjection checks for common SQL injection patterns
func noSQLInjection(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())

	// Common SQL injection patterns
	dangerous := []string{
		"select",
		"insert",
		"update",
		"delete",
		"drop",
		"union",
		"exec",
		"execute",
		"script",
		"javascript:",
		"<script",
		"onclick",
		"onerror",
		"--",
		"/*",
		"*/",
		"xp_",
		"sp_",
		"0x",
		"\\x",
	}

	for _, pattern := range dangerous {
		if strings.Contains(value, pattern) {
			return false
		}
	}

	return true
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors formats validation errors
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			var message string

			switch e.Tag() {
			case "required":
				message = "This field is required"
			case "email":
				message = "Invalid email format"
			case "min":
				message = "Too short"
			case "max":
				message = "Too long"
			case "strong_password":
				message = "Password must be 8-50 characters with uppercase, lowercase, number and special character"
			case "phone":
				message = "Invalid phone number format"
			case "slug":
				message = "Invalid slug format (lowercase letters, numbers and hyphens only)"
			case "no_sql_injection":
				message = "Input contains potentially dangerous characters"
			default:
				message = "Invalid value"
			}

			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: message,
			})
		}
	}

	return errors
}