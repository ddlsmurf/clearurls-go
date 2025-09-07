package clearurls

// Match and transform based on a `providerJSON`, a `RunnableProvider`

import "regexp"

// implements RunnableProvider
func (provider *providerJSON) matchURL(url string) (bool, error) {
	matches, err := regexp.MatchString(caseInsensitiveRXStrPrefix+provider.URLPattern, url)
	if err != nil || !matches {
		return false, err
	}
	isException, err := matchAnyOfRegexenStringCaseInsensitive(provider.Exceptions, url)
	if err != nil || isException {
		return false, err
	}
	return true, nil
}

// implements RunnableProvider
func (provider *providerJSON) getName() string {
	return provider.name
}

// implements RunnableProvider
func (provider *providerJSON) isComplete() bool {
	return provider.CompleteProvider
}

// implements RunnableProvider
func (provider *providerJSON) hasRedirect(url string) ([][]string, error) {
	for _, redirectionRXStr := range provider.Redirections {
		redirectionRX, errCompilingRedirRx := regexp.Compile(caseInsensitiveRXStrPrefix + redirectionRXStr)
		if errCompilingRedirRx != nil {
			return nil, errCompilingRedirRx
		}
		redirMatches := redirectionRX.FindAllStringSubmatch(url, -1)
		if len(redirMatches) > 0 {
			return redirMatches, nil
		}
	}
	return nil, nil
}

// implements RunnableProvider
func (provider *providerJSON) applyRawRules(url string) (string, error) {
	for _, ruleRXStr := range provider.RawRules {
		regex, err := regexp.Compile(caseInsensitiveRXStrPrefix + ruleRXStr)
		if err != nil {
			return "", err
		}
		url = regex.ReplaceAllString(url, "")
	}
	return url, nil
}

// implements RunnableProvider
func (provider *providerJSON) rulesKeyFilter(key string, dontFilterReferrals bool) (bool, error) {
	matchAnyRx := func(rulesRXStr []string, key string) (bool, error) {
		for _, ruleRXStr := range rulesRXStr {
			matches, err := regexp.MatchString(caseInsensitiveRXStrPrefix+"^"+ruleRXStr+"$", key)
			if matches || err != nil {
				return matches, err
			}
		}
		return false, nil
	}
	shouldFilter, err := matchAnyRx(provider.Rules, key)
	if err != nil {
		return false, err
	}
	if shouldFilter && dontFilterReferrals {
		shouldIncludeAnyway, err := matchAnyRx(provider.ReferralMarketing, key)
		if shouldIncludeAnyway || err != nil {
			return !shouldIncludeAnyway, err
		}
	}
	return shouldFilter, nil
}
