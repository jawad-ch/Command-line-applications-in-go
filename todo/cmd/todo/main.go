package main

import (
	"flag"
	"fmt"
	"github.com/jawad-ch/Command-line-applications-in-go/todo"
	"os"
	"time"
)

// HardCoding the file name
const todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "%s tool. \nDeveloped for The Pragmatic Bookshelf\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Copyright %d\n", time.Now().Year())
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parsing command line flags
	task := flag.String("task", "", "Task to be included in the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	// Define an items list
	l := &todo.List{}

	// Use the Get method to read to do items from file
	if err := l.Get(todoFileName); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the provided flags
	switch {
	case *list:
		// List current to do items
		fmt.Print(l)
		//for _, item := range *l {
		//	if !item.Done {
		//		fmt.Println(item.Task)
		//	}
		//}
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
	case *task != "":
		// Add the task
		l.Add(*task)
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
