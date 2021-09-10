package filters

import (
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestDropInput(tt *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "not input",
			args: args{
				bs: []byte("hello world"),
			},
			want: []byte("hello world"),
		},
		{
			name: "is input",
			args: args{
				bs: []byte("> hello world"),
			},
			want: []byte(nil),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)
			got := DropInput(test.args.bs)
			t.R.Equal(test.want, got)
		})
	}
}

func TestLineNumberStripper(tt *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "no line number",
			args: args{[]byte("no line number")},
			want: []byte("no line number"),
		},
		{
			name: "has line number",
			args: args{[]byte("  [14]  has line number")},
			want: []byte("has line number"),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)
			got := LineNumberStripper(test.args.bs)
			t.R.Equal(test.want, got)
		})
	}
}

func TestNewPrefixFilter(tt *testing.T) {
	tests := []struct {
		name   string
		prefix string
		input  []byte
		want   []byte
	}{
		{
			name:   "no prefix",
			prefix: "",
			input:  []byte("hello world"),
			want:   []byte("hello world"),
		},
		{
			name:   "with prefix",
			prefix: "prefix",
			input:  []byte("hello world"),
			want:   []byte("prefix hello world"),
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)
			fn := NewPrefixFilter(test.prefix)
			got := fn(test.input)
			t.R.Equal(test.want, got)
		})
	}
}
