package server

import (
	"net/url"
	"strings"
)

// consentContentSecurityPolicy allows the validated OAuth callback as well as
// the local consent POST. Browsers can apply form-action to POST redirects too.
func consentContentSecurityPolicy(redirectURI string) string {
	formAction := "'self'"
	if u, err := url.Parse(redirectURI); err == nil {
		var source string
		switch u.Scheme {
		case "http", "https":
			if u.Host != "" {
				source = u.Scheme + "://" + u.Host
			}
		case "cstars":
			source = "cstars:"
		}
		// Never let URI punctuation introduce another CSP directive, policy,
		// or wildcard. Invalid sources leave form submissions restricted to self.
		if source != "" && !strings.ContainsAny(source, " \t\r\n;,'\"\\*") {
			formAction += " " + source
		}
	}
	return "default-src 'none'; style-src 'unsafe-inline'; form-action " + formAction + "; frame-ancestors 'none'"
}
