package clearurls

// Use a `providerCompiled`, a `RunnableProvider`, to run the ClearURLs match and transform

// implements RunnableProvider
func (provider *providerCompiled) matchURL(url string) (bool, error) {
	if !provider.URLPattern.MatchString(url) {
		return false, nil
	}
	return !provider.Exceptions.MatchString(url), nil
}

// implements RunnableProvider
func (provider *providerCompiled) getName() string {
	return provider.name
}

// implements RunnableProvider
func (provider *providerCompiled) isComplete() bool {
	return provider.CompleteProvider
}

// implements RunnableProvider
func (provider *providerCompiled) hasRedirect(url string) ([][]string, error) {
	for _, redirectionRX := range provider.Redirections {
		redirMatches := redirectionRX.FindAllStringSubmatch(url, -1)
		if len(redirMatches) > 0 {
			return redirMatches, nil
		}
	}
	return nil, nil
}

// implements RunnableProvider
func (provider *providerCompiled) applyRawRules(url string) (string, error) {
	if provider.RawRules == nil {
		return url, nil
	}
	return provider.RawRules.ReplaceAllString(url, ""), nil
}

// implements RunnableProvider
func (provider *providerCompiled) rulesKeyFilter(key string, dontFilterReferrals bool) (bool, error) {
	shouldFilter := provider.Rules.MatchString(key)
	if dontFilterReferrals && shouldFilter && provider.ReferralMarketing.MatchString(key) {
		return false, nil
	}
	return shouldFilter, nil
}
