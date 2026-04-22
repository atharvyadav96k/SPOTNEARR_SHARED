package neon

import (
	"errors"

	"gorm.io/gorm"
)

// IsRecordNotFound checks if the error is a "record not found" error.
// This is useful when you want to distinguish between "not found" and other errors.
//
// Example:
//
//	var user User
//	err := neonService.FindByID(ctx, &user, userID)
//	if neon.IsRecordNotFound(err) {
//	    // Handle not found case
//	    return fmt.Errorf("user not found")
//	}
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsInvalidData checks if the error indicates invalid data or uninitialized service.
func IsInvalidData(err error) bool {
	return errors.Is(err, gorm.ErrInvalidData)
}

// HandleQueryError provides consistent error handling for database queries.
// It returns a user-friendly error message based on the actual error type.
//
// Example:
//
//	err := neonService.FindByID(ctx, &user, userID)
//	if err != nil {
//	    userMsg := neon.HandleQueryError(err, "user", "123")
//	    // userMsg could be "user with ID 123 not found" or "database query failed"
//	}
func HandleQueryError(err error, entityType string, identifier string) string {
	if err == nil {
		return ""
	}

	if IsRecordNotFound(err) {
		return entityType + " with ID " + identifier + " not found"
	}

	if IsInvalidData(err) {
		return "database service not initialized"
	}

	return "failed to query " + entityType
}
