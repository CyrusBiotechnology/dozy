package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Lock struct {
	Path     string
	Modified time.Time
}

func getLockFilesRecursive(path string) ([]Lock, error) {
	files := []Lock{}
	err := filepath.Walk(strings.TrimRight(path, "/"), func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return errors.New("path doesn't exist")
		}
		if !f.IsDir() {
			files = append(files, Lock{
				Path:     path,
				Modified: f.ModTime(),
			})
		}
		return nil
	})
	return files, err
}

// Get locks in a directory `root` or `root` itself if it is a file
func getLocks(root string, maxAge time.Duration) (valid []Lock, invalid []Lock, err error) {
	locks, err := getLockFilesRecursive(root)
	if err != nil {
		return valid, invalid, err
	}

	if len(locks) == 0 {
		return valid, invalid, nil
	} else {
		now := time.Now()
		counter := 0
		for i := range locks {
			counter++
			if locks[i].Modified.After(now.Add(-maxAge)) {
				valid = append(valid, locks[i])
			} else {
				invalid = append(invalid, locks[i])
				continue
			}
		}
		Trace.Println(fmt.Sprintf("%v lock(s) checked, %v valid %v invalid", counter, len(valid), len(invalid)))
	}
	return valid, invalid, nil
}
