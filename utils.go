// Contains code from Stack Overflow: http://stackoverflow.com/a/35615565/1342445
package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52 possibilities
	letterIdxBits = 6                                                      // 6 bits to represent 64 possibilities / indexes
	letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
)

// SecureRandomBytes returns the requested number of bytes using crypto/rand
func SecureRandomBytes(length int) ([]byte, error) {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return randomBytes, errors.New("Unable to generate random bytes")
	}
	return randomBytes, nil
}

func SecureRandomAlphaString(length int) (str string, err error) {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes, err = SecureRandomBytes(bufferSize)
			if err != nil {
				return string(result), err
			}
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}

	return string(result), nil
}

func getVersion() string {
	versionStr := fmt.Sprintf("v%v", VERSION[0])
	for i := 1; i < len(VERSION); i++ {
		versionStr = fmt.Sprintf("%v.%v", versionStr, VERSION[i])
	}
	return versionStr
}

func getUptime() (time.Duration, error) {
	uptime := time.Duration(0)
	uptime_str, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return uptime, err
	}
	uptime_seconds, err := strconv.Atoi(strings.Split(string(uptime_str), ".")[0])
	if err != nil {
		return uptime, err
	}
	uptime = time.Duration(int(uptime_seconds)) * time.Second
	return uptime, nil
}

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
