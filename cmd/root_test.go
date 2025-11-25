package cmd

import "testing"

func TestGenerateCommand(t *testing.T) {
	t.Parallel()

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

			if got := generateCommand(tt.prompt); got != tt.expected {
				t.Fatalf("generateCommand(%q) = %q, want %q", tt.prompt, got, tt.expected)
			}
		})
	}
}
