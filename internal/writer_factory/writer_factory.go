package writer_factory

import (
	"io"
	"os"
)

// TODO: rename?
type WriterFactory interface {
	NewWriter(string) (io.WriteCloser, error)
}

type STDOUTWriter struct{}

func (sw *STDOUTWriter) NewWriter(_ string) (io.WriteCloser, error) {
	return &stdoutWriter{}, nil
}

// TODO: workaround to avoid closing stdout - that is already a workaround
// to make this interface work for stdout :/
type stdoutWriter struct{}

func (std *stdoutWriter) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

func (std *stdoutWriter) Close() error {

	return nil
}
