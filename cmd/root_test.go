package cmd

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"default", "", "Hello, world!"},
		{"custom name", "Mira", "Hello, Mira!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := greet(tt.input); got != tt.expected {
				t.Fatalf("greet(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
