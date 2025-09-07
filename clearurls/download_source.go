package clearurls

// Get the JSON from the web sources, validate the checksum and test it

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

// Pair of URLs with a json to parse and hash to check
type DownloadSource struct {
	data, hash256 string
}

// URLs for distributed versions based on https://docs.clearurls.xyz/1.27.3/specs/rules/
// for use with [DownloadSource.DownloadWithCache] or [DownloadSource.Download]
var (
	SourceGitHub = &DownloadSource{
		data:    "https://rules2.clearurls.xyz/data.minify.json",
		hash256: "https://rules2.clearurls.xyz/rules.minify.hash",
	}
	SourceGitLab = &DownloadSource{
		data:    "https://rules1.clearurls.xyz/data.minify.json",
		hash256: "https://rules1.clearurls.xyz/rules.minify.hash",
	}
)

func (source *DownloadSource) downloadJSON(checkHash bool) ([]byte, error) {
	bodyResponseReader := asyncGetHTTPBody(source.data, "application/json")
	expectedSHA256Str := ""
	if checkHash {
		hashTextBytes, err := asyncGetHTTPBody(source.hash256, "application/octet-stream")()
		if err != nil {
			return nil, err
		}
		expectedSHA256Str = strings.TrimSpace(string(hashTextBytes))
	}
	jsonData, err := bodyResponseReader()
	if err != nil {
		return nil, err
	}
	if checkHash {
		sum := fmt.Sprintf("%x", sha256.Sum256(jsonData))
		if !strings.EqualFold(sum, expectedSHA256Str) {
			return nil, fmt.Errorf(
				"invalid checksum for %q (against %q):\n"+
					"  expected: %q\n"+
					"       got: %q\n",
				source.data, source.hash256, expectedSHA256Str, sum,
			)
		}
		verbose("    Valid hash %q at %s", expectedSHA256Str, time.Now().Format(time.RFC3339))
	}
	testParsed := parseJSON(jsonData)
	if len(testParsed) == 0 {
		return nil, fmt.Errorf("Invalid JSON, no providers found in %q", string(jsonData))
	}
	if _, err := Compile(testParsed); err != nil {
		return nil, fmt.Errorf("Invalid JSON, %w in %q", err, string(jsonData))
	}
	return jsonData, nil
}

func (source *DownloadSource) cachedDownloadJSON(cacheFileName string, cacheMaxAgeM int, checkHash bool) ([]byte, error) {
	return getCachedData(cacheFileName, cacheMaxAgeM, func() ([]byte, error) {
		return source.downloadJSON(checkHash)
	})
}
