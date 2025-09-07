package clearurls

// Generic method of describing a source from which to obtain `RunnableProvider`s

import (
	"fmt"
	"regexp"
	"strconv"
)

type parseSource struct {
	sourceName    string
	cacheFilename string
	cacheMaxAgeM  int
}

// Parse arguments of the format `<source>[:<cache_filename>[:<cache_max_age_minutes>]]`
func parseSourceArgument(source string) (*parseSource, error) {
	matches := regexp.MustCompile("(?i)^([^:]+):?(?:(.+?)(?::(\\d*))?)?$").FindStringSubmatch(source)
	if matches == nil {
		return nil, fmt.Errorf("Invalid source argument %q", source)
	}
	result := &parseSource{
		sourceName:    matches[1],
		cacheFilename: matches[2],
		cacheMaxAgeM:  -1,
	}
	if matches[3] != "" {
		num, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("Invalid source argument %q (%w)", source, err)
		}
		if num > 0 {
			result.cacheMaxAgeM = num
		}
	}
	return result, nil
}

var sources = map[string]*DownloadSource{
	"github": SourceGitHub,
	"gitlab": SourceGitLab,
}

func downloadSource(source, cache string, cacheMaxAgeM int) ([]RunnableProvider, error) {
	sourceURLs := sources[source]
	if sourceURLs == nil {
		return nil, fmt.Errorf("Invalid source %q", source)
	}
	if cache == "" {
		return sourceURLs.Download(true)
	}
	return sourceURLs.DownloadWithCache(cache, cacheMaxAgeM, true)
}

// Get providers from a string of the format `<source>[:<cache_filename>[:<cache_max_age_minutes>]]`
//
// Where `<source>` can be one of `hardcoded`, `github` or `gitlab`
//
// Warning: If not `hardcoded`, the providers returned are not compiled
//
// Examples:
//
// - Load hardcoded providers. Will fail if they weren't generated. Result should be compiled.
//
//	clearurls.GetProvidersFromSourceArgument("hardcoded")
//	// Equivalent to: clearurls.HardcodedProviders()
//
// - Download from the web each time
//
//	clearurls.GetProvidersFromSourceArgument("github")
//	clearurls.GetProvidersFromSourceArgument("gitlab")
//	// Equivalent to: clearurls.SourceGitHub.Download(true)
//
// - Download using a local cache file (never updating)
//
//	clearurls.GetProvidersFromSourceArgument("gitlab:/var/run/clearurls_cache.json")
//	// Equivalent to: clearurls.SourceGitLab.DownloadWithCache("/var/run/clearurls_cache.json", -1, true)
//
// - Download using a local cache file, refreshing if the file is older than 1 hour
//
//	clearurls.GetProvidersFromSourceArgument("github:/var/run/clearurls_cache.json:60")
//	// Equivalent to: clearurls.SourceGitHub.DownloadWithCache("/var/run/clearurls_cache.json", 60, true)
func GetProvidersFromSourceArgument(source string) ([]RunnableProvider, error) {
	parsedSource, err := parseSourceArgument(source)
	if err != nil {
		return nil, err
	}
	if parsedSource.sourceName == "hardcoded" {
		return MustHaveHardcodedProviders()
	} else {
		return downloadSource(parsedSource.sourceName, parsedSource.cacheFilename, parsedSource.cacheMaxAgeM)
	}
}
