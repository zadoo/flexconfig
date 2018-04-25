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

func Test_args_noArgs(t *testing.T) {
	v := make(map[string]string)
	readCommandLineArgs(v, []string{"test"})

	if len(v) > 0 {
		t.Errorf("Unexpected properties found")
	}
}

func Test_args_singleDash(t *testing.T) {
	v := make(map[string]string)
	readCommandLineArgs(v, []string{"test", "-v"})

	if len(v) > 0 {
		t.Errorf("Unexpected properties found")
	}
}

func Test_args_noEqual(t *testing.T) {
	v := make(map[string]string)
	readCommandLineArgs(v, []string{"test", "--test.config.one"})

	if len(v) > 0 {
		t.Errorf("Unexpected properties found")
	}
}

func Test_args_found(t *testing.T) {
	v := make(map[string]string)
	readCommandLineArgs(v, []string{"test", "-v", "--test.config.one=singleWord", "-z"})

	if len(v) != 1 {
		t.Errorf("Wrong number of properties")
	}

	if v["test.config.one"] != "singleWord" {
		t.Errorf("Unexpected value for property: %s", v["test.config.one"])
	}
}

func Test_args_wrongCharacters(t *testing.T) {
	v := make(map[string]string)
	readCommandLineArgs(v, []string{"test", "--Test.Config.One=singleWord"})

	if len(v) > 0 {
		t.Errorf("Unexpected properties found")
	}
}

func Test_args_foundLongString(t *testing.T) {
	v := make(map[string]string)
	readCommandLineArgs(v, []string{"test", "--test.config.one=\"String Value\""})

	if len(v) != 1 {
		t.Errorf("Wrong number of properties")
	}

	if v["test.config.one"] != "String Value" {
		t.Errorf("Unexpected value for property: %s", v["test.config.one"])
	}
}

func Test_args_foundTwo(t *testing.T) {
	v := make(map[string]string)
	readCommandLineArgs(v, []string{"test", "--test.config.one='String Value'", "--test.config.two=42"})

	if len(v) != 2 {
		t.Errorf("Wrong number of properties")
	}

	if v["test.config.one"] != "String Value" {
		t.Errorf("Unexpected value for property: %s", v["test.config.one"])
	}

	if v["test.config.two"] != "42" {
		t.Errorf("Unexpected value for property: %s", v["test.config.two"])
	}
}

func Test_args_doubleMinusStops(t *testing.T) {
	v := make(map[string]string)
	readCommandLineArgs(v, []string{"test", "--test.config.one=\"String Value\"", "--", "--test.config.two=42"})

	if len(v) != 1 {
		t.Errorf("Wrong number of properties")
	}

	if v["test.config.one"] != "String Value" {
		t.Errorf("Unexpected value for property: %s", v["test.config.one"])
	}
}
