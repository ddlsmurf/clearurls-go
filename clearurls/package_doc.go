// Library to download and apply rules from [ClearURLs] to remove tracking parameters
// from common URL formats.
//
// In cases where the [source] doesn't match the documentation, efforts were made to
// reproduce the [source]'s behaviour.
//
//  1. The [RunnableProvider] interface represents one entry in the ClearURLs JSON that can
//     "ran", i.e. tested to match a given url, then transform it. To obtain them:
//
//     - Download from distributed source either on github or gitlab (see [DownloadSource.DownloadCompiled])
//
//     - Downloaded with a local cache (see [DownloadSource.DownloadWithCacheCompiled])
//
//     - Obtain by parsing a descriptive string, convenient for configurations (see [GetProvidersFromSourceArgument])
//
//     - If `go generate` was ran in this package, it includes a hardcoded version (see [clearurls.MustHaveHardcodedProviders])
//
//  2. For each URL to clean, call [clearurls.ClearURL]. If the result is an empty string and no error,
//     the URL is just completely blocked.
//
// [ClearURLs]: https://docs.clearurls.xyz/1.27.3/
// [source]: https://github.com/ClearURLs/Addon
package clearurls
