// Package realtimeauth issues and validates short-lived realtime access tickets.
package realtimeauth

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

// ErrOriginDenied indicates that the websocket request origin is not allowlisted.
var ErrOriginDenied = errors.New("websocket origin denied")

// ValidateOrigin validates that the provided WebSocket Origin header matches an entry in the allowlist. It returns nil if a match is found, and ErrOriginDenied otherwise.
func ValidateOrigin(requestOrigin string, allowedOrigins []string) error {
	normalizedOrigin, ok := normalizeOrigin(requestOrigin)
	if !ok {
		return ErrOriginDenied
	}
	for _, allowed := range allowedOrigins {
		normalizedAllowed, allowedOK := normalizeOrigin(allowed)
		if allowedOK && strings.EqualFold(normalizedOrigin, normalizedAllowed) {
			return nil
		}
	}
	return ErrOriginDenied
}

func normalizeOrigin(raw string) (string, bool) {
	origin := strings.TrimSpace(raw)
	if origin == "" {
		return "", false
	}
	parsedOrigin, err := url.Parse(origin)
	if err != nil || parsedOrigin.Scheme == "" || parsedOrigin.Host == "" {
		return "", false
	}
	if parsedOrigin.Path != "" || parsedOrigin.RawQuery != "" || parsedOrigin.Fragment != "" {
		return "", false
	}

	scheme := strings.ToLower(parsedOrigin.Scheme)
	host := strings.ToLower(parsedOrigin.Hostname())
	if host == "" {
		return "", false
	}

	port := parsedOrigin.Port()
	switch {
	case port != "":
		return scheme + "://" + net.JoinHostPort(host, port), true
	case scheme == "http":
		return scheme + "://" + host + ":80", true
	case scheme == "https":
		return scheme + "://" + host + ":443", true
	default:
		return scheme + "://" + host, true
	}
}
