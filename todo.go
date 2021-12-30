package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

// List represents a list of ToDo items
type List []item

// Add creates a new ToDo item and appens it to the list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Complete marks a ToDo item as completed by setting
// Done = true and CompletedAt to the current time
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}
	// adjusting index for 0 based index
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

func (l *List) Incomplete() string {
	result := ""

	for k, task := range *l {
		prefix := "  "
		if !task.Done {
			result += fmt.Sprintf("%s%d: %s\n", prefix, k+1, task.Task)
		}
	}
	return result
}

func (l *List) Delete(i int) error {
	// switch to value for slice appending
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}
	// adjust index for 0 based index
	// use *l because you are changing the slice
	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save encodes the List as JSON and saves it using
// the provided file name
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	// ioutil is deprecated in 1.16. Use os.WriteFile
	// https://go.dev/doc/go1.16#ioutil
	return ioutil.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
	file, err := ioutil.ReadFile(filename)
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
func (l *List) String() string {
	formatted := ""

	for k, t := range *l {
		prefix := "  "
		if t.Done {
			prefix = "X "
		}
		// Adjust the item number k to print numbers starting from 1 instead of 0
		// Sprintf returns a formatted string
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
	}
	return formatted
}
