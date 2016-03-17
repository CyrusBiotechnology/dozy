package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

type dummyTask struct {
	Lock *os.File // Location of lock file
}

func (dt *dummyTask) Cleanup() {
	dt.Lock.Close()
}

func newDummyTask(tempDir string, fakeStartTime time.Time) (dummyTask, error) {
	task := dummyTask{}

	tempFile, err := ioutil.TempFile(tempDir, "dummy-task")
	if err != nil {
		return task, err
	}
	err = os.Chtimes(tempFile.Name(), fakeStartTime, fakeStartTime)
	task.Lock = tempFile
	return task, err
}

func TestGettingValidLocks(t *testing.T) {
	tasks := []dummyTask{}
	tempDir, err := ioutil.TempDir(os.TempDir(), "dozy-test")
	if err != nil {
		t.Fatal("failed to acquire lock directory")
	}
	totalTasks := 10 // 1-indexed
	for i := 1; i <= totalTasks; i++ {
		mockStartTime := time.Duration(i) * time.Minute
		checkDuration := mockStartTime - time.Second
		task, err := newDummyTask(tempDir, time.Now().Local().Add(-mockStartTime))
		if err != nil {
			t.Fatal(err)
		}
		tasks = append(tasks, task)

		valid, invalid, err := getLocks(tempDir, checkDuration)
		if err != nil {
			t.Fatal("failed to get locks")
		}
		if len(valid)+len(invalid) != i {
			t.Fatal(fmt.Sprintf("lock counts don't total correctly (expected %v, got %v)", i, len(valid)+len(invalid)))
		}
		if i > 1 && len(valid) != i-1 {
			t.Fatal(fmt.Sprintf("incorrect number of valid locks found (expected %v, got %v)", i-1, len(valid)))
		}
	}
	// make sure all "tasks" are cleaned up on exit
	for i := range tasks {
		tasks[i].Cleanup()
	}
}

func TestMain(m *testing.M) {
	logging()
	os.Exit(m.Run())
}
