package main
import (
	"flag"
	"os"
	"fmt"
	"path/filepath"
	"os/exec"
	"bytes"
	"strings"
	"time"
	"errors"
	"log"
)

var locks = flag.String("lock", "/tmp/lockfiles/", "where to look for lockfiles")
var lockDur = flag.Int("duration", 30, "duration in minutes for which lock files are considered valid")
var sleepTime = flag.Int("sleep", 0, "duration to sleep at the end of script before exit 0")

func getRunningContainers() ([]string, error) {
	containers := make([]string, 0)
	cmd := exec.Command("docker", "ps", "-q")
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return containers, err
	}
	containers = strings.Split(out.String(), "\n")
	for i := range(containers) {
		// remove empty container IDs
		if containers[i] == "" {
			containers = append(containers[:i], containers[i+1:]...)
		}
	}
	return containers, nil
}

func killContainer(containerId string) error {
	out := bytes.Buffer{}
	cmd := exec.Command("docker", "kill", containerId)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func killAllRunningContainers() error {
	containers, err := getRunningContainers()
	if err != nil {
		return err
	}
	for i := range(containers) {
		killContainer(containers[i])
	}
	return nil
}

// kill all containers or die trying
func firstDegree() {
	for {
		running, err:= getRunningContainers()
		if err != nil {
			log.Println(err)
		}
		if len(running) == 0 {
			break
		}
		err = killAllRunningContainers()
		if err != nil {
			log.Println(err)
		}
	}
}

func lockIsStale(lockFile string) (bool, error) {
	info, err := os.Stat(lockFile)
	if err != nil {
		return false, err
	}
	lockAge := time.Now().Sub(info.ModTime())
	if lockAge > time.Minute * time.Duration(*lockDur) {
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
	time.Sleep(time.Second * time.Duration(*sleepTime))
}
