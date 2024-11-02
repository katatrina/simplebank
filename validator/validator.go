package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

type FieldViolation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength {
		return fmt.Errorf("must contain from %d to %d characters", minLength, maxLength)
	}
	
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 25); err != nil {
		return err
	}
	
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only letters, numbers, or underscores")
	}
	
	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 100); err != nil {
		return err
	}
	
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 6, 200); err != nil {
		return err
	}
	
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 0, 100); err != nil {
		return err
	}
	
	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	
	return nil
}
