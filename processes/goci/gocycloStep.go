package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type gocycloStep struct {
	step
}

func newGocycloStep(name, exe, message, proj string, args []string) gocycloStep {
	s := gocycloStep{}
	s.step = newStep(name, exe, message, proj, args)
	return s
}

func (s gocycloStep) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.proj

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "Failed to execute",
			cause: err,
		}
	}

	if out.Len() > 0 {
		return "", &stepErr{
			step:  s.name,
			msg:   fmt.Sprintf("gocyclo err: %s", out.String()),
			cause: nil,
		}
	}

	return s.message, nil
}
