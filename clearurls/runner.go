package clearurls

// Implement the ClearURLs matching and transformation algorithm,
// using `RunnableProvider`s repeatedly to remove tracking parameters

import (
	"fmt"
	"net/url"
)

// If a redirect in the provider matches, return that url
func getRedirect(provider RunnableProvider, urlToSearch string) (string, error) {
	redirMatches, err := provider.hasRedirect(urlToSearch)
	if err != nil || redirMatches == nil {
		return "", err
	}
	if len(redirMatches) > 1 || len(redirMatches[0]) < 2 {
		return "", fmt.Errorf("URL %q: Provider %v: More than one redirection match, or no groups in %#v", urlToSearch, provider, redirMatches)
	}
	return url.QueryUnescape(redirMatches[0][1])
}

func runProviderRuleOnValues(provider RunnableProvider, values url.Values, dontFilterReferrals bool) (string, error) {
	keysToDelete := make([]string, 0, 3)
	for key := range values {
		shouldFilter, err := provider.rulesKeyFilter(key, dontFilterReferrals)
		if err != nil {
			return "", err
		} else if shouldFilter {
			keysToDelete = append(keysToDelete, key)
		}
	}
	for _, keyToDelete := range keysToDelete {
		values.Del(keyToDelete)
	}
	return values.Encode(), nil
}

// Run on query then fragments. Order of Addon is not respected here, it does foreach rule { foreach [query, fragments] { apply() } }
func runProviderRule(provider RunnableProvider, parsedURL *url.URL, dontFilterReferrals bool) error {
	queryValues, err := runProviderRuleOnValues(provider, parsedURL.Query(), dontFilterReferrals)
	if err != nil {
		return err
	}
	parsedURL.RawQuery = queryValues
	fragmentValues, err := url.ParseQuery(parsedURL.Fragment)
	if err != nil {
		return err
	}
	if len(fragmentValues) > 0 {
		fragStr, err := runProviderRuleOnValues(provider, fragmentValues, dontFilterReferrals)
		if err != nil {
			return err
		}
		parsedURL.Fragment = fragStr
	}
	return nil
}

// Go through every provider (except if one returns a redirection), updating the URL
func runProviders(providers []RunnableProvider, runningURL string, dontFilterReferrals bool) (string, error) {
	// Equivalent to _cleaning @ https://github.com/ClearURLs/Addon/blob/master/core_js/pureCleaning.js#L43
	for _, provider := range providers {
		matched, err := provider.matchURL(runningURL)
		if err != nil {
			return "", err
		}
		if !matched {
			continue
		}

		if redirectionURL, err := getRedirect(provider, runningURL); err != nil || redirectionURL != "" {
			return redirectionURL, err
		}

		if provider.isComplete() {
			// Addon code contradicts doc at https://docs.clearurls.xyz/1.27.3/specs/rules/#completeprovider - redirections are processed before
			continue
		}

		parsedURL, err := url.Parse(runningURL)
		if err != nil {
			return "", err
		}
		if err := runProviderRule(provider, parsedURL, dontFilterReferrals); err != nil {
			return "", err
		}

		runningURL = parsedURL.String()
	}
	return runningURL, nil
}

// Clean the provided `url` from tracking etc as per [ClearURLs].
// Prioritises the [source] when it disagrees with the documentation.
// If `keepMarketingReferrals` is `true`, parameters matching `referralMarketing`
// regexen will be left included.
//
// To obtain a list of [RunnableProvider]: see [GetProvidersFromSourceArgument] for examples
//
// Example:
//
//	providers, err := clearurls.MustHaveHardcodedProviders()
//	// if err != nil ....
//	clearedURL, err := clearurls.ClearURL(providers, "http://example.com?eviltrackytracktrack=true", false)
//	// if err != nil ....
//
// [ClearURLs]: https://docs.clearurls.xyz/1.27.3/
// [source]: https://github.com/ClearURLs/Addon
func ClearURL(providers []RunnableProvider, url string, keepMarketingReferrals bool) (string, error) {
	// Equivalent to pureCleaning @ https://github.com/ClearURLs/Addon/blob/master/core_js/pureCleaning.js#L28
	var prev string
	for changed := true; changed; changed = prev != url {
		prev = url
		var err error
		url, err = runProviders(providers, url, keepMarketingReferrals)
		if err != nil {
			return "", err
		}
	}
	return url, nil
}
