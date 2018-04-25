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

func Test_yaml_flat(t *testing.T) {
	v := make(map[string]string)
	contents := "flat.string: \"string value\"\n" +
		"flat.bool: true\n" +
		"flat.int: 42\n" +
		"flat.float: 3.14159\n"

	err := parseYaml(v, contents)
	if err != nil {
		t.Errorf("Unexpected failure: %v", err)
	}

	if v["flat.string"] != "string value" {
		t.Errorf("Missing string value in result")
	}

	if v["flat.bool"] != "true" {
		t.Errorf("Missing bool value in result")
	}

	if v["flat.int"] != "42" {
		t.Errorf("Missing int value in result")
	}

	if v["flat.float"] != "3.14159" {
		t.Errorf("Missing float value in result")
	}
}

func Test_yaml_embeddedStruct(t *testing.T) {
	v := make(map[string]string)
	contents := "zero: false\n" +
		"one:\n" +
		"  string: singleWord\n" +
		"  bool: true\n" +
		"two: -4.296e-99\n" +
		"three:\n" +
		"  int: 42\n" +
		"  float: 3.14159\n" +
		"four: 1234567890\n"

	err := parseYaml(v, contents)
	if err != nil {
		t.Errorf("Unexpected failure: %v", err)
	}

	if v["zero"] != "false" {
		t.Errorf("Missing bool value in result")
	}

	if v["one.string"] != "singleWord" {
		t.Errorf("Missing string value in result")
	}

	if v["one.bool"] != "true" {
		t.Errorf("Missing bool value in result")
	}

	if v["two"] != "-4.296e-99" {
		t.Errorf("Missing bool value in result")
	}

	if v["three.int"] != "42" {
		t.Errorf("Missing int value in result")
	}

	if v["three.float"] != "3.14159" {
		t.Errorf("Missing float value in result")
	}

	if v["four"] != "1234567890" {
		t.Errorf("Missing bool value in result")
	}
}

func Test_yaml_semiFlat(t *testing.T) {
	v := make(map[string]string)
	contents := "zero: false\n" +
		"one.with.several.fields:\n" +
		"  string: \"string value\"\n" +
		"  bool: true\n" +
		"two: -4.296e-99\n"

	err := parseYaml(v, contents)
	if err != nil {
		t.Errorf("Unexpected failure: %v", err)
	}

	if v["zero"] != "false" {
		t.Errorf("Missing bool value in result")
	}

	if v["one.with.several.fields.string"] != "string value" {
		t.Errorf("Missing string value in result")
	}

	if v["one.with.several.fields.bool"] != "true" {
		t.Errorf("Missing bool value in result")
	}

	if v["two"] != "-4.296e-99" {
		t.Errorf("Missing bool value in result")
	}
}

func Test_yaml_simpleArray(t *testing.T) {
	v := make(map[string]string)
	contents := "zero: false\n" +
		"one:\n" +
		"  - 11\n" +
		"  - 22\n" +
		"two: 42\n"

	err := parseYaml(v, contents)
	if err != nil {
		t.Errorf("Unexpected failure: %v", err)
	}

	if v["zero"] != "false" {
		t.Errorf("Missing bool value in result")
	}

	if v["one.0"] != "11" {
		t.Errorf("Missing array value in result")
	}

	if v["one.1"] != "22" {
		t.Errorf("Missing array value in result")
	}

	if v["two"] != "42" {
		t.Errorf("Missing int value in result")
	}
}

func Test_yaml_arrayOfStruct(t *testing.T) {
	v := make(map[string]string)
	contents := "zero: false\n" +
		"one:\n" +
		"  - a: 11\n" +
		"    b: 22\n" +
		"  - a: 33\n" +
		"    b: 44\n" +
		"two: 42\n"

	err := parseYaml(v, contents)
	if err != nil {
		t.Errorf("Unexpected failure: %v", err)
	}

	if v["zero"] != "false" {
		t.Errorf("Missing bool value in result")
	}

	if v["one.0.a"] != "11" {
		t.Errorf("Missing array value in result")
	}

	if v["one.0.b"] != "22" {
		t.Errorf("Missing array value in result")
	}

	if v["one.1.a"] != "33" {
		t.Errorf("Missing array value in result")
	}

	if v["one.1.b"] != "44" {
		t.Errorf("Missing array value in result")
	}

	if v["two"] != "42" {
		t.Errorf("Missing int value in result")
	}
}

func Test_yaml_arrayOfArray(t *testing.T) {
	v := make(map[string]string)
	contents := "zero: false\n" +
		"one:\n" +
		"  - \n" +
		"    - 11\n" +
		"    - 22\n" +
		"  - \n" +
		"    - 33\n" +
		"    - 44\n" +
		"two: 42\n"

	err := parseYaml(v, contents)
	if err != nil {
		t.Errorf("Unexpected failure: %v", err)
	}

	if v["zero"] != "false" {
		t.Errorf("Missing bool value in result")
	}

	if v["one.0.0"] != "11" {
		t.Errorf("Missing array value in result")
	}

	if v["one.0.1"] != "22" {
		t.Errorf("Missing array value in result")
	}

	if v["one.1.0"] != "33" {
		t.Errorf("Missing array value in result")
	}

	if v["one.1.1"] != "44" {
		t.Errorf("Missing array value in result")
	}

	if v["two"] != "42" {
		t.Errorf("Missing int value in result")
	}
}

func Test_yaml_json(t *testing.T) {
	v := make(map[string]string)
	contents := "{" +
		"  \"zero\": false," +
		"  \"one\": [" +
		"    {" +
		"      \"a\": 11," +
		"      \"b\": 22" +
		"    }," +
		"    {" +
		"      \"a\": 33," +
		"      \"b\": 44" +
		"    }" +
		"  ]," +
		"  \"test.conf.two\": 42" +
		"}"

	err := parseYaml(v, contents)
	if err != nil {
		t.Errorf("Unexpected failure: %v", err)
	}

	if v["zero"] != "false" {
		t.Errorf("Missing bool value in result")
	}

	if v["one.0.a"] != "11" {
		t.Errorf("Missing array value in result")
	}

	if v["one.0.b"] != "22" {
		t.Errorf("Missing array value in result")
	}

	if v["one.1.a"] != "33" {
		t.Errorf("Missing array value in result")
	}

	if v["one.1.b"] != "44" {
		t.Errorf("Missing array value in result")
	}

	if v["test.conf.two"] != "42" {
		t.Errorf("Missing int value in result")
	}
}

func Test_notJsonOrYaml(t *testing.T) {
	v := make(map[string]string)
	contents := "[section]\n" +
		"key=value\n"

	err := parseYaml(v, contents)
	if err == nil {
		t.Errorf("Unexpected success")
	}
}
