package flexconfig

import (
	"strings"
)

// dotsToSlash converts a canonical key name from using dots as field name
// separators to using slashes as field name separators and assuring that
// an initial slash is present.
func dotsToSlash(key string) string {
	if len(key) == 0 {
		return key
	}

	name := strings.Replace(key, ".", "/", -1)
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}

	return name
}

// slashToDots converts a key name using slashes to separate fields to the
// cannonical key name using dots as field name separators.
func slashToDots(key string) string {
	if len(key) == 0 {
		return key
	}

	name := strings.Replace(key, "/", ".", -1)

	return name
}
