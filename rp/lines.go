package rp

import (
	"bytes"

	"github.com/metrumresearchgroup/filter"

	"github.com/metrumresearchgroup/rcmd/v2/filters"
)

// ScanLines scans lines from Rscript output and returns an array with
// the line numbers removed and whitespace trimmed.
func ScanLines(b []byte) []byte {
	return NewROutputFilter(false)(b)
}

// OutputOnly retrieves lines without > from interactive sessions.
func OutputOnly(b []byte) []byte {
	return NewROutputFilter(true)(b)
}

// NewROutputFilter scans lines from RScript output and returns an array with
// the line numbers removed, whitespace trimmed, and (optionally)
// with all input-like lines (which start with ">") excluded.
func NewROutputFilter(outputOnly bool) func([]byte) []byte {
	return func(b []byte) []byte {
		if !bytes.HasSuffix(b, []byte{'\n'}) {
			b = append(b, '\n')
		}

		var fns filter.Funcs
		if outputOnly {
			fns = append(fns, filters.DropInput)
		}
		fns = append(fns, filters.LineNumberStripper)

		return fns.Apply(b)
	}
}
