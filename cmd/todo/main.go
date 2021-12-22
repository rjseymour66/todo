package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rjseymour66/todo"
)

// Hardcoding the file name
const todoFileName = ".todo.json"

func main() {
	// Assign l the address of a list because it is
	// a lot to store in memory during operations.

	// Also, the methods all have a ptr receiver.
	l := &todo.List{}

	// Use the Get method to read ToDo items from file
	if err := l.Get(todoFileName); err != nil {
		// use STDERR for CLI tools
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the number of arguments provided
	switch {
	// If not Args, print the list
	case len(os.Args) == 1:
		// List current ToDo items
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	// Concatenate all provided args with a space and add
	// to the list as an item
	default:
		// Concatenate all args with a space
		item := strings.Join(os.Args[1:], " ")

		// Add the task
		l.Add(item)

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
