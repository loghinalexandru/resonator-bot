package audio

import "io"

type Provider interface {
	Audio(path string) (io.ReadCloser, error)
}
