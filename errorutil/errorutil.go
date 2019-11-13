package errorutil

import (
	"errors"
	"fmt"
)

type statusError struct {
	message string
	code    int
}

func (s *statusError) Error() string {
	return fmt.Sprintf("%d: %s", s.code, s.message)
}

func New(code int, msg string) error {
	return &statusError{
		message: msg,
		code:    code,
	}
}

func Newf(code int, msg string, a ...interface{}) error {
	return New(code, fmt.Sprintf(msg, a...))
}

func FromError(code int, err error) error {
	return &statusError{
		message: err.Error(),
		code:    code,
	}
}

func Code(err error) *int {
	s, ok := err.(*statusError)
	if !ok {
		return nil
	}
	return &s.code
}

func Annotate(err error, msg string) error {
	newMsg := fmt.Sprintf("%s: %s", msg, err.Error())
	s, ok := err.(*statusError)
	if !ok {
		return errors.New(newMsg)
	}
	s.message = newMsg
	return s
}

func Annotatef(err error, msg string, a ...interface{}) error {
	return Annotate(err, fmt.Sprintf(msg, a...))
}
