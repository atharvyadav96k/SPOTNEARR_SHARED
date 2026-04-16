package res

import (
	"net/http"

	common_error "github.com/atharvyadav96k/gcp/common/error"
)

func Success(w http.ResponseWriter, message string, data interface{}) {
	Send(w, http.StatusOK, message, data)
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	Send(w, http.StatusCreated, message, data)
}

func NoContent(w http.ResponseWriter) {
	Send(w, http.StatusNoContent, "No Content", nil)
}

func InternalServerError(w http.ResponseWriter, err []error) {
	var errorMessages []string
	errorMessages = common_error.ErrorsToString(err)
	Send(w, http.StatusInternalServerError, "Internal Server Error", errorMessages)
}

func BadRequest(w http.ResponseWriter, err []error) {
	var errorMessages []string
	errorMessages = common_error.ErrorsToString(err)
	Send(w, http.StatusBadRequest, "Bad Request", errorMessages)
}

func NotFound(w http.ResponseWriter, err []error) {
	var errorMessages []string
	errorMessages = common_error.ErrorsToString(err)
	Send(w, http.StatusNotFound, "Not Found", errorMessages)
}

func Forbidden(w http.ResponseWriter, err []error) {
	var errorMessages []string
	errorMessages = common_error.ErrorsToString(err)
	Send(w, http.StatusForbidden, "Forbidden", errorMessages)
}

func Unauthorized(w http.ResponseWriter, err []error) {
	var errorMessages []string
	errorMessages = common_error.ErrorsToString(err)
	Send(w, http.StatusUnauthorized, "Unauthorized", errorMessages)
}
