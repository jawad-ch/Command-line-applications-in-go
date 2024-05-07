package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jawad-ch/Command-line-applications-in-go/todo"
)

// Default file name
var todoFileName = ".todo.json"

func main() {

	// Check if the user defined the ENV VAR for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "%s tool. \nDeveloped for The Pragmatic Bookshelf\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Copyright %d\n", time.Now().Year())
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Item to be deleted from toDo list")
	limit := flag.Int("limit", 0, "Max Items to show")

	flag.Parse()

	// Define an items list
	l := &todo.List{}

	// Use the Get method to read to do items from file
	if err := l.Get(todoFileName); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%w : hnnn", err)
		os.Exit(1)
	}

	// Decide what to do based on the provided flags
	switch {
	case *list:
		// List current to do items
		if *limit > 0 {
			fmt.Println(l.String(*limit))
		} else {
			fmt.Println(l.String(len(*l)))
			//fmt.Print(l)
		}

	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *del > 0:
		// Complete the given item
		if err := l.Delete(*del); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		// When any arguments (excluding flags) are provided, they will be
		//used as the new task
		t, err := getTask(os.Stdin, flag.Args()...)

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Add the task
		l.Add(t)
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	// Concatenate all provided arguments
	default:
		// Invalid flag provided
		_, _ = fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)

	s.Scan()

	if err := s.Err(); err != nil {
		return "", nil
	}
	text := s.Text()
	if len(text) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return text, nil
}
