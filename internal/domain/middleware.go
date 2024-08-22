package domain

import "io"

type Middleware interface {
	Write(w io.Writer) (wp io.Writer, err error)
	Read(r io.Reader) (rp io.Reader, err error)
}

func ChainMiddlewareRead(m []Middleware) func(io.Reader) (io.Reader, error) {
	return func(r io.Reader) (io.Reader, error) {
		reader := r

		var err error

		for _, middleware := range m {
			reader, err = middleware.Read(reader)
			if err != nil {
				return nil, err
			}
		}

		return reader, nil
	}
}

func ChainMiddlewareWrite(m []Middleware) func(io.Writer) (io.Writer, error) {
	return func(r io.Writer) (io.Writer, error) {
		writer := r

		var err error

		for _, middleware := range m {
			writer, err = middleware.Write(writer)
			if err != nil {
				return nil, err
			}
		}

		return writer, nil
	}
}
