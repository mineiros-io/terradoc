package hclreader_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mineiros-io/terradoc/internal/srcreader/hclreader"
)

var content = `
root {
  section {
    title = "Test"
    description = "Test description"
  }
}
`

func TestRead(t *testing.T) {
	sf, err := hclreader.Read([]byte(content))
	if err != nil {
		t.Errorf("Expected no error. Got %q instead", err.Error())
	}

	if !bytes.Equal(sf.HCLFile.Bytes, []byte(content)) {
		t.Errorf("Expected content %q to be read from input. Got %q instead", content, string(sf.HCLFile.Bytes))
	}
}

func TestReadFromFile(t *testing.T) {
	t.Run("WhenFileExists", func(t *testing.T) {
		f, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal(err)
		}

		_, err = f.WriteString(content)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())

		sf, err := hclreader.ReadFromFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(sf.HCLFile.Bytes, []byte(content)) {
			t.Errorf("Expected content %q to be read from file. Got %q instead", content, string(sf.HCLFile.Bytes))
		}
	})

	t.Run("WhenFileDoesNotExist", func(t *testing.T) {
		sf, err := hclreader.ReadFromFile("random-file-name")
		if err == nil {
			t.Error("Expected an error trying to open a non existing file. Got none instead")
		}

		if sf != nil {
			t.Errorf("Expected no data to be read from non existing file. Got %+v instead", sf)
		}
	})
}
