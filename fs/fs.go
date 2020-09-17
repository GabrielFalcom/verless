// Package fs provides functions for filesystem operations.
package fs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

var (
	// MarkdownOnly is a filter that only lets pass Markdown files.
	MarkdownOnly = func(file string) bool {
		return filepath.Ext(file) == ".md"
	}

	// NoUnderscores is a predefined filter that doesn't let pass
	// files starting with an underscore.
	NoUnderscores = func(file string) bool {
		filename := filepath.Base(file)
		return !strings.HasPrefix(filename, "_")
	}

	// ErrStreaming is returned from StreamFiles.
	ErrStreaming error = nil
)

// StreamFiles sends all relative file paths inside a given path that
// match the given filters through the files channel.
func StreamFiles(path string, files chan<- string, filters ...func(file string) bool) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		close(files)
		return nil
	}

	ErrStreaming = filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, filter := range filters {
			if !filter(file) {
				return nil
			}
		}

		// A filepath like `./about-me.md` will be shortened to `about-me.md`.
		// In that case, the filepath already is a relative path.
		if path == "." {
			files <- file
		} else {
			// For paths other than `.`, slice the filepath to a relative path.
			files <- file[len(path):]
		}

		return nil
	})

	close(files)
	return ErrStreaming
}

// MkdirAll creates one or more directories inside the given path.
func MkdirAll(path string, dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(path, dir), 0755); err != nil {
			return err
		}
	}

	return nil
}

// Rmdir removes an entire directory along with its contents. If the
// directory does not exist, nothing happens.
func Rmdir(fs afero.Fs, path string) error {
	_, err := fs.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	return fs.RemoveAll(path)
}
