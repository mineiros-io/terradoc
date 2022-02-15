package main_test

import (
	"io"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/test"
)

const (
	generateInput          = "generate/golden-input.tfdoc.hcl"
	expectedGenerateOutput = "generate/golden-readme.md"
)

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
