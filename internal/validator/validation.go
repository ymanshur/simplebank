package validator

import (
	"fmt"
	"net/mail"
	"regexp"

	"github.com/ymanshur/simplebank/internal/common"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func validateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func validateRange(value int, min int, max int) error {
	if value < min || value > max {
		return fmt.Errorf("must be no less than %d and greater than %d", min, max)
	}
	return nil
}

func validateUsername(value string) error {
	if err := validateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

func validateFullName(value string) error {
	if err := validateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func validatePassword(value string) error {
	return validateString(value, 8, 100)
}

func validateEmail(value string) error {
	if err := validateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}

func validateCurrency(value string) error {
	if !common.IsSupportedCurrency(value) {
		return fmt.Errorf("%s is not supported currency", value)
	}
	return nil
}

func validateSecretCode(value string) error {
	return validateString(value, 32, 128)
}
