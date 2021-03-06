package flexconfig

/*
Copyright 2018-2019 The flexconfig Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"strings"
)

const (
	allowedCharacters = "abcdefghijklmnopqrstuvwxyz.-_0123456789"
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

// searchArguments iterates through an array of strings representing command
// line arguments to locate a value having the form --<argumentName>=<value>.
// If a member is found matching the pattern, the value of that argument is
// returned. If no member matches the pattern and empty string is returned.
func searchArgument(args []string, argumentName string) string {
	val := ""

	if len(argumentName) == 0 {
		return val
	}

	for _, a := range args {
		if strings.HasPrefix(a, "--"+argumentName+"=") {
			index := len(argumentName) + 3
			if len(a) == index {
				return val
			}

			return a[index:]
		}
	}

	return val
}

// conformsToKey checks whether a given string conforms to the character set
// allowed for property keys.
func conformsToKey(arg string) bool {
	for _, c := range []rune(arg) {
		if strings.Index(allowedCharacters, string(c)) < 0 {
			return false
		}
	}

	return true
}
