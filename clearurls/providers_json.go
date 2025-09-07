package clearurls

// Parse JSON data into `providerJSON`, a `RunnableProvider` with just the
// the raw data as it came in from the JSON distribution

import (
	"encoding/json"
	"fmt"
)

// Unprocessed JSON source data rulesets called "providers" in ClearURL lingo
type providerJSON struct {
	name              string
	URLPattern        string
	CompleteProvider  bool
	Rules             []string
	RawRules          []string
	ReferralMarketing []string
	Exceptions        []string
	Redirections      []string
	// ForceRedirection  bool // Applies only to web
}

// Debug print for `providerJSON`
func (provider *providerJSON) String() string {
	header := fmt.Sprintf("providerJSON %q", provider.name)
	fieldCountsString := ""
	totalCount := 0
	appendLenOf := func(name string, items []string) {
		if count := len(items); count > 0 {
			fieldCountsString = fieldCountsString + fmt.Sprintf(" %s:%d", name, count)
			totalCount += count
		}
	}
	appendLenOf("r", provider.Rules)
	appendLenOf("rr", provider.RawRules)
	appendLenOf("rm", provider.ReferralMarketing)
	appendLenOf("e", provider.Exceptions)
	appendLenOf("re", provider.Redirections)
	if provider.CompleteProvider {
		if totalCount == 0 {
			fieldCountsString = " CompleteProvider"
		} else {
			fieldCountsString = fmt.Sprintf(" Very weird, it's completeProvider=true, but has%s. Inconsistent with documentation but might be ok in Addon code", fieldCountsString)
		}
	}
	return header + fieldCountsString + "\n"
}

// Parse the JSON into an array of `providerJSON`
func parseJSON(jsonData []byte) []RunnableProvider {
	type clearURLsRoot struct {
		Providers map[string]providerJSON
	}
	var parsedRules clearURLsRoot
	json.Unmarshal(jsonData, &parsedRules)
	providers := make([]RunnableProvider, len(parsedRules.Providers))
	i := 0
	for key, provider := range parsedRules.Providers {
		provider.name = key
		providers[i] = &provider
		i++
	}
	return providers
}

// Download either source, optionally checking hash, does not use any cached file.
// The returned value must have the valid hash if requested, and must be valid JSON
// that can be compiled. This allows for caching and dealing with only valid values.
func (source *DownloadSource) Download(checkHash bool) ([]RunnableProvider, error) {
	data, err := source.downloadJSON(checkHash)
	if err != nil {
		return nil, err
	}
	return parseJSON(data), nil
}

// Download from the provided `source` (`SourceGitHub` or `SourceGitLab`) the latest rules file, and return
// a list of "providers" as they are called that can be used to process URLs, but are un-processed (no precompilation etc)
//
// ## Cache
//
//   - The parent folder of the provided `cacheFileName` is immediately created
//   - If the file exists, and `cacheMaxAgeM` is `-1`, or is older than the file's modification time,
//     it is read and returned as is
//   - If the file doesn't exist, or is older than `cacheMaxAgeM`, the file is retreived as per [DownloadSource.Download]
//   - The cache file is written with the raw json
func (source *DownloadSource) DownloadWithCache(cacheFileName string, cacheMaxAgeM int, checkHash bool) ([]RunnableProvider, error) {
	data, err := source.cachedDownloadJSON(cacheFileName, cacheMaxAgeM, checkHash)
	if err != nil {
		return nil, err
	}
	return parseJSON(data), nil
}
