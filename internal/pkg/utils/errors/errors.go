package errors

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e CustomError) Error() string {
	return fmt.Sprintf("status: %d, message: %s", e.Status, e.Message)
}

func BadRequest(message string) CustomError {
	return CustomError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func Unauthorized(message string) CustomError {
	return CustomError{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func NotFound(message string) CustomError {
	return CustomError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func Internal(message string) CustomError {
	return CustomError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}
