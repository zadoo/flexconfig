package flexconfig

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
