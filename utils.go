package main

import (
	"fmt"
	"os"
)

func getVersion() string {
	versionStr := fmt.Sprintf("v%v", VERSION[0])
	for i := 1; i < len(VERSION); i++ {
		versionStr = fmt.Sprintf("%v.%v", versionStr, VERSION[i])
	}
	return versionStr
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
