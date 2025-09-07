package clearurls

import (
	"regexp"
	"strings"
)

const caseInsensitiveRXStrPrefix = "(?i)"

// Return true if any of the provided patterns match `needle`
func matchAnyOfRegexenStringCaseInsensitive(patterns []string, needle string) (bool, error) {
	for _, rx := range patterns {
		if match, err := regexp.MatchString(caseInsensitiveRXStrPrefix+rx, needle); err != nil || match {
			return match, err
		}
	}
	return false, nil
}

// Generate a regex string that would match any of the regex strings in `alternateRegexen`
func regexStrForAnyOf(alternateRegexen []string, prefix, suffix string) string {
	if len(alternateRegexen) == 0 {
		return ""
	}
	return prefix + "(?:(?:" + strings.Join(alternateRegexen, ")|(?:") + "))" + suffix
}
