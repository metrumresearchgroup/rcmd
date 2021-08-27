package rcmd_test

import (
	"testing"

	"github.com/metrumresearchgroup/wrapt"

	"github.com/metrumresearchgroup/rcmd"
)

func TestRVersionExecution(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := rcmd.NewRSettings("")
	t.A.NoError(err)
	t.A.NotEmpty(rs.Version)

	// this test expects a machine with R 3.5.2 available on the default System Path
	t.R.Equal(rcmd.RVersion{4, 1, 0}, rs.Version)

	actual, _, _, err := rcmd.GetRVersionPlatformPath("")
	t.A.NoError(err)
	t.R.Equal(&rcmd.RVersion{4, 1, 0}, actual)
}
