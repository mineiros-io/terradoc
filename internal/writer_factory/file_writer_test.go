package writer_factory_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mineiros-io/terradoc/internal/writer_factory"
)

func TestNewFileWriter(t *testing.T) {
	t.Run("DirectoryExists", func(t *testing.T) {
		dirName, err := ioutil.TempDir("", "test-terradoc*")
		if err != nil {
			t.Fatal(err)
		}

		defer os.RemoveAll(dirName)

		filename := "test-file"
		fpath := filepath.Join(dirName, filename)

		assertFileDoesNotExist(t, fpath)

		fileWriter := writer_factory.NewFileWriter(dirName)

		w, err := fileWriter.NewWriter(filename)
		defer w.Close()

		if err != nil {
			t.Fatal(err)
		}
		// w.Close()

		assertFileExists(t, fpath)
	})

	t.Run("DirectoryDoesNotExist", func(t *testing.T) {
		newDirName := "new-dir"
		defer os.RemoveAll(newDirName)

		_, err := os.Stat(newDirName)
		if !os.IsNotExist(err) {
			t.Fatal(err)
		}

		newFilename := "new-file"
		newFilepath := filepath.Join(newDirName, newFilename)

		assertFileDoesNotExist(t, newFilepath)

		fileWriter := writer_factory.NewFileWriter(newDirName)

		w, err := fileWriter.NewWriter(newFilename)
		if err != nil {
			t.Fatal(err)
		}
		defer w.Close()

		assertFileExists(t, newFilepath)
	})
}

func assertFileDoesNotExist(t *testing.T, filename string) {
	t.Helper()

	fi, err := os.Stat(filename)
	if fi != nil {
		t.Errorf("File %q should not exist", fi)
	}

	if err == nil {
		t.Error("Expected error but got none")
	}

	if !os.IsNotExist(err) {
		t.Errorf("Expected error %q. Got %q instead", os.ErrNotExist, err)
	}
}

func assertFileExists(t *testing.T, filename string) {
	t.Helper()

	_, err := os.Stat(filename)
	if err != nil {
		t.Errorf("Expected no error. Got %q", err)
	}
}
