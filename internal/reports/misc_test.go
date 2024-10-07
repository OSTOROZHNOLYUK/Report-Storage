package reports

import (
	"reflect"
	"testing"
)

func Test_generateFileNameJPEG(t *testing.T) {
	tests := []struct {
		name    string
		notWant string
	}{
		{
			name:    "OK",
			notWant: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateFileNameJPEG(); got == tt.notWant {
				t.Errorf("generateFileName() = %s", got)
			}
		})
	}
}

func TestSliceDiff(t *testing.T) {
	type args struct {
		origin []string
		new    []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Origin more than New",
			args: args{
				origin: []string{"one", "two", "three"},
				new:    []string{"one", "three"},
			},
			want: []string{"two"},
		},
		{
			name: "New slice empty",
			args: args{
				origin: []string{"one", "two", "three"},
				new:    []string{},
			},
			want: []string{"one", "two", "three"},
		},
		{
			name: "Both slices empty",
			args: args{
				origin: []string{},
				new:    []string{},
			},
			want: nil,
		},
		{
			name: "Equal slices",
			args: args{
				origin: []string{"one", "two", "three"},
				new:    []string{"one", "two", "three"},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceDiff(tt.args.origin, tt.args.new); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckMedia() = %v, want %v", got, tt.want)
			}
		})
	}
}
