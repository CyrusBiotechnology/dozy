package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func GetVersion(t *testing.T) {
	number_re, err := regexp.Compile("([0-9]+)")
	if err != nil {
		t.Fatal("test encountered an error")
	}
	dot_re, err := regexp.Compile("\\.")
	if err != nil {
		t.Fatal("test encountered an error")
	}

	version_from_function := getVersion()
	numbers := number_re.FindAllString(version_from_function, -1)
	dots := dot_re.FindAllString(version_from_function, -1)

	if len(numbers) != 3 {
		t.Fatal("too many version groups (should be MAJOR.MINOR.BUGFIX)")
	}

	reconstructed := []string{
		numbers[0], dots[0],
		numbers[1], dots[1],
		numbers[2], dots[2],
	}
	check_string := strings.Join(reconstructed, "")

	if version_from_function != check_string {
		t.Fatal("unexpected results from version function")
	}
}

func TestGettingUptime(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("can't get uptime on darwin yet...")
	}
	uptime_function_output, err := getUptime()
	uptime_str, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		t.Fatal(err)
	}
	uptime_seconds, err := strconv.Atoi(strings.Split(string(uptime_str), ".")[0])
	if err != nil {
		t.Fatal(err)
	}
	uptime := time.Duration(int(uptime_seconds)) * time.Second
	delta := uptime.Seconds() - uptime_function_output.Seconds()
	if delta > 2 {
		t.Fatal(fmt.Sprintf("shouldn't be more than a two second delta between uptime checks. There was a delta of %v", delta))
	}
}
