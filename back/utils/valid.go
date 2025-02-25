package utils

import (
	"errors"
	"strings"
)

const allowedDomain = "yourcompany.com"

func IsCorporateEmail(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}
	domain := parts[1]

	if domain != allowedDomain {
		return errors.New("only corporate email addresses are allowed")
	}

	return nil
}
