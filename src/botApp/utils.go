package main

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"math/rand"
)

func Random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min) + min
}

func DownloadFile(url string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("File download error: " + resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}