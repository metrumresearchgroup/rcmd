// +build !windows

package rcmd

import (
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestRVersionExecution(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("")
	t.A.NoError(err)

	// this test expects a machine with R 3.5.2 available on the default System Path
	expected := RVersion{3, 5, 2}
	t.A.Empty(rs.Version, "uncollected value")

	actual, err := rs.getRVersion()
	t.A.NoError(err)

	t.RunFatal("R version", func(t *wrapt.T) {
		t.A.Equal(expected, actual)
	})
	t.RunFatal("rs.Version", func(t *wrapt.T) {
		t.A.Equal(expected, rs.Version)
	})
}
