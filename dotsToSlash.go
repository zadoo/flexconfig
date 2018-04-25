package flexconfig

/*
Copyright 2018 The flexconfig Authors

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
