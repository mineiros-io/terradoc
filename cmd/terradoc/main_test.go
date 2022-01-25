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
	binName                = "terradoc"
	generateInput          = "golden-input.tfdoc.hcl"
	formatInput            = "unformatted-input.tfdoc.hcl"
	expectedGenerateOutput = "golden-readme.md"
	expectedFormatOutput   = "formatted-input.tfdoc.hcl"
)

var dir, terradocBinPath string

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	var err error

	dir, err = os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get current directory: %v", err)
		os.Exit(1)
	}
	terradocBinPath = filepath.Join(dir, binName)

	build := exec.Command("go", "build", "-o", binName)

	err = build.Run()
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

func TestGenerate(t *testing.T) {
	inputContent := test.ReadFixture(t, generateInput)
	// create tempfile
	inputFile, err := ioutil.TempFile("", "terradoc-input-")
	assert.NoError(t, err)
	// write content to tempfile
	_, err = inputFile.Write(inputContent)
	assert.NoError(t, err)

	defer inputFile.Close()

	expectedOutput := test.ReadFixture(t, expectedGenerateOutput)

	t.Run("ReadFromFile", func(t *testing.T) {
		cmd := exec.Command(terradocBinPath, "generate", inputFile.Name())

		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		if diff := cmp.Diff(expectedOutput, output); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})

	t.Run("ReadFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(terradocBinPath, "generate", "-")

		cmdStdIn, err := cmd.StdinPipe()
		assert.NoError(t, err)

		_, err = io.WriteString(cmdStdIn, string(inputContent))
		assert.NoError(t, err)

		cmdStdIn.Close()

		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		if diff := cmp.Diff(expectedOutput, output); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})

	t.Run("WriteToStdout", func(t *testing.T) {
		cmd := exec.Command(terradocBinPath, "generate", inputFile.Name())

		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		if diff := cmp.Diff(expectedOutput, output); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})

	t.Run("WriteToFile", func(t *testing.T) {
		f, err := ioutil.TempFile("", "terradoc-output-")
		assert.NoError(t, err)
		defer f.Close()

		cmd := exec.Command(terradocBinPath, "generate", "-o", f.Name(), inputFile.Name())

		err = cmd.Run()
		assert.NoError(t, err)

		result, err := ioutil.ReadAll(f)
		assert.NoError(t, err)

		if diff := cmp.Diff(expectedOutput, result); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})
}

func TestFormat(t *testing.T) {
	unformattedInput := test.ReadFixture(t, formatInput)
	// create another file with unformattedFile content to test overwrites
	inputFile, err := ioutil.TempFile("", "terradoc-fmt-output-")
	assert.NoError(t, err)
	// write unformatted input to file
	_, err = inputFile.Write(unformattedInput)
	assert.NoError(t, err)

	defer inputFile.Close()

	expectedFormattedOutput := test.ReadFixture(t, expectedFormatOutput)

	t.Run("WriteToStdout", func(t *testing.T) {
		cmd := exec.Command(terradocBinPath, "fmt", inputFile.Name())

		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		if diff := cmp.Diff(expectedFormattedOutput, output); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})

	t.Run("OverwriteFile", func(t *testing.T) {
		cmd := exec.Command(terradocBinPath, "fmt", "-w", inputFile.Name())

		_, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		_, err = inputFile.Seek(0, 0)
		assert.NoError(t, err)

		persistedResult, err := ioutil.ReadAll(inputFile)
		assert.NoError(t, err)

		if diff := cmp.Diff(expectedFormattedOutput, persistedResult); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})
}
