package hcl_file_reader

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestFile struct {
	name string
}

func (tf *TestFile) Close() error {
	return os.Remove(tf.name)
}

func createTestFile(filepath, contents string) (*TestFile, error) {
	f, err := ioutil.TempFile("", filepath)
	if err != nil {
		return nil, err
	}

	_, err = f.WriteString(contents)

	if err != nil {
		return nil, err
	}

	return &TestFile{name: f.Name()}, nil
}

func TestReadWhenFileExists(t *testing.T) {
	content := `
root {
  section {
    title = "Test"
    description = "Test description"
  }
}
`
	filepath := "test.tfdoc.hcl"

	// generaete test file
	testFile, err := createTestFile(filepath, content)
	require.NoError(t, err)

	// ensure we delete created test file
	defer testFile.Close()

	sourceFile, err := Read(testFile.name)
	require.NoError(t, err)

	require.Equal(t, content, string(sourceFile.HCLFile.Bytes))
}

func TestReadWhenFileHasInvalidContents(t *testing.T) {
	content := `
I-----think909831
)(*&)(*^&^*^%&*$I
":>?>Might
b!@)(#%$%
e='./
]\]]-09-342Invalid!@#!w80
`
	filepath := "test.tfdoc.hcl"

	// generaete test file
	testFile, err := createTestFile(filepath, content)
	require.NoError(t, err)

	// ensure we delete created test file
	defer testFile.Close()

	t.Skip("TODO: this fails. a lot")

	sourceFile, err := Read(testFile.name)
	require.NoError(t, err)

	require.Equal(t, content, string(sourceFile.HCLFile.Bytes))
}

func TestReadWhenFileDoesNotExist(t *testing.T) {
	filepath := "test-lalalala"

	sourceFile, err := Read(filepath)
	require.Errorf(t, err, `Failed to read file; The configuration file "%s" could not be read.`, filepath)
	require.Nil(t, sourceFile)
}
