package main_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/test"
)

var (
	binName                   = "terradoc"
	inputFixtureName          = "golden-input.tfdoc.hcl"
	expectedOutputFixtureName = "golden-readme.md"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	err := build.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build %q: %v", binName, err)
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
	assert.NoError(t, err)

	// set up input file and content - TODO: structure this better
	inputContent := test.ReadFixture(t, inputFixtureName)

	inputFile, err := ioutil.TempFile("", "terradoc-input-")
	assert.NoError(t, err)

	defer inputFile.Close()

	_, err = inputFile.Write(inputContent)
	assert.NoError(t, err)

	cmdPath := filepath.Join(dir, binName)

	expectedOutput := test.ReadFixture(t, expectedOutputFixtureName)

	t.Run("ReadFromFile", func(t *testing.T) {
		cmd := exec.Command(cmdPath, inputFile.Name())

		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		if diff := cmp.Diff(output, expectedOutput); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})

	t.Run("ReadFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath)

		cmdStdIn, err := cmd.StdinPipe()
		assert.NoError(t, err)

		_, err = io.WriteString(cmdStdIn, string(inputContent))
		assert.NoError(t, err)

		cmdStdIn.Close()

		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		if diff := cmp.Diff(output, expectedOutput); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})

	t.Run("WriteToStdout", func(t *testing.T) {
		cmd := exec.Command(cmdPath, inputFile.Name())

		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		if diff := cmp.Diff(output, expectedOutput); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})

	t.Run("WriteToFile", func(t *testing.T) {
		f, err := ioutil.TempFile("", "terradoc-output-")
		assert.NoError(t, err)
		defer f.Close()

		cmd := exec.Command(cmdPath, "-o", f.Name(), inputFile.Name())

		err = cmd.Run()
		assert.NoError(t, err)

		result, err := ioutil.ReadAll(f)
		assert.NoError(t, err)

		if diff := cmp.Diff(result, expectedOutput); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})
}
