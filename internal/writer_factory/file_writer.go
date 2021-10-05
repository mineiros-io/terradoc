package writer_factory

import (
	"io"
	"os"
	"path/filepath"
)

type FileWriter struct {
	DirPath string
}

func NewFileWriter(dirPath string) *FileWriter {
	return &FileWriter{DirPath: dirPath}
}

func (fw *FileWriter) NewWriter(fileName string) (io.WriteCloser, error) {
	// TODO: refactor
	_, err := os.Stat(fw.DirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		err = os.MkdirAll(fw.DirPath, os.ModePerm)

		if err != nil {
			return nil, err
		}
	}

	f, err := os.Create(filepath.Join(fw.DirPath, fileName))
	if err != nil {
		return nil, err
	}

	return f, nil
}
