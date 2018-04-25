package flexconfig

import (
	"os"
	"testing"
)

func Test_env_emptyPrefixList(t *testing.T) {
	v := make(map[string]string)
	readEnvVars(v, nil)

	if len(v) > 0 {
		t.Errorf("Unexpected properties found")
	}
}

func Test_env(t *testing.T) {
	os.Setenv("TEST_CONFIG_ENV_ONE", "string value")
	os.Setenv("TEST_CONFIG_ENV_TWO", "42")

	v := make(map[string]string)
	readEnvVars(v, []string{"TEST_CONFIG_ENV_"})

	if len(v) != 2 {
		t.Errorf("Unexpected properties found")
	}

	if v["test.config.env.one"] != "string value" {
		t.Errorf("Unexpected value: %s", v["test.config.env.one"])
	}

	if v["test.config.env.two"] != "42" {
		t.Errorf("Unexpected value: %s", v["test.config.env.two"])
	}

	os.Unsetenv("TEST_CONFIG_ENV_ONE")
	os.Unsetenv("TEST_CONFIG_ENV_TWO")
}

func Test_env_someAccepted(t *testing.T) {
	os.Setenv("TEST_CONFIG_ENV_ONE", "string value")
	os.Setenv("DEBUG_CONFIG_ENV_TWO", "42")

	v := make(map[string]string)
	readEnvVars(v, []string{"TEST_CONFIG_ENV_"})

	if len(v) != 1 {
		t.Errorf("Unexpected properties found")
	}

	if v["test.config.env.one"] != "string value" {
		t.Errorf("Unexpected value: %s", v["test.config.env.one"])
	}

	os.Unsetenv("TEST_CONFIG_ENV_ONE")
	os.Unsetenv("DEBUG_CONFIG_ENV_TWO")
}

func Test_evaluate_nilPrefix(t *testing.T) {
	v := make(map[string]string)
	evaluateEnvVar(v, nil, "TEST_IT=42")

	if len(v) > 0 {
		t.Errorf("Unexpected properties found")
	}
}

func Test_evaluate_emptyPrefix(t *testing.T) {
	v := make(map[string]string)
	evaluateEnvVar(v, []string{}, "TEST_IT=42")

	if len(v) > 0 {
		t.Errorf("Unexpected properties found")
	}
}

func Test_evaluate_otherPrefix(t *testing.T) {
	v := make(map[string]string)
	evaluateEnvVar(v, []string{"DEBUG_", "OTHER_", "TEST_CONFIG_"}, "TEST_IT=42")

	if len(v) > 0 {
		t.Errorf("Unexpected properties found")
	}
}

func Test_evaluate_found(t *testing.T) {
	v := make(map[string]string)
	evaluateEnvVar(v, []string{"DEBUG_", "OTHER_", "TEST_"}, "TEST_IT=42")

	if len(v) != 1 {
		t.Errorf("Unexpected properties found")
	}

	if v["test.it"] != "42" {
		t.Errorf("Unexpected value: %s", v["test.it"])
	}
}
