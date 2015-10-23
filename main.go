package main
import (
	"flag"
	"os"
	"fmt"
	"path/filepath"
	"time"
	"errors"
	"log"
)

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

	log.Println("")
	log.Println("         `( ◔ ౪◔)´")
	log.Println("                   dozey")
	log.Println("")
	log.Println(fmt.Sprintf("locks valid for %vm, lock: %v", *lockDur, *locks))

	fh, err := os.Open(*locks)
	if err != nil {
		log.Println(err)
	}
	defer fh.Close()
	locksFh, err := fh.Stat()
	if err != nil {
		panic(err)
		return
	}

	if locksFh.IsDir() {
		// locks is a directory
		log.Println("lock is a directory, scanning")
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
			log.Println("no locks found, killing containers")
			firstDegree()
		} else {
			counter := 0
			for i := range (fileList) {
				isStale, err := lockIsStale(fileList[i])
				if err != nil {
					panic(err)
				}
				if isStale {
					log.Println(fmt.Sprintf("found stale lock: %v, continuing..", fileList[i]))
				} else {
					panic("valid lock found, exiting...")
				}
				counter++
			}
			log.Println(fmt.Sprintf("%v lock(s) checked", counter))
		}
	} else if locksFh.Mode().IsRegular() {
			// locks is a file
			log.Println("lock is a file, checking...")
			isStale, err := lockIsStale(*locks)
			if err != nil {
				panic(err)
			}
			if isStale {
				log.Println("WARNING: stale lock, killing containers and shutting down")
				firstDegree()
			} else {
				panic("valid lock found, exiting.")
			}
	} else {
		panic(errors.New(fmt.Sprintf("fucked up lockfile. -lock=%s", *locks)))
	}
	log.Println("shutting down...")
	time.Sleep(*sleepDur)
}
