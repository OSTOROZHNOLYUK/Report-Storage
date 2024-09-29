package reports

import "testing"

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
