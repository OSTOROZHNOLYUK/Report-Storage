package api

import (
	"Report-Storage/internal/storage"
	"reflect"
	"testing"
)

func Test_splitStatus(t *testing.T) {
	var nilSlice []storage.Status
	tests := []struct {
		name  string
		param string
		want  []storage.Status
	}{
		{
			name:  "One status",
			param: "1",
			want:  []storage.Status{1},
		},
		{
			name:  "Two statuses",
			param: "1,2",
			want:  []storage.Status{1, 2},
		},
		{
			name:  "Empty string",
			param: "",
			want:  nilSlice,
		},
		{
			name:  "Fully incorrect",
			param: "asdf",
			want:  nilSlice,
		},
		{
			name:  "Partial incorrect",
			param: "1,asdf,3",
			want:  []storage.Status{1, 3},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := splitStatus(tt.param); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
