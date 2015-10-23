package main
import (
	"os/exec"
	"bytes"
	"strings"
	"log"
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
