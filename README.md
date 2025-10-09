This is an approximate implementation of [CleanURLs](https://docs.clearurls.xyz/1.27.3/) in Go.

It allows you to run CleanURLs rules to clean up a URL from tracking URL parameters and similar.

It can get those rules by downloading them (optional local cache), or by hardcoding them using
`go generate` into the library's source code directly.
This is done in the [published tags](https://github.com/ddlsmurf/clearurls-go/tags) of this package.

## CLI use

On mac, replace URL in clipboard with a cleaned URL:

```sh
go run github.com/ddlsmurf/clearurls-go/tools/cleanurls@latest clean hardcoded "$(pbpaste)" | pbcopy
```

## Go package `clearurls/`

Overall description here: [clearurls/package_doc.go](clearurls/package_doc.go).

### Example

```go
import "github.com/ddlsmurf/clearurls-go/clearurls"
// Get providers
providers, err := clearurls.SourceGitHub.DownloadCompiled(true)
providers, err := clearurls.SourceGitLab.DownloadWithCacheCompiled("filename", 60, true)
providers, err := clearurls.HardcodedProviders()
providers, err := clearurls.GetProvidersFromSourceArgument("github")
if err != nil {
  panic(err)
}
fmt.Printf("Loaded %d providers (compiled: %v)\n", len(providers), (providers[0].IsCompiled())

// Clean a URL
clearurls.ClearURL(providers, "http://example.com/", false)

```

### Generate hardcoded source file

Run `go generate github.com/ddlsmurf/clearurls-go/clearurls` in this repository. This will
create the file `clearurls/providers_hardcoded_data.go`.

## Tool `tools/cleanurls/`

Run the CleanURL process against URLs in the command line, or line-by-line over stdin.

# Disclaimer

I have not found any good test data for this, and the CleanURLs source is not always in
agreement with my interpretation of the documentation. So this result is probably not
production ready.
