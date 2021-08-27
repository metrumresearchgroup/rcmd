package rcmd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
	"github.com/spf13/afero"

	"github.com/metrumresearchgroup/rcmd/v2"
)

func TestRunRBatch(tt *testing.T) {
	type args struct {
		fs      afero.Fs
		rs      rcmd.RSettings
		cmdArgs []string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test R Version",
			args: args{
				fs: afero.NewOsFs(),
				rs: func() rcmd.RSettings {
					rs, _ := rcmd.NewRSettings("")

					return *rs
				}(),
				cmdArgs: []string{
					"--version",
				},
			},
			want: []byte(`R version 3.6.3 (2020-02-29) -- "Holding the Windsock"
 Copyright (C) 2020 The R Foundation for Statistical Computing
Platform: x86_64-apple-darwin15.6.0 (64-bit)
                                
R is free software and comes with ABSOLUTELY NO WARRANTY.
You are welcome to redistribute it under the terms of the
GNU General Public License versions 2 or 3.
For more information about these matters see
https://www.gnu.org/licenses/.


`),
			wantErr: false,
		},
	}

	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			cmd, err := rcmd.New(context.Background(), "", test.args.cmdArgs...)
			t.R.NoError(err)
			co, err := cmd.CombinedOutput()

			t.R.NoError(err)
			t.R.True(bytes.HasPrefix(co, []byte("R version")), "missing prefix 'R version'")
		})
	}
}
