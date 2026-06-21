// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

// Package keys provides stable state-store key sanitization helpers.
package keys

import "strings"

var segmentReplacer = strings.NewReplacer(" ", "-", "/", "-", "\\", "-", ":", "-", ".", "-")

// Segment normalizes a raw key segment for state-store storage. If the input is empty or becomes empty after normalization, the fallback value is returned instead.
func Segment(value string, fallback string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return fallback
	}

	sanitized := segmentReplacer.Replace(trimmed)
	if sanitized == "" {
		return fallback
	}
	return sanitized
}
