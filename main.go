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
)

var locks = flag.String("lock", "/tmp/lockfiles/", "where to look for lockfiles")

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s\n", path)
  return nil
}

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
			fmt.Println(err)
		}
		if len(running) == 0 {
			break
		}
		err = killAllRunningContainers()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func lockIsStale(lockFile string) (bool, error) {
	info, err := os.Stat(lockFile)
	if err != nil {
		return false, err
	}
	lockAge := time.Now().Sub(info.ModTime())
	if lockAge > time.Hour * 2 {
		return true, nil
	} else {
		return false, nil
	}
}

func main() {
	flag.Parse()

	f, err := os.Open(*locks)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
		fi, err := f.Stat()
		if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
		case mode.IsDir():
			// locks is a directory
			fmt.Println("lock is a directory, scanning")
			fileList := []string{}
			err := filepath.Walk(*locks, func(path string, f os.FileInfo, err error) error {
				fileList = append(fileList, path)
				return nil
			})
			if (len(fileList) == 1) {
				fmt.Println("killing containers")
				firstDegree()
				fmt.Println("shutting down...")
			}
			fmt.Printf("filepath.Walk() returned %v\n", err)
		case mode.IsRegular():
			// locks is a file
			fmt.Println("lock is a file, checking...")
			isStale, err := lockIsStale(*locks)
			if err != nil {
				panic(err)
			}
			if isStale {
				fmt.Println("WARNING: stale lock, killing containers and shutting down")
				firstDegree()
			} else {
				fmt.Println("valid lock found, exiting.")
			}
	}
}
