package flexconfig

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
