package provider

import (
	"io"
	"os"
)

type LocalProvider struct{}

func (l *LocalProvider) Fetch(path string) (io.ReadCloser, error) {
	res, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return res, nil
}
