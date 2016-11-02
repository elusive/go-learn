package service

import (
	"net/http"

	"v.io/x/lib/vlog"
)

// ApiError Wraps an error with additional properties such as error code
type ApiError struct {
	// Internal full error
	Error error
	// User friendly error message that can be returned to clients
	Message string
	// Http error code to set in response
	HTTPErrorCode int
}

func badRequest(err error) *ApiError {
	vlog.Errorf("Encountered bad request error '%v'", err.Error())
	return &ApiError{err, "Invalid request body", http.StatusBadRequest}
}

func badRequestString(err error, msg string) *ApiError {
	vlog.Errorf("Encountered bad request error %v. Cause: %v", msg, err)
	return &ApiError{err, msg, http.StatusBadRequest}
}

func internalError(err error, msg string) *ApiError {
	vlog.Errorf("Encountered internal error %v. Cause: %v", msg, err)
	return &ApiError{err, msg, http.StatusInternalServerError}
}

func genericError(err error) *ApiError {
	vlog.Errorf("Encountered error %v", err)
	return &ApiError{err, "Internal server error", http.StatusInternalServerError}
}
