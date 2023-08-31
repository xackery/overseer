package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Start starts a download of provided url
func Save(url string, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("status: %s", resp.Status)
	}

	// Create the file
	w, err := os.Create(path)
	if err != nil {
		return false, fmt.Errorf("create: %w", err)
	}

	// Write the body to file
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return false, fmt.Errorf("copy: %w", err)
	}

	return false, nil
}

// Get fetches a url and returns the body
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("readall: %w", err)
	}

	return body, nil
}
