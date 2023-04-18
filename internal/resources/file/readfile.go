package file

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

// ReadLines If src != nil, ReadFile converts src to a []byte if possible;
// otherwise it returns an error. If src == nil, readSource returns
// the result of reading the file specified by filename.
func ReadLines(filename string, src any, lines int) ([]byte, error) {
	var err error
	if src == nil {
		src, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
	}
	switch t := src.(type) {
	case string:
		return str2Lines(t, lines)
	case []byte:
		return str2Lines(string(t), lines)
	case *bytes.Buffer:
		return buf2Lines(t, lines)
	case io.Reader:
		return rdr2Lines(t, lines)
	}
	return nil, errors.New("invalid source")
}

func str2Lines(str string, lines int) ([]byte, error) {
	c := 0
	for i := 0; i < len(str); i++ {
		if str[i] == '\n' {
			c++
		}
		if c == lines {
			return []byte(str[:i+1]), nil
		}
	}
	return []byte(str), nil
}

func buf2Lines(buf *bytes.Buffer, lines int) ([]byte, error) {
	var out []rune
	c := 0
	for {
		r, _, err := buf.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		out = append(out, r)
		if r == '\n' {
			c++
		}
		if c == lines {
			return []byte(string(out)), nil
		}
	}
	if len(out) > 0 {
		return []byte(string(out)), nil
	}
	return buf.Bytes(), nil
}

func rdr2Lines(rdr io.Reader, lines int) ([]byte, error) {
	buf := bufio.NewReader(rdr)
	var out []rune
	c := 0
	for {
		r, _, err := buf.ReadRune()
		if err == io.EOF {
			return []byte(string(out)), nil
		} else if err != nil {
			return nil, err
		}
		out = append(out, r)
		if r == '\n' {
			c++
		}
		if c == lines {
			return []byte(string(out)), nil
		}
	}
}

// ReadFile If src != nil, ReadFile converts src to a []byte if possible;
// otherwise it returns an error. If src == nil, readSource returns
// the result of reading the file specified by filename.
func ReadFile(filename string, src any) ([]byte, error) {
	var err error
	if src == nil {
		src, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
	}
	switch s := src.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	case *bytes.Buffer:
		// is io.Reader, but src is already available in []byte form
		if s != nil {
			return s.Bytes(), nil
		}
	case io.Reader:
		return io.ReadAll(s)
	}
	return nil, errors.New("invalid source")
}
