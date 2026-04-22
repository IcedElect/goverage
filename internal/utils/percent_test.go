package utils

import "testing"

func TestPercent(t *testing.T) {
	tests := []struct {
		name    string
		covered int64
		total   int64
		want    float64
	}{
		{"zero total", 0, 0, 0},
		{"zero covered", 0, 100, 0},
		{"half covered", 50, 100, 50},
		{"full covered", 100, 100, 100},
		{"rounding", 33, 100, 33},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Percent(tt.covered, tt.total); got != tt.want {
				t.Errorf("Percent() = %v, want %v", got, tt.want)
			}
		})
	}
}
