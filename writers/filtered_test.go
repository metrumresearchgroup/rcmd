package writers_test

import (
	"bytes"
	"testing"

	. "github.com/metrumresearchgroup/wrapt"

	"github.com/metrumresearchgroup/rcmd/writers"
)

func Test_Filters(tt *testing.T) {
	type args struct {
		filters []writers.FilterFunc
	}

	tests := []struct {
		name        string
		msg         string
		args        args
		want        string
		wantErr     bool
		expectPanic bool
	}{
		{
			name: "no filters",
			msg:  "no filters\nno problems",
			want: "no filters\nno problems",
		},
		{
			name: "no prefix",
			args: args{filters: []writers.FilterFunc{writers.NewPrefixFilter("")}},
			msg:  "no prefix\nno problems",
			want: "no prefix\nno problems\n",
		},
		{
			name: "prefix",
			args: args{
				filters: []writers.FilterFunc{writers.NewPrefixFilter("prefix")},
			},
			msg:  "no problems\n",
			want: "prefix no problems\n",
		},
		{
			name: "combined",
			args: args{filters: []writers.FilterFunc{
				writers.InputFilter,
				writers.LineNumberStripper,
				bytes.TrimSpace,
				writers.NewPrefixFilter("out"),
			}},
			msg:  "> input\n \n [10000] output",
			want: "out output\n",
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := WrapT(tt)
			if test.expectPanic {
				defer func() {
					p := recover()
					if p == nil {
						t.Fatal("panic not raised")
					}
				}()
			}
			w := &nopcloser{&bytes.Buffer{}}

			sut := writers.NewFilter(w, test.args.filters...)
			n, err := sut.Write([]byte(test.msg))
			t.R.WantError(test.wantErr, err)
			if !test.wantErr {
				t.A.NotZero(n)
			}
			bs := w.Bytes()
			t.A.Equal(test.want, string(bs))
		})
	}
}

type nopcloser struct {
	*bytes.Buffer
}

func (nc *nopcloser) Close() error {
	// nop
	return nil
}
