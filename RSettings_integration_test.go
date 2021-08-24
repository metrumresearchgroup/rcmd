package rcmd

import (
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestRVersionExecution(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("")
	t.A.NoError(err)
	t.A.NotEmpty(rs.Version)

	// this test expects a machine with R 3.5.2 available on the default System Path
	t.R.Equal(RVersion{4, 1, 0}, rs.Version)

	actual, err := getRVersion("")
	t.A.NoError(err)
	t.R.Equal(&RVersion{4, 1, 0}, actual)
}
