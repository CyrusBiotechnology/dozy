// +build darwin linux
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var VERSION = [...]int{1, 0, 6}

func getVersion() string {
	versionStr := fmt.Sprintf("v%v", VERSION[0])
	for i := 1; i < len(VERSION); i++ {
		versionStr = fmt.Sprintf("%v.%v", versionStr, VERSION[i])
	}
	return versionStr
}

var minUptime = flag.Duration("minuptime", 0, "will not exit 0 before uptime >= <minuptime>")
var locks = flag.String("lock", "/tmp/lockfiles/", "where to look for lockfiles")
var lockDur = flag.Duration("duration", time.Minute*10, "duration for which lock files are considered valid")
var sleepDur = flag.Duration("sleep", 0, "duration to sleep at the end of script before exit 0")

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

func lockIsStale(lockFile string) (bool, error) {
	info, err := os.Stat(lockFile)
	if err != nil {
		return false, err
	}
	lockAge := time.Now().Sub(info.ModTime())
	if lockAge > *lockDur {
		return true, nil
	} else {
		return false, nil
	}
}

func exit(message string) {
	Info.Println(message)
	os.Exit(1)
}

func exitE(message error) {
	Error.Println(message)
	os.Exit(2)
}

func main() {
	flag.Parse()

	logging("/var/log/dozy", fmt.Sprint(" `( ◔ ౪◔)´  dozy ", getVersion()))
	Info.Println(fmt.Sprintf("minimum uptime: %v, locks valid for %v, lock: %v", *minUptime, *lockDur, *locks))

	if locksPlaceExists, _ := exists(*locks); !locksPlaceExists {
		exitE(errors.New("specified locks location doesn't exist"))
	}

	fh, err := os.Open(*locks)
	if err != nil {
		Error.Println(err)
	}
	defer fh.Close()
	locksFh, err := fh.Stat()
	if err != nil {
		exitE(err)
	}

	if locksFh.IsDir() {
		// locks is a directory
		Info.Println("input is a directory, scanning")
		fileList := []string{}
		err := filepath.Walk(*locks, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				fileList = append(fileList, path)
			}
			return nil
		})
		if err != nil {
			exitE(err)
		}
		if len(fileList) == 0 {
			Info.Println("no locks found, stopping containers")
			err := firstDegree()
			if err != nil {
				exitE(err)
			}
		} else {
			counter := 0
			for i := range fileList {
				isStale, err := lockIsStale(fileList[i])
				if err != nil {
					exitE(err)
				}
				if isStale {
					Error.Println(fmt.Sprintf("found stale lock: %v, continuing..", fileList[i]))
				} else {
					exitE(errors.New("valid lock found."))
				}
				counter++
			}
			Info.Println(fmt.Sprintf("%v lock(s) checked", counter))
		}
	} else if locksFh.Mode().IsRegular() {
		// locks is a file
		Info.Println("lock is a file, checking...")
		isStale, err := lockIsStale(*locks)
		if err != nil {
			exitE(err)
		}
		if isStale {
			Error.Println("WARNING: stale lock, killing stopping and shutting down")
			firstDegree()
		} else {
			exitE(errors.New("valid lock found, exiting."))
		}
	} else {
		exitE(errors.New(fmt.Sprintf("fucked up lock(s). -lock=%s", *locks)))
	}
	Info.Println("shutting down...")
	time.Sleep(*sleepDur)
}
