// +build darwin linux
package main
import (
	"flag"
	"os"
	"fmt"
	"time"
	"errors"
	"strings"
	"strconv"
	"io/ioutil"
	"path/filepath"
)

var minUptime = flag.Duration("minuptime", 0, "will not exit 0 before uptime >= <minuptime>")
var locks = flag.String("lock", "/tmp/lockfiles/", "where to look for lockfiles")
var lockDur = flag.Duration("duration", time.Minute * 10, "duration for which lock files are considered valid")
var sleepDur = flag.Duration("sleep", 0, "duration to sleep at the end of script before exit 0")

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

func main() {
	flag.Parse()

	logging("/var/log/dozy", " `( ◔ ౪◔)´  dozy")
	Info.Println(fmt.Sprintf("minimum uptime: %v, locks valid for %vm, lock: %v", *minUptime, *lockDur, *locks))

	uptime_str, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		panic(err)
	}
	uptime_secnds, err := strconv.Atoi(strings.Split(string(uptime_str), ".")[0])
	if err != nil {
		panic(err)
	}
	uptime := time.Duration(int(uptime_secnds)) * time.Second

	if uptime < *minUptime {
		panic(fmt.Sprintf("uptime not Ok (%v < %v)", uptime, *minUptime))
	} else {
		fmt.Println("uptime is Ok (%v > %v)", uptime, *minUptime)
	}

	fh, err := os.Open(*locks)
	if err != nil {
		Error.Println(err)
	}
	defer fh.Close()
	locksFh, err := fh.Stat()
	if err != nil {
		panic(err)
		return
	}

	if locksFh.IsDir() {
		// locks is a directory
		Info.Println("input is a directory, scanning")
		fileList := []string{}
		err := filepath.Walk(*locks, func(path string, f os.FileInfo, err error) error {
			if ! f.IsDir() {
				fileList = append(fileList, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		if (len(fileList) == 0) {
			Info.Println("no locks found, killing containers")
			firstDegree()
		} else {
			counter := 0
			for i := range (fileList) {
				isStale, err := lockIsStale(fileList[i])
				if err != nil {
					panic(err)
				}
				if isStale {
					Error.Println(fmt.Sprintf("found stale lock: %v, continuing..", fileList[i]))
				} else {
					panic("valid lock found.")
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
				panic(err)
			}
			if isStale {
				Error.Println("WARNING: stale lock, killing containers and shutting down")
				firstDegree()
			} else {
				panic("valid lock found, exiting.")
			}
	} else {
		panic(errors.New(fmt.Sprintf("fucked up lock(s). -lock=%s", *locks)))
	}
	Info.Println("shutting down...")
	time.Sleep(*sleepDur)
}
