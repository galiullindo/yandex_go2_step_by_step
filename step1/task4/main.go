package main

import "io"

func Copy(r io.Reader, w io.Writer, n uint) error {
	var ln int
	var err error
	buf := make([]byte, n)

	ln, err = r.Read(buf)
	if err != nil {
		if err == io.EOF {
		} else {
			return err
		}
	}

	if ln > 0 {
		ln, err = w.Write(buf[:ln])
	}
	if err != nil {
		return err
	}

	return nil
}
