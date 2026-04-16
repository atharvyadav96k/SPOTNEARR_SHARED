package entity

import (
	"errors"
	"fmt"
	"net/mail"
	"reflect"
	"strings"
	"unicode"
)

func ValidateStruct(s interface{}) []error {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return []error{errors.New("input must be a struct")}
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name
		errList := []error{}
		if validator, ok := field.Addr().Interface().(Validator); ok {
			if err := validator.Validate(); err != nil {
				errList = append(errList, err)
				return []error{fmt.Errorf("validation failed for field: %s", fieldName)}
			}
		}
	}

	return nil
}

func (e *Email) String() string {
	return string(*e)
}

func (e *Email) Validate() error {
	emailStr := strings.TrimSpace(string(*e))
	if emailStr == "" {
		return errors.New("email address cannot be empty")
	}

	addr, err := mail.ParseAddress(emailStr)
	if err != nil {
		return errors.New("invalid email format")
	}

	if addr.Address != emailStr {
		return errors.New("email contains invalid characters or extra spaces")
	}

	parts := strings.Split(emailStr, "@")
	domain := parts[1]
	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return errors.New("email domain is invalid")
	}

	return nil
}

func (p *Password) String() string {
	return string(*p)
}

func (p *Password) Validate() error {
	pStr := string(*p)
	if len(pStr) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(pStr) > 72 {
		return errors.New("password is too long (maximum 72 characters)")
	}

	var hasUpper, hasLower, hasNumber bool
	for _, char := range pStr {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

func (pn *PhoneNumber) String() string {
	return pn.CountryCode + pn.Number
}

func (pn *PhoneNumber) Validate() error {
	if !strings.HasPrefix(pn.CountryCode, "+") || len(pn.CountryCode) < 2 {
		return errors.New("invalid country code (must start with +)")
	}

	cleanNumber := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, pn.Number)

	if len(cleanNumber) < 6 {
		return errors.New("phone number is too short")
	}
	if len(cleanNumber) > 15 {
		return errors.New("phone number is too long")
	}

	return nil
}
