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
func getValidLocks(root string, maxAge time.Duration) ([]Lock, error) {
	validated := []Lock{}
	locks, err := getLockFilesRecursive(root)
	if err != nil {
		return validated, err
	}

	if len(locks) == 0 {
		Info.Println("no locks found, stopping all containers")
		stopAllRunningContainersWithRetry()
	} else {
		now := time.Now()
		counter := 0
		for i := range locks {
			if locks[i].Modified.Before(now.Add(-maxAge)) {
				Info.Println(fmt.Sprintf("found stale lock: %v", locks[i].Path))
				validated = append(validated, locks[i])
			} else {
				continue
			}
			counter++
		}
		Trace.Println(fmt.Sprintf("%v lock(s) checked", counter))
	}
	return validated, nil
}
