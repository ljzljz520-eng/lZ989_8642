package model

import "strings"

func NormalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	seen := map[string]bool{}
	for _, t := range tags {
		t = strings.TrimSpace(strings.ToLower(t))
		if t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}
func IsTerminal(status string) bool { return status == StateArchived || status == StateStopped }
func ValidTransition(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case "":
		return to == "processed" || to == StateStopped
	case "processed":
		return to == StateArchived
	case StateRunning:
		return to == StateStopped
	}
	return false
}
