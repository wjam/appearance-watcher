package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteToFile(t *testing.T) {
	tests := []struct {
		name          string
		existingFiles map[string]string
		expectedFiles map[string]string
		subject       func() error
	}{
		{
			name:          "creates new file",
			existingFiles: nil,
			expectedFiles: map[string]string{
				"file.txt": "dark",
			},
			subject: writeToFile{
				Path:  "file.txt",
				Light: "light",
				Dark:  "dark",
			}.ToDark,
		},
		{
			name: "updates file",
			existingFiles: map[string]string{
				"file.txt": "dark",
			},
			expectedFiles: map[string]string{
				"file.txt": "light",
			},
			subject: writeToFile{
				Path:  "file.txt",
				Light: "light",
				Dark:  "dark",
			}.ToLight,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()

			for k, v := range test.existingFiles {
				require.NoError(t, os.WriteFile(filepath.Join(dir, k), []byte(v), 0644))
			}

			t.Chdir(dir)

			err := test.subject()
			require.NoError(t, err)

			for k, v := range test.expectedFiles {
				file := filepath.Join(dir, k)

				content, err := os.ReadFile(file)
				require.NoError(t, err)
				assert.Equal(t, v, string(content))
			}
		})
	}
}
