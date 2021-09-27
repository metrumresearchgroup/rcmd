package rp

import (
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
		var c filter.Chain
		if outputOnly {
			c = append(c, filters.DropInput)
		}
		c = append(c, filters.LineNumberStripper)

		f := filter.NewFlow(c)

		return f.Apply(b)
	}
}
