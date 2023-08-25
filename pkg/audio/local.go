package audio

import (
	"io"
	"os"
)

type Local struct{}

func (l Local) Audio(path string) (io.ReadCloser, error) {
	res, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return res, nil
}
