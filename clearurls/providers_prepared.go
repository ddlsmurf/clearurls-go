package clearurls

import (
	"fmt"
	"regexp"
)

// Intermediary provider with the full regexen strings ready for compilation
// after interpretation from the JSON.
type providerWithPreparedRegexStr struct {
	name              string
	URLPattern        string
	CompleteProvider  bool
	Rules             string
	RawRules          string
	Exceptions        string
	Redirections      []string
	ReferralMarketing string
}

// implements RunnableProvider
func (provider *providerWithPreparedRegexStr) prepare() *providerWithPreparedRegexStr {
	return provider
}

// From a JSON provider entry, create a intermediary representation with adapted
// regexen in string form
//
// implements RunnableProvider
func (provider *providerJSON) prepare() *providerWithPreparedRegexStr {
	makeCaseInsensitive := func(rxStr string) string {
		if rxStr == "" {
			return ""
		}
		return caseInsensitiveRXStrPrefix + rxStr
	}
	result := &providerWithPreparedRegexStr{
		// ForceRedirection
		// ReferralMarketing
		name:              provider.name,
		CompleteProvider:  provider.CompleteProvider,
		URLPattern:        makeCaseInsensitive(provider.URLPattern),
		Rules:             makeCaseInsensitive(regexStrForAnyOf(provider.Rules, "^", "$")),
		RawRules:          makeCaseInsensitive(regexStrForAnyOf(provider.RawRules, "", "")),
		Exceptions:        makeCaseInsensitive(regexStrForAnyOf(provider.Exceptions, "", "")),
		ReferralMarketing: makeCaseInsensitive(regexStrForAnyOf(provider.ReferralMarketing, "", "")),
		Redirections:      make([]string, len(provider.Redirections)),
	}
	for i, redirRXStr := range provider.Redirections {
		result.Redirections[i] = makeCaseInsensitive(redirRXStr)
	}
	return result
}

// From a compiled provider entry, create a intermediary representation with adapted
// regexen in string form
//
// I don't know why anyone would use this, but it's there anyway
//
// implements RunnableProvider
func (provider *providerCompiled) prepare() *providerWithPreparedRegexStr {
	safeString := func(regex *regexp.Regexp) string {
		if regex == nil {
			return ""
		}
		return regex.String()
	}
	result := &providerWithPreparedRegexStr{
		// ForceRedirection
		// ReferralMarketing
		name:              provider.name,
		CompleteProvider:  provider.CompleteProvider,
		URLPattern:        safeString(provider.URLPattern),
		Rules:             safeString(provider.Rules),
		RawRules:          safeString(provider.RawRules),
		Exceptions:        safeString(provider.Exceptions),
		ReferralMarketing: safeString(provider.ReferralMarketing),
		Redirections:      make([]string, len(provider.Redirections)),
	}
	for i, redirRXStr := range provider.Redirections {
		result.Redirections[i] = redirRXStr.String()
	}
	return result
}

// ## Prevent using intermediary providerWithPreparedRegexStr for actual matching:
// (but stay a valid RunnableProvider for facility of not splitting prepare/compile)
//
// `providerWithPreparedRegexStr` is an intermediary representation between
// either JSON or the compiled versions

// implements RunnableProvider - kinda
func (provider *providerWithPreparedRegexStr) matchURL(url string) (bool, error) {
	panic(fmt.Errorf("providerWithPreparedRegexStr can't matchURL"))
}

// implements RunnableProvider
func (provider *providerWithPreparedRegexStr) getName() string {
	return provider.name
}

// implements RunnableProvider
func (provider *providerWithPreparedRegexStr) isComplete() bool {
	return provider.CompleteProvider // I suppose this much we can :)
}

// implements RunnableProvider - kinda
func (provider *providerWithPreparedRegexStr) hasRedirect(url string) ([][]string, error) {
	panic(fmt.Errorf("providerWithPreparedRegexStr can't hasRedirect"))
}

// implements RunnableProvider - kinda
func (provider *providerWithPreparedRegexStr) applyRawRules(url string) (string, error) {
	panic(fmt.Errorf("providerWithPreparedRegexStr can't applyRawRules"))
}

// implements RunnableProvider - kinda
func (provider *providerWithPreparedRegexStr) rulesKeyFilter(key string, dontFilterReferrals bool) (bool, error) {
	panic(fmt.Errorf("providerWithPreparedRegexStr can't rulesKeyFilter"))
}
