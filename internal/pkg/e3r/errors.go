package e3r

import "net/http"

func BadRequest(message string) *Error {
	return New(message, http.StatusBadRequest)
}

func Forbidden(message string) *Error {
	return New(message, http.StatusForbidden)
}

func Internal(message string) *Error {
	return New(message, http.StatusInternalServerError)
}

func NotFound(message string) *Error {
	return New(message, http.StatusNotFound)
}

func TooManyRequests(message string) *Error {
	return New(message, http.StatusTooManyRequests)
}
