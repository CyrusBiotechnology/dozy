package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func dummyContainer() (name string, stderr string, err error) {
	stdout := bytes.Buffer{}
	stderrB := bytes.Buffer{}
	cmd := exec.Command("docker", "run", "-d", "busybox", "sh", "-c", "sleep 30")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderrB
	err = cmd.Run()
	if err != nil {
		return stdout.String(), stderrB.String(), err
	}
	return stdout.String(), stderrB.String(), nil
}

func TestGettingRunningContainers(t *testing.T) {
	if os.Getenv("DOCKER_HOST") == "" {
		t.Skip("you need to set DOCKER_HOST to test docker interactions")
	}
	_, err := getRunningContainers()
	if err != nil {
		t.Fatal("failed to get docker containers")
	}
}

func TestStoppingContainer(t *testing.T) {
	if os.Getenv("DOCKER_HOST") == "" {
		t.Skip("you need to set DOCKER_HOST to test docker interactions")
	}
	containers, err := getRunningContainers()
	if err != nil {
		t.Skip("problem getting containers")
	}
	if len(containers) != 0 {
		t.Skip("you have other containers running on this system, I'm not very picky. Skipping...")
	}
	container, stderr, err := dummyContainer()
	if err != nil {
		t.Log(err)
		t.Log(container)
		t.Log(stderr)
		t.Fatal("container creation failed")
	}
	stopAllRunningContainersWithRetry()
	containers, err = getRunningContainers()
	if err != nil {
		t.Fatal("error getting containers at the end of the test, containers may still be running")
	}
	if len(containers) != 0 {
		t.Fatal("containers left running")
	}
}
