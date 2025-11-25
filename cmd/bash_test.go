package cmd

import "testing"

func TestGenerateDummyCommand(t *testing.T) {
	t.Parallel()

	// Save and restore original flags
	prevDummy := useDummy
	useDummy = true
	defer func() { useDummy = prevDummy }()

	tests := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "copy txt files",
			prompt:   "copy the .txt files over",
			expected: "cp *.txt /tmp/backup/",
		},
		{
			name:     "disk usage",
			prompt:   "show me disk usage",
			expected: "df -h",
		},
		{
			name:     "fallback",
			prompt:   "do something else",
			expected: "echo 'Dummy command for: do something else'",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := generateDummyCommand(tt.prompt)
			if got != tt.expected {
				t.Fatalf("generateDummyCommand(%q) = %q, want %q", tt.prompt, got, tt.expected)
			}
		})
	}
}

