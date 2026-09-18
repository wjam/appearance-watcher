package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSymLinkFile(t *testing.T) {
	tests := []struct {
		name          string
		existingFiles map[string]string
		existingLinks map[string]string
		subject       func() error
		expectedLinks map[string]string
	}{
		{
			name: "creates links",
			existingFiles: map[string]string{
				"light.txt": "light",
				"dark.txt":  "dark",
			},
			existingLinks: nil,
			subject: symLinkFile{
				Path:  "link.txt",
				Light: "light.txt",
				Dark:  "dark.txt",
			}.ToDark,
			expectedLinks: map[string]string{
				"link.txt": "dark.txt",
			},
		},
		{
			name: "updates existing links",
			existingFiles: map[string]string{
				"light.txt": "light",
				"dark.txt":  "dark",
			},
			existingLinks: map[string]string{
				"link.txt": "dark.txt",
			},
			subject: symLinkFile{
				Path:  "link.txt",
				Light: "light.txt",
				Dark:  "dark.txt",
			}.ToLight,
			expectedLinks: map[string]string{
				"link.txt": "light.txt",
			},
		},
		{
			name: "overwrites tmp file",
			existingFiles: map[string]string{
				"light.txt":    "light",
				"dark.txt":     "dark",
				"link.txt.tmp": "foo",
			},
			existingLinks: nil,
			subject: symLinkFile{
				Path:  "link.txt",
				Light: "light.txt",
				Dark:  "dark.txt",
			}.ToDark,
			expectedLinks: map[string]string{
				"link.txt": "dark.txt",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()

			for k, v := range test.existingFiles {
				require.NoError(t, os.WriteFile(filepath.Join(dir, k), []byte(v), 0644))
			}
			for k, v := range test.existingLinks {
				require.NoError(t, os.Symlink(filepath.Join(dir, v), filepath.Join(dir, k)))
			}

			t.Chdir(dir)

			err := test.subject()
			require.NoError(t, err)

			for k, v := range test.expectedLinks {
				file := filepath.Join(dir, k)

				link, err := os.Readlink(file)
				require.NoError(t, err)

				assert.Equal(t, v, link)
			}
		})
	}
}
