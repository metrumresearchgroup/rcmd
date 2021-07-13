package rcmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRunRBatch(t *testing.T) {
	type args struct {
		fs      afero.Fs
		rs      RSettings
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
				rs: func() RSettings {
					rs, _ := NewRSettings("")

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.args.rs.RunR(context.Background(), "", tt.args.cmdArgs...)
			assert.Equal(t, nil, err, "error")

			msg := fmt.Sprintf("\ngot<\n%v\n>\nwant<\n%v\n>", string(got.Output), string(tt.want))
			assert.True(t, bytes.HasPrefix(got.Output, []byte("R version")), msg)

		})
	}
}

func BenchmarkRunR(b *testing.B) {
	rs, err := NewRSettings("R")
	if err != nil {
		b.Fatalf("non-nil err: %v", err)
	}
	for n := 0; n < b.N; n++ {
		_, err := rs.RunR(context.Background(), "", "--version")
		if err != nil {
			b.Fatalf("caught error: %v", err)
		}
	}
}
