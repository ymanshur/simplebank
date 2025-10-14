package validator

import validation "github.com/go-ozzo/ozzo-validation"

func ValidUsername(value any) error {
	s, ok := value.(string)
	if ok && s != "" {
		return validateUsername(s)
	}
	return nil
}

func ValidFullName(value any) error {
	s, ok := value.(string)
	if ok && s != "" {
		return validateFullName(s)
	}
	return nil
}

func ValidPassword(value any) error {
	s, ok := value.(string)
	if ok && s != "" {
		return validatePassword(s)
	}
	return nil
}

func ValidEmail(value any) error {
	s, ok := value.(string)
	if ok && s != "" {
		return validateEmail(s)
	}
	return nil
}

func ValidCurrency(value any) error {
	s, ok := value.(string)
	if ok && s != "" {
		return validateCurrency(s)
	}
	return nil
}

func ValidPageSize(value any) error {
	n, ok := value.(int)
	if ok && n != 0 {
		return validateRange(n, 5, 10)
	}
	return nil
}

func ValidSecretCode(value any) error {
	s, ok := value.(string)
	if ok && s != "" {
		return validateSecretCode(s)
	}
	return nil
}

func ValidID() validation.Rule {
	return validation.Min(1)
}
