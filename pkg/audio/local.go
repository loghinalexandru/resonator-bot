package audio

import (
	"io"
	"os"
)

type localProvider struct{}

func NewLocal() *localProvider {
	return &localProvider{}
}

func (l localProvider) Audio(path string) (io.ReadCloser, error) {
	res, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return res, nil
}
