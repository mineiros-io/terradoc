package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName   = "terradoc"
	inputFile = "testdata/input.tfdoc.hcl"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	err := build.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")

	result := m.Run()

	fmt.Println("Cleaning up...")

	os.Remove(binName)

	os.Exit(result)
}

func TestTerradocCLI(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("ReadFromFile", func(t *testing.T) {
		cmd := exec.Command(cmdPath, inputFile)

		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ReadFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath)

		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Logf("ERROR: %+v", err.Error())
			t.Fatal(err)
		}

		content := `
section {
	  title       = "Module Argument Reference"
	  description = "See [variables.tf] and [examples/] for details and use-cases."

	  section {
	    title = "Main Resource Configuration"

	    variable "local_secondary_indexes" {
	      type        = any
	      readme_type = "list(local_secondary_index)"

	      description = "Describe an LSI on the table; these can only be allocated creation so you cannot change this definition after you have created the resource."
	      default     = []

	      required = true

	      forces_recreation = true

	      readme_example = {
	        local_secondary_indexes = [
	          {
	            range_key = "someKey"
	          }
	        ]
	      }

	      attribute "range_key" {
	        type = string

	        description = "The attribute to use as the range (sort) key. Must also be defined as an attribute, see below."

	        forces_recreation = true
	      }
	    }
	  }

	}
`

		io.WriteString(cmdStdIn, content)
		cmdStdIn.Close()

		err = cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("OutputToSTDOUT", func(t *testing.T) {
		cmd := exec.Command(cmdPath, inputFile)

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n", "foo")
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead.\n", expected, string(out))
		}
	})

	t.Run("OutputToDirectory", func(t *testing.T) {
	})
}
