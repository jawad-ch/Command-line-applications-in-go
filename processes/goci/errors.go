package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("valdation failed")
	ErrSignal     = errors.New("received signal")
)

type stepErr struct {
	step  string
	msg   string
	cause error
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}
func (s *stepErr) Is(traget error) bool {
	t, ok := traget.(*stepErr)

	if !ok {
		return false
	}

	return t.step == s.step
}

func (s *stepErr) Unwarp() error {
	return s.cause
}
