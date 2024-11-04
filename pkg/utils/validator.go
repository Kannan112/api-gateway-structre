package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate

	// Common regex patterns
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};:'",.<>/?]{8,}$`)
	phoneRegex    = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,16}$`)
)

type ValidationError struct {
	Field   string
	Message string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

// InitValidator initializes the validator with custom validations
func InitValidator() {
	validate = validator.New()

	// Register custom validation tags
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("username", validateUsername)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}

// Validate struct based on tags
func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError
	err := GetValidator().Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Message: getErrorMsg(err),
			})
		}
	}
	return errors
}

// Common validation functions
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return &ValidationError{
			Field:   "email",
			Message: "invalid email format",
		}
	}
	return nil
}

func ValidatePassword(password string) error {
	if !passwordRegex.MatchString(password) {
		return &ValidationError{
			Field:   "password",
			Message: "password must be at least 8 characters and contain valid characters",
		}
	}

	var (
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

	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		return &ValidationError{
			Field:   "password",
			Message: "password must contain at least one uppercase letter, lowercase letter, number, and special character",
		}
	}

	return nil
}

// Custom validator functions
func validatePassword(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

func validateUsername(fl validator.FieldLevel) bool {
	return usernameRegex.MatchString(fl.Field().String())
}

// Helper function to get error messages
func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "password":
		return "invalid password format"
	case "phone":
		return "invalid phone number format"
	case "username":
		return "username must be 3-16 characters long and contain only letters, numbers, underscores, or hyphens"
	default:
		return fmt.Sprintf("validation failed on '%s' tag", err.Tag())
	}
}

// Example struct with validation tags
type User struct {
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Phone    string `json:"phone" validate:"required,phone"`
}
