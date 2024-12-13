package e3r

import (
	"forum/internal/pkg/httphelper"
	"log"
	"net/http"
)

type Error struct {
	Message string
	Code    int
}

func New(message string, code int) *Error {
	return &Error{Message: message, Code: code}
}

func Wrap(err error, code int) *Error {
	return &Error{Message: err.Error(), Code: code}
}

func GetCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if e, ok := err.(*Error); ok {
		return e.Status()
	}

	return http.StatusInternalServerError
}

func (e Error) Status() int {
	return e.Code
}

func (e Error) Error() string {
	return e.Message
}

func ErrorEncoder(err error, w http.ResponseWriter, authN bool) {
	status := getStatus(err)
	log.Printf("Response status: %d, Message: %s", status, err.Error())

	// if bad request return warning with prefilled data

	httphelper.Render(w, status, "error", httphelper.GetTmplData(Error{
		Message: err.Error(),
		Code:    status,
	}, authN))
}

func getStatus(err error) int {
	if e, ok := err.(*Error); ok {
		return e.Status()
	}

	return http.StatusInternalServerError
}
