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
	"testing"
)

func Test_dts(t *testing.T) {
	s := dotsToSlash("")
	if len(s) > 0 {
		t.Errorf("Unexpected value: %s", s)
	}

	s = dotsToSlash("a.b.c")
	if s != "/a/b/c" {
		t.Errorf("Unexpected value: %s", s)
	}

	s = dotsToSlash(".a.b.c.")
	if s != "/a/b/c/" {
		t.Errorf("Unexpected value: %s", s)
	}
}

func Test_std(t *testing.T) {
	s := slashToDots("")
	if len(s) > 0 {
		t.Errorf("Unexpected value: %s", s)
	}

	s = slashToDots("a/b/c")
	if s != "a.b.c" {
		t.Errorf("Unexpected value: %s", s)
	}

	s = slashToDots("/a/b/c/")
	if s != ".a.b.c." {
		t.Errorf("Unexpected value: %s", s)
	}
}
