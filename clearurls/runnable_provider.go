package clearurls

// A rule from ClearURLs software than can be applied, in their lingo a Provider
type RunnableProvider interface {
	// Match semantics of matchURL @ https://github.com/ClearURLs/Addon/blob/master/clearurls.js#L404
	matchURL(url string) (bool, error)
	// Return original (unique) name of this entry
	getName() string
	// Return `completeProvider` field of provider
	isComplete() bool
	// If `redirections` field matches, return match list for validation outside
	hasRedirect(url string) ([][]string, error)
	// Apply `rawRules`
	applyRawRules(url string) (string, error)
	// Run rules (will be called for each key in query and fragment `url.Values`)
	// Return `true` if that `key` should be removed.
	rulesKeyFilter(key string, dontFilterReferrals bool) (bool, error)

	// Get a version of this provider with all regexen as ready strings for go's regexp package
	prepare() *providerWithPreparedRegexStr

	// Get a version of this provider with all regexen compiled
	compile() (*providerCompiled, error)

	// `true` if this instance is only regexen
	IsCompiled() bool
}
