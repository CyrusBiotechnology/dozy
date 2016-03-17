package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"time"
)

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
	for i := range containers {
		// remove empty container IDs
		if containers[i] == "" {
			containers = append(containers[:i], containers[i+1:]...)
		}
	}
	return containers, nil
}

func stopContainer(containerId string) error {
	if len(containerId) == 0 {
		return errors.New("may not provide empty containerID")
	}
	out := bytes.Buffer{}
	cmd := exec.Command("docker", "stop", "—-time=30", containerId)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return errors.New(out.String())
	}
	return nil
}

func stopAllRunningContainers() error {
	containers, err := getRunningContainers()
	if err != nil {
		return err
	}
	for i := range containers {
		stopContainer(containers[i])
	}
	return nil
}

func sleepOnError(err error, d time.Duration) {
	if err != nil {
		time.Sleep(d)
	}
}

// stop all containers
func stopAllRunningContainersWithRetry() {
	for {
		running, err := getRunningContainers()
		if err != nil {
			sleepOnError(err, time.Second*5)
			continue
		}
		if len(running) != 0 {
			err = stopAllRunningContainers()
			sleepOnError(err, time.Second*5)
			continue
		} else {
			break
		}
	}
	return
}
