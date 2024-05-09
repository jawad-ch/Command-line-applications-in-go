package pomodoro_test

import (
	"testing"

	"github.com/jawad-ch/Command-line-application-in-go/interactiveTools/pomo/pomodoro"
	"github.com/jawad-ch/Command-line-application-in-go/interactiveTools/pomo/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}
