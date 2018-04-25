package flexconfig

import (
	"os"
	"strings"
)

// readEnvVars iterates through all environment variables, selecting those with
// the one of the specified prefixes, and adds properties having a name
// converted to canonical format, and a value of the environment variable.
func readEnvVars(vars map[string]string, prefixes []string) {
	envs := os.Environ()
	for _, e := range envs {
		evaluateEnvVar(vars, prefixes, e)
	}
}

// evaluateEnvVar checks a single environment variable against accepted
// prefixes to decide if a configuramtion property should be added.
func evaluateEnvVar(vars map[string]string, prefixes []string, envvar string) {
	pair := strings.Split(envvar, "=")
	for _, prefix := range prefixes {
		if strings.HasPrefix(pair[0], prefix) {
			key := transformEnvName(pair[0])
			vars[key] = pair[1]
			break
		}
	}
}

// transformEnvName converts an environment variable name into the canonical
// configuration proeprty key form.
func transformEnvName(envName string) string {
	key := strings.ToLower(envName)
	key = strings.Replace(key, "_", ".", -1)
	return key
}
