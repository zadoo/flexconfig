package flexconfig

import (
	"strings"
)

// readCommandLineArgs iterates through an array of strings representiing
// command line arguments, selecting which strings represent property
// definitions, and setting the selected properties in the specified
// configuration map.
func readCommandLineArgs(vars map[string]string, args []string) {
	for i, arg := range args {
		if i == 0 {
			// In command line args, the first arg is the command
			continue
		}

		if arg == "--" {
			// An argument of -- indicates all other arguments
			// are to skip processing
			break
		}

		if strings.HasPrefix(arg, "--") && strings.Index(arg, "=") > 0 {
			pair := strings.Split(arg[2:], "=")
			val := pair[1]
			if val[0:1] == "\"" && val[len(val)-1:] == "\"" {
				val = val[1 : len(val)-1]
			} else if val[0:1] == "'" && val[len(val)-1:] == "'" {
				val = val[1 : len(val)-1]
			}

			if conformsToKey(pair[0]) {
				vars[pair[0]] = val
			}
		}
	}
}

// conformsToKey checks whether a given string conforms to the character set
// allowed for property keys.
func conformsToKey(arg string) bool {
	for _, c := range []rune(arg) {
		if strings.Index("abcdefghijklmnopqrstuvwxyz.-_0123456789", string(c)) < 0 {
			return false
		}
	}

	return true
}
