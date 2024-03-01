package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")
	//if runtime.GOOS == "windows" {
	//	binName += ".exe"
	//}

	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)

		os.Exit(1)
	}
	fmt.Println("Running tests....")
	result := m.Run()
	fmt.Println("Cleaning up...")
	_ = os.Remove(binName)
	_ = os.Remove(fileName)
	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)
	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")

		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		_, _ = io.WriteString(cmdStdIn, task2)
		_ = cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("\n\t1: %s\n\t2: %s\n\n", task, task2)

		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

}
