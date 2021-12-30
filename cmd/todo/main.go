package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rjseymour66/todo"
)

// Hardcoding the file name
var todoFileName = ".todo.json"

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed for Command Line Applications in Golang\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2021\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		fmt.Fprintln(flag.CommandLine.Output(), "To add a todo, use the '-add' flag and supply")
		fmt.Fprintln(flag.CommandLine.Output(), "a task string to the command line, or pipe")
		fmt.Fprintln(flag.CommandLine.Output(), "a task to todo with the '-add' option")
		flag.PrintDefaults()
	}

	// assign flags. resulting vars are ptrs
	add := flag.Bool("add", false, "Add task to the ToDo list")
	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Int("delete", 0, "Delete task from the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	incomp := flag.Bool("incomp", false, "Lists incomplete tasks")

	flag.Parse()

	// Check if there is an ENV VAR for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// Assign l the address of a list because of the amount
	// of memory it takes up, and the methods all have a ptr
	// receiver
	l := &todo.List{}

	// Use the Get method to read ToDo items from file
	if err := l.Get(todoFileName); err != nil {
		// use STDERR for CLI tools
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the provided flags
	switch {
	case *list:
		// list current ToDo items
		fmt.Print(l)
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// When any arguments (excluding flags) are provided, they are
		// used as the new task. .Args() returns all non-flag input provided
		// by user as input
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Add the task
		l.Add(t)

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		// Delete the task by number
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *incomp:
		fmt.Print(l.Incomplete())
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// getTask decides whether to get the description for a new
// task from the arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
	// if there are args provided to the command, concatenate
	// them and return
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	// create a scanner and read a line (reads lines by default)
	s := bufio.NewScanner(r)
	s.Scan()
	// return if there is an error reading the line
	if err := s.Err(); err != nil {
		return "", err
	}

	// if the scanner didn't read anything, return error
	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}
