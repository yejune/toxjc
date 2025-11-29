package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yejune/toxjc/lib"
)

func TestDetect(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		content  []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "CSV file",
			content:  []byte("a,b,c\n1,2,3"),
			expected: "csv",
		},
		{
			name:     "JSON array",
			content:  []byte(`[{"a":1}]`),
			expected: "json",
		},
		{
			name:     "JSON object",
			content:  []byte(`{"a":1}`),
			expected: "json",
		},
		{
			name:     "XLSX file",
			content:  []byte{0x50, 0x4B, 0x03, 0x04},
			expected: "xlsx",
		},
		{
			name:     "XLS file",
			content:  []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1},
			expected: "xls",
		},
		{
			name:    "Binary file",
			content: []byte{0x00, 0x01, 0x02},
			wantErr: true,
		},
		{
			name:    "Plain text (no delimiter)",
			content: []byte("hello world"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name)
			if err := os.WriteFile(path, tt.content, 0644); err != nil {
				t.Fatal(err)
			}

			result, err := toxjc.Detect(path)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}
