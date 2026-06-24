package realtime

import "strings"

// NormalizeTopic trims one topic and returns an empty string for invalid blank values.
func NormalizeTopic(topic string) string {
	return strings.TrimSpace(topic)
}
