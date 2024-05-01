package main

import (
	"context"
	"os/exec"
	"time"
)

type timeoutStep struct {
	step
	timeout time.Duration
}

var command = exec.CommandContext

func newTimeoutStep(name, exe, message, proj string, args []string, timeout time.Duration) timeoutStep {
	s := timeoutStep{}

	s.step = step{
		name:    name,
		exe:     exe,
		message: message,
		args:    args,
		proj:    proj,
	}

	s.timeout = timeout

	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}

	return s
}

func (s timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	cmd := command(ctx, s.exe, s.args...)
	cmd.Dir = s.proj
	if err := cmd.Run(); err != nil {

		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:  s.name,
				msg:   "failed to execute",
				cause: context.DeadlineExceeded,
			}
		}

		return "", &stepErr{
			step:  s.name,
			msg:   "Failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}
