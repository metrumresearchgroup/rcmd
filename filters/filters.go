package filters

import (
	"bytes"
	"regexp"
)

func NewPrefixFilter(pfx string) func([]byte) []byte {
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

func DropInput(bs []byte) []byte {
	if bytes.HasPrefix(bs, []byte{'>'}) {
		return nil
	}

	return bs
}

// LineNumberRegex applies matches any [n] line prefix.
var LineNumberRegex = regexp.MustCompile(`^\s*\[\d+]\s*`)

// LineNumberStripper removes the line markers matching LineNumberRegex.
// R uses this format to indicate output lines. It is not useful to have
// if you're trying to use the data in the line, hence LineNumberStripper.
func LineNumberStripper(bs []byte) []byte {
	buf := bytes.Buffer{}

	if loc := LineNumberRegex.FindIndex(bs); loc != nil {
		buf.Write(bs[loc[1]:])
	} else {
		buf.Write(bs)
	}

	return bytes.TrimSpace(buf.Bytes())
}
