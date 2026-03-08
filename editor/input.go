package editor

import "strings"

func ParseOpts(opts []string) map[string]string {
	result := make(map[string]string)
	for _, opt := range opts {
		if strings.HasPrefix(opt, "-") {
			key := strings.TrimPrefix(opt, "-")
			if idx := strings.Index(key, "="); idx > 0 {
				result[key[:idx]] = key[idx+1:]
			} else {
				result[key] = ""
			}
		}
	}
	return result
}
