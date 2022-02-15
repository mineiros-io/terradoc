package main_test

import (
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/test"
)

const (
	formatInput          = "format/unformatted-input.tfdoc.hcl"
	expectedFormatOutput = "format/formatted-input.tfdoc.hcl"
)

func TestFormat(t *testing.T) {
	unformattedInput := test.ReadFixture(t, "format/"+formatInput)
	// create another file with unformattedFile content to test overwrites

	expectedFormattedOutput := test.ReadFixture(t, expectedFormatOutput)

	t.Run("WriteToStdout", func(t *testing.T) {
		inputFile, err := ioutil.TempFile(t.TempDir(), "terradoc-fmt-output-")
		assert.NoError(t, err)
		// write unformatted input to file
		_, err = inputFile.Write(unformattedInput)
		assert.NoError(t, err)

		defer inputFile.Close()

		cmd := exec.Command(terradocBinPath, "fmt", inputFile.Name())

		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)

		if diff := cmp.Diff(expectedFormattedOutput, output); diff != "" {
			t.Errorf("Result is not expected (-want +got):\n%s", diff)
		}
	})

	t.Run("OverwriteFile", func(t *testing.T) {
		inputFile, err := ioutil.TempFile(t.TempDir(), "terradoc-fmt-output-")
		assert.NoError(t, err)
		// write unformatted input to file
		_, err = inputFile.Write(unformattedInput)
		assert.NoError(t, err)

		defer inputFile.Close()

		cmd := exec.Command(terradocBinPath, "fmt", "-w", inputFile.Name())

		_, err = cmd.CombinedOutput()
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
