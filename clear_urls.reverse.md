
/*
  My understanding of the logic in https://github.com/ClearURLs/Addon on which this code is based

  pureCleaning(url, quiet = false) // https://github.com/ClearURLs/Addon/blob/master/core_js/pureCleaning.js#L28
    - keep calling _cleaning(url) until result stops changing
  _cleaning(url, quiet = false) // https://github.com/ClearURLs/Addon/blob/master/core_js/pureCleaning.js#L43
    - var cleanURL = url
    - foreach provider
      - if providers[i].matchURL(cleanURL)
        - result = removeFieldsFormURL(providers[i], cleanURL, quiet)
        - cleanURL = result.url
        - if has redirection
          - return result.url
    - return cleanURL

    this.matchURL = function (url) // https://github.com/ClearURLs/Addon/blob/master/clearurls.js#L404
      - urlPattern.test(url) && !(this.matchException(url));
    removeFieldsFormURL(provider, pureUrl, quiet = false, request = null) // https://github.com/ClearURLs/Addon/blob/master/clearurls.js#L40
      - if provider.getRedirection(url) // yes, processed before `isCaneling`
        - return decoded URL from match with redirection=true
      - if provider.isCaneling() // isCanceling == isCompleteProvider
        - return url and cancel=true

      - foreach rawRules
        - url = url.replace(/rawRule/gi, "")
        - if changed set changes=true
      - foreach rule
        - foreach query then fragment KV pairs
          - if key matches ^ + rule + $
            - delete it, changes = true
      - rebuild url based on new query&fragment KV pairs
*/