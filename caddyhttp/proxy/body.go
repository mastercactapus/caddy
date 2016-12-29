package proxy

import (
	"bytes"
	"io"
	"io/ioutil"
)

type rewindableReader interface {
	io.ReadCloser
	rewind() error
}

type bufferedBody struct {
	*bytes.Reader
}

var _ rewindableReader = &bufferedBody{}

func (*bufferedBody) Close() error {
	return nil
}

// rewind allows bufferedBody to be read again.
func (b *bufferedBody) rewind() error {
	if b == nil {
		return nil
	}
	_, err := b.Seek(0, io.SeekStart)
	return err
}

type unbufferedBody struct {
	io.ReadCloser
}

var _ rewindableReader = &unbufferedBody{}

func (b *unbufferedBody) rewind() error {
	panic("cannot rewind unbuffered body")
}

// newBufferedBody returns *bufferedBody to use in place of src. Closes src
// and returns Read error on src. All content from src is buffered.
func newBufferedBody(src io.ReadCloser) (rewindableReader, error) {
	if src == nil {
		return nil, nil
	}
	b, err := ioutil.ReadAll(src)
	src.Close()
	if err != nil {
		return nil, err
	}
	return &bufferedBody{
		Reader: bytes.NewReader(b),
	}, nil
}

// newUnbufferedBody returns *unbufferedBody as an interface
// to interchangably use it with newBufferedBody().
func newUnbufferedBody(src io.ReadCloser) rewindableReader {
	return &unbufferedBody{src}
}
