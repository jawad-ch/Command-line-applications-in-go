package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// item struct represents a ToDo item
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

// Add creates a new todo item and appends it to the list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Complete method marks a ToDo item as completed by
//setting Done = true and CompletedAt to the current time

func (l *List) Complete(i int) error {
	ls := *l

	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	// Adjusting index for 0 based index
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete method deletes a ToDo item from the list
func (l *List) Delete(i int) error {
	ls := *l

	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	// Adjusting index for 0 based index
	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save method encodes the List as JSON and saves it
// using the provided file name
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

// Get method opens the provided file name, decodes
// the JSON data and parses it into a List
func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}

// String prints out a formatted list
// Implements the fmt.Stringer interface
func (l *List) String(limit int) string {
	formatted := ""
	ls := *l
	if limit >= 0 {
		if limit > len(ls) {
			limit = len(ls)
		}
		ls = ls[:limit]
	}
	for k, t := range ls {
		prefix := "\t"
		createdAt := t.getDate(false)
		completedAt := ""
		if t.Done {
			prefix = "X\t"
			completedAt = fmt.Sprintf("[completed : %s]", t.getDate(true))
		}

		formatted += fmt.Sprintf("%s%d: %s [created : %s]%s\n", prefix, k+1, t.Task, createdAt, completedAt)
	}
	return formatted
}

func (task *item) getDate(completeTime bool) string {
	d := task.CreatedAt.String()
	if completeTime {
		d = task.CompletedAt.String()
	}
	parsedDate, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 -07", d)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "error parsing date"
	}

	// Calculate the time difference
	duration := time.Since(parsedDate)

	// Convert the duration to minutes
	minutes := int(duration.Minutes())

	// Print the human-readable format
	if minutes < 1 {
		return fmt.Sprintf("Just now")
	} else {
		return fmt.Sprintf("%d min ago", minutes)
	}
}
