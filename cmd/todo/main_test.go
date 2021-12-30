package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

var (
	binName = "todo"
	// fileName = ".todo.json"
)

// Integration tests that build the tool, run the tests, and clean up.
func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	// Build a windows .exe if Windows machine
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	fileName := os.Getenv("TODO_FILENAME")

	if os.Getenv("TODO_FILENAME") == "" {
		fmt.Fprintln(os.Stderr, "Cannot locate TODO_FILENAME environment variable")
		os.Exit(1)
	}

	// build is the "go build -o binName" command
	build := exec.Command("go", "build", "-o", binName)

	// run build command and wait for it to complete
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	// run the tests and returns an exit code to result
	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {

	// dir stores a rooted path name
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// create a path with dir and the executable from main()
	cmdPath := filepath.Join(dir, binName)

	// t.Run creates subtests that depend on each other

	// Add task from args
	task := "test task number 1"
	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		// add first task
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// Add task from STDIN
	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		// add task2
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// List tasks
	t.Run("ListTasks", func(t *testing.T) {
		// Command returns a Cmd struct to execute the named program
		// with the given args
		cmd := exec.Command(cmdPath, "-list")
		// out stores the combined STDOUT and STDERR output
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
	// List incomplete tasks
	incTask := "Incomplete task"
	t.Run("ListIncompleteTasks", func(t *testing.T) {
		// Build add completed command
		cmd := exec.Command(cmdPath, "-add", incTask)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		// Build complete command
		compCmd := exec.Command(cmdPath, "-complete", strconv.Itoa(3))

		if err := compCmd.Run(); err != nil {
			t.Fatal(err)
		}

		// Build the incomplete command
		cmdInc := exec.Command(cmdPath, "-incomp")

		// out stores STDOUT and STDERR output
		out, err := cmdInc.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	// Complete tasks
	t.Run("ListCompleteTasks", func(t *testing.T) {
		// Build -list command that includes complete tasks
		cmd := exec.Command(cmdPath, "-list")

		// out stores STDOUT and STDERR output
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n  2: %s\nX 3: %s\n", task, task2, incTask)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	// Delete tasks
	delTask := strconv.Itoa(1)
	t.Run("DeleteTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-delete", delTask)

		// If -delete returns a nil, there is an error
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
}
