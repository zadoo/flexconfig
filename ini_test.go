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

func Test_ini(t *testing.T) {
	prefix := "my.prefix"
	v := make(map[string]string)
	content := "defaultName=defaultValue\n" +
		"; this is a comment\n" +
		"[sectionA]\n" +
		"name=sectionA-name\n" +
		"[sectionB]\n" +
		"name=sectionB-name\n" +
		"other=other-value\n"

	err := parseIniFile(v, prefix, content)
	if err != nil {
		t.Errorf("Error calling parseIniFile: %v", err)
	}

	if len(v) != 3 {
		t.Errorf("Unexpected number of proeprties found: %d", len(v))
	}

	if v[prefix+".sectiona.name"] != "sectionA-name" {
		t.Errorf("Unexpected property value: %s", v[prefix+".sectiona.name"])
	}

	if v[prefix+".sectionb.name"] != "sectionB-name" {
		t.Errorf("Unexpected property value: %s", v[prefix+".sectionb.name"])
	}

	if v[prefix+".sectionb.other"] != "other-value" {
		t.Errorf("Unexpected property value: %s", v[prefix+".sectionb.other"])
	}
}
