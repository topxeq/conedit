package editor

import "testing"

func TestScreen_CalculateVisualWidth(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"你好", 4},
		{"hi 你好", 7},
	}
	for _, tt := range tests {
		result := CalculateVisualWidth(tt.input)
		if result != tt.expected {
			t.Errorf("CalculateVisualWidth(%q) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}
