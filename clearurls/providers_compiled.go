package clearurls

// Create a `providerCompiled`, a `RunnableProvider`, where each field is a
// pre-compiled `regexp.Regexp`, for matching performance

import (
	"regexp"
)

// Provider with as much pre-compilation of regexen as useful
// based on `providerJSON`
type providerCompiled struct {
	name              string
	URLPattern        *regexp.Regexp
	CompleteProvider  bool
	Rules             *regexp.Regexp
	RawRules          *regexp.Regexp
	Exceptions        *regexp.Regexp
	ReferralMarketing *regexp.Regexp
	Redirections      []*regexp.Regexp
}

// implements RunnableProvider
func (provider *providerWithPreparedRegexStr) compile() (*providerCompiled, error) {
	compileRegexpIfNotEmpty := func(rxStr string) (*regexp.Regexp, error) {
		if rxStr == "" {
			return nil, nil
		}
		return regexp.Compile(rxStr)
	}
	result := &providerCompiled{
		// ForceRedirection
		// ReferralMarketing
		name:             provider.name,
		CompleteProvider: provider.CompleteProvider,
	}
	rx, err := compileRegexpIfNotEmpty(provider.URLPattern)
	if err != nil {
		return nil, err
	}
	result.URLPattern = rx

	rx, err = compileRegexpIfNotEmpty(provider.Rules)
	if err != nil {
		return nil, err
	}
	result.Rules = rx

	rx, err = compileRegexpIfNotEmpty(provider.RawRules)
	if err != nil {
		return nil, err
	}
	result.RawRules = rx

	rx, err = compileRegexpIfNotEmpty(provider.Exceptions)
	if err != nil {
		return nil, err
	}
	result.Exceptions = rx

	rx, err = compileRegexpIfNotEmpty(provider.ReferralMarketing)
	if err != nil {
		return nil, err
	}
	result.ReferralMarketing = rx

	result.Redirections = make([]*regexp.Regexp, len(provider.Redirections))
	for i, redirRXStr := range provider.Redirections {
		rx, err = compileRegexpIfNotEmpty(redirRXStr)
		if err != nil {
			return nil, err
		}
		result.Redirections[i] = rx
	}

	return result, nil
}

// implements RunnableProvider
func (provider *providerCompiled) compile() (*providerCompiled, error) {
	return provider, nil
}

// implements RunnableProvider
func (provider *providerJSON) compile() (*providerCompiled, error) {
	return provider.prepare().compile()
}

// implements RunnableProvider
func (provider *providerJSON) IsCompiled() bool {
	return false
}

// Compile all the regex strings to match faster
func Compile(providers []RunnableProvider) ([]RunnableProvider, error) {
	result := make([]RunnableProvider, len(providers))
	for i, provider := range providers {
		compiled, err := provider.compile()
		if err != nil {
			return nil, err
		}
		result[i] = compiled
	}
	return result, nil
}

// Same as [DownloadSource.Download] but [Compile] the providers before returning
func (source *DownloadSource) DownloadCompiled(checkHash bool) ([]RunnableProvider, error) {
	result, err := source.Download(checkHash)
	if err != nil {
		return nil, err
	}
	return Compile(result)
}

// Same as [DownloadSource.DownloadWithCache] but [Compile] the providers before returning
func (source *DownloadSource) DownloadWithCacheCompiled(cacheFileName string, cacheMaxAgeM int, checkHash bool) ([]RunnableProvider, error) {
	result, err := source.DownloadWithCache(cacheFileName, cacheMaxAgeM, checkHash)
	if err != nil {
		return nil, err
	}
	return Compile(result)
}
