package entity

import (
	"errors"
	"strings"
)

func NewEmail(email string) (Email, error) {
	trimmed := strings.TrimSpace(email)
	if trimmed == "" {
		return "", errors.New("email cannot be empty")
	}
	return Email(trimmed), nil
}

func NewPassword(password string) (Password, error) {
	trimmed := strings.TrimSpace(password)
	if trimmed == "" {
		return "", errors.New("password cannot be empty")
	}
	return Password(trimmed), nil
}

func NewPhoneNumber(countryCode, number string) (PhoneNumber, error) {
	trimmedCountryCode := strings.TrimSpace(countryCode)
	trimmedNumber := strings.TrimSpace(number)
	if trimmedCountryCode == "" {
		return PhoneNumber{}, errors.New("country code cannot be empty")
	}
	if trimmedNumber == "" {
		return PhoneNumber{}, errors.New("phone number cannot be empty")
	}
	return PhoneNumber{
		CountryCode: trimmedCountryCode,
		Number:      trimmedNumber,
	}, nil
}
