package clearurls

// Handle HTTP, and file cache aspects of obtaining the raw JSON from ClearURLs distribution

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ensureParentFolderExists(filename string) error {
	dir, _ := filepath.Split(filename)
	if dir == "" {
		return nil
	}
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		} else {
			return err
		}
	}
	return nil
}

// Setting this to `true` will enable some stderr printing
// when downloading and using the cache
var Verbose = false

func verbose(format string, a ...any) {
	if Verbose {
		fmt.Fprintf(os.Stderr, format+"\n", a...)
	}
}

// Gets the return of running `miss()` using a cache file `cacheRootFolder/cacheFileName`.
//
// Will create the folder `cacheRootFolder` even if `miss()` fails, leaving it
// empty (this helps to early check write permissions).
func getCachedData(cacheFileName string, cacheMaxAgeM int, miss func() ([]byte, error)) ([]byte, error) {
	if err := ensureParentFolderExists(cacheFileName); err != nil {
		return nil, err
	}
	filename := filepath.Join(cacheFileName)
	stat, statErr := os.Stat(filename)
	var bytesFromCache []byte
	if statErr != nil {
		if !os.IsNotExist(statErr) {
			return nil, statErr
		}
	} else {
		fileAge := time.Since(stat.ModTime()).Truncate(time.Second)
		verbose("Cache: Using %q (%s old)", cacheFileName, fileAge.String())
		var errReading error
		bytesFromCache, errReading = os.ReadFile(filename)
		if errReading != nil {
			return nil, errReading
		}
		if cacheMaxAgeM < 1 || cacheMaxAgeM > int(fileAge.Minutes()) {
			return bytesFromCache, nil
		} else {
			verbose("Cache: Expired - re-downloading %q", cacheFileName)
		}
	}
	text, err := miss()
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(filename, []byte(text), 0644)
	if err != nil {
		return nil, err
	}
	return text, nil
}

// Download `url` with `GET`, check `Content-Type` matches `expectedMIME` and return the body
func getHTTPBody(url, expectedMIME string) ([]byte, error) {
	var client http.Client
	verbose("GET %q started", url)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to GET %q: %v", url, err)
	}
	defer resp.Body.Close()
	verbose("    %q ended with %d", url, resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to GET %q: unexpected response code %d", url, resp.StatusCode)
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), expectedMIME) {
		return nil, fmt.Errorf("failed to GET %q: wrong mime type (%q - expected %q)", url, resp.Header.Get("Content-Type"), expectedMIME)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func asyncGetHTTPBody(url, expectedMIME string) func() ([]byte, error) {
	type getHTTPBodyFuncResponse struct {
		data []byte
		err  error
	}
	resultChannel := make(chan getHTTPBodyFuncResponse)
	go func() {
		data, err := getHTTPBody(url, expectedMIME)
		resultChannel <- getHTTPBodyFuncResponse{data, err}
	}()
	return func() ([]byte, error) {
		result := <-resultChannel
		close(resultChannel)
		return result.data, result.err
	}
}
