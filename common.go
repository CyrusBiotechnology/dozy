package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Get files recursively
func getFilesRecursive(path string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(strings.TrimRight(path, "/"), func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return errors.New("file / folder doesn't exist")
		}
		if !f.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
