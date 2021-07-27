package writers

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"sync"
)

type Filter struct {
	w   io.Writer
	ffs []FilterFunc
	m   *sync.Mutex
}

type FilterFunc func([]byte) []byte


func NewFilter(w io.Writer, filters ...FilterFunc) *Filter {
	f := &Filter{
		w:   w,
		ffs: filters,
		m:   &sync.Mutex{},
	}

	return f
}

func (f *Filter) Write(p []byte) (n int, err error) {
	f.m.Lock()
	defer f.m.Unlock()

	if len(f.ffs) == 0 {
		return f.w.Write(p)
	}

	return f.filterLines(p)
}

var LineRegEx = regexp.MustCompile("^\\s*\\[\\d+]\\s*")

func LineNumberStripper(bs []byte) []byte {
	buf := bytes.Buffer{}

	loc := LineRegEx.FindIndex(bs)
	if loc != nil {
		buf.Write(bs[loc[1]:])
	} else {
		buf.Write(bs)
	}

	return buf.Bytes()
}

func NewPrefixFilter(pfx string) FilterFunc {
	if len(pfx) == 0 {
		return func(bs []byte) []byte {
			return bs
		}
	}

	return func(bs []byte) []byte {
		buf := bytes.Buffer{}

		if len(pfx) > 0 {
			buf.WriteString(pfx)
			buf.WriteByte(' ')
		}
		buf.Write(bs)

		return buf.Bytes()
	}
}

func InputStripper(bs []byte) []byte {
	buf := bytes.Buffer{}

	if !bytes.HasPrefix(bs, []byte{'>'}) {
		buf.Write(bs)
	}

	return buf.Bytes()
}


func (f *Filter) filterLines(p []byte) (n int, err error) {
	buf := bytes.NewReader(p)
	s := bufio.NewScanner(buf)

	var written int

scan:
	for s.Scan() {
		bs := s.Bytes()

		for _, lw := range f.ffs {
			bs = lw(bs)
			if len(bs) == 0 {
				continue scan
			}
		}

		for _, v := range [][]byte{
			bs,
			{'\n'},
		} {
			written, err = f.w.Write(v)
			if err != nil {
				return n, err
			}
			n += written
		}
	}

	return n, nil
}
