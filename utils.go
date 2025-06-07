package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func deleteFilesInDir(dirPath string) error {
	files, err := os.ReadDir(dirPath) // Read directory contents
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name()) // Construct full file path
		if !file.IsDir() {                              // Skip directories
			err := os.Remove(filePath) // Delete file
			if err != nil {
				fmt.Println("Error deleting:", filePath, err)
			}
		}
	}
	return nil
}

// Writes mangaCache to given yaml file
func writeYamlToFile(mangaCache MangaCache, file *os.File) {
	bytes, err := yaml.Marshal(mangaCache)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(bytes)

	if err != nil {
		panic(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir() // Ensure it's not a directory
}

func httpGetRespWithRetries(url string, retries int) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers (some servers require a User-Agent)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	for attempt := 1; attempt <= retries+1; attempt++ {
		resp, err := http.DefaultClient.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		return resp, nil
	}
	return nil, fmt.Errorf("could not download after retries")
}

func httpGetAndReturnBytes(url string) []byte {
	resp, err := httpGetRespWithRetries(url, 20)

	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	time.Sleep(200 * time.Millisecond)
	return bytes
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false // Path does not exist or error occurred
	}
	return info.IsDir()
}
