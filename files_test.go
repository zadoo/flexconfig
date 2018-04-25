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
	"io/ioutil"
	"os"
	"testing"
)

const (
	appName = "testConfig"
)

func Test_files(t *testing.T) {
	v := make(map[string]string)
	readConfigFiles(v, appName, []string{".conf"}, "test")

	if len(v) != 6 {
		t.Errorf("Unexpected number of properties: %d", len(v))
	}

	if v["test.conf.one"] != "written by yaml" {
		t.Errorf("Unexpected value: %s", v["test.conf.one"])
	}

	if v["test.conf.two"] != "written by yaml" {
		t.Errorf("Unexpected value: %s", v["test.conf.two"])
	}

	if v["test.conf.three"] != "written by json" {
		t.Errorf("Unexpected value: %s", v["test.conf.three"])
	}

	if len(v["test.conf.seven"]) > 0 {
		t.Errorf("Unexpected presence: %s", v["test.conf.seven"])
	}

	if len(v["test.conf.four"]) > 0 {
		t.Errorf("Unexpected presence: %s", v["test.conf.four"])
	}

	if v["test.section1.name"] != "section1-name" {
		t.Errorf("Unexpected value: %s", v["test.section1.name"])
	}

	if v["test.section2.name"] != "section2-name" {
		t.Errorf("Unexpected value: %s", v["test.section2.name"])
	}

	if v["test.section2.other"] != "otherValue" {
		t.Errorf("Unexpected value: %s", v["test.section2.other"])
	}
}

func Test_files_differentSuffix(t *testing.T) {
	v := make(map[string]string)
	readConfigFiles(v, appName, []string{".xyz"}, "")

	if len(v) != 1 {
		t.Errorf("Unexpected number of properties: %d", len(v))
	}

	if v["test.conf.seven"] != "written by xyz" {
		t.Errorf("Unexpected value: %s", v["test.conf.seven"])
	}
}

func Test_files_emptyDir(t *testing.T) {
	name, err := ioutil.TempDir(".", ".testTempFiles")
	if err != nil {
		t.Errorf("Can't create temporary directory")
	}

	defer os.Remove(name)

	localAppName := name[1:]

	v := make(map[string]string)
	readConfigFiles(v, localAppName, []string{".conf"}, "")

	if len(v) > 0 {
		t.Errorf("Unexpected properties found in empty directory: %d", len(v))
	}
}
