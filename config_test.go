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
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_config_getBeforeInitialized(t *testing.T) {
	c := GetConfiguration()
	if c == nil {
		t.Errorf("Uninitialized configuration should create empty configuration")
		return
	}

	if len(configuration.appName) > 0 {
		t.Errorf("Empty configuration should not have appName set")
	}

	if len(configuration.config) > 0 {
		t.Errorf("Empty configuration should not have any properties set")
	}
}

func Test_config_topLevel(t *testing.T) {
	_, err := NewFlexibleConfiguration(ConfigurationParameters{})
	if err != nil {
		t.Errorf("Error calling NewFlexibleConfiguration")
	}

	key := "abc"
	setval := "123"
	val := Get(key)
	if len(val) > 0 {
		t.Errorf("Get of non-existent key has value: %s", val)
	}

	if Exists(key) {
		t.Errorf("Non-existent key exists")
	}

	Set(key, setval)
	if !Exists(key) {
		t.Errorf("Key does not have value")
	}

	val = Get(key)
	if val != setval {
		t.Errorf("Key does not have value set: %s", val)
	}
}

func Test_config_emptyDefault(t *testing.T) {
	os.Setenv("TEST_CONFIG_ENV", "singleWord")

	cfg, err := NewFlexibleConfiguration(ConfigurationParameters{})
	if err != nil {
		t.Errorf("Error calling NewFlexibleConfiguration")
	}

	if cfg == nil {
		t.Errorf("Initialized configuration should not be nil")
		return
	}

	c := cfg.(*flexibleConfiguration)
	if len(c.config) > 0 {
		t.Errorf("Config has unintended properties set")
	}

	if c.Exists("test.nonexistent.property") {
		t.Errorf("Reports existence of nonexistent property")
	}

	val := c.Get("test.config.env")
	if len(val) > 0 {
		t.Errorf("Unexpected value for property: %s", val)
	}

	c.Set("test.config.env", "setValue")
	if !c.Exists("test.config.env") {
		t.Errorf("Set property does not exist")
	}

	val = c.Get("test.config.env")
	if val != "setValue" {
		t.Errorf("Expecting %v, but found %v", "setValue", val)
	}

	c.Set("", "emptyValue")
	if c.Exists("") {
		t.Errorf("Empty property exists")
	}

	val = c.Get("")
	if len(val) > 0 {
		t.Errorf("Get of empty key returns value: %s", val)
	}
}

func Test_config_onlyEnv(t *testing.T) {
	os.Setenv("TEST_TEST_ENV", "singleWord")

	cfg, err := NewFlexibleConfiguration(ConfigurationParameters{
		EnvironmentVariablePrefixes: []string{"TEST_TEST_"}})
	if err != nil {
		t.Errorf("Error calling NewFlexibleConfiguration")
	}

	if cfg == nil {
		t.Errorf("Initialized configuration should not be nil")
		return
	}

	c := cfg.(*flexibleConfiguration)
	if len(c.config) != 1 {
		t.Errorf("Config size incorrect: %v", len(c.config))
	}

	if c.Exists("test.nonexistent.property") {
		t.Errorf("Reports existence of nonexistent property")
	}

	if !c.Exists("test.test.env") {
		t.Errorf("Environment property does not have value")
	}

	val := c.Get("test.test.env")
	if len(val) == 0 {
		t.Errorf("Property has no value")
	}

	if val != "singleWord" {
		t.Errorf("Expecting %v, but found %v", "singleWord", val)
	}

	if configuration == nil {
		t.Errorf("Initialized configuration not saved")
		return
	}

	ccfg := GetConfiguration()
	cc := ccfg.(*flexibleConfiguration)
	if len(cc.config) != 1 {
		t.Errorf("Config size incorrect: %v", len(cc.config))
	}

	val = cc.Get("test.test.env")
	if val != "singleWord" {
		t.Errorf("Expecting %v, but found %v", "singleWord", val)
	}

	cfg, err = NewFlexibleConfiguration(ConfigurationParameters{})
	if err != nil {
		t.Errorf("Error calling NewFlexibleConfiguration")
	}

	if cfg == nil {
		t.Errorf("Reinitialized configuration should not be nil")
		return
	}

	c = cfg.(*flexibleConfiguration)
	if len(c.config) > 0 {
		t.Errorf("Config has unintended properties set")
	}
}

func Test_config_onlyFiles(t *testing.T) {
	os.Setenv("TEST_CONFIG_ENV", "singleWord")

	cfg, err := NewFlexibleConfiguration(ConfigurationParameters{
		IniNamePrefix:   "test",
		ApplicationName: "testConfig"})
	if err != nil {
		t.Errorf("Error calling NewFlexibleConfiguration")
	}

	if cfg == nil {
		t.Errorf("Initialized configuration should not be nil")
		return
	}

	c := cfg.(*flexibleConfiguration)
	if len(c.config) != 6 {
		t.Errorf("Wrong number of properties read from files: %d", len(c.config))
	}

	if c.Exists("test.config.env") {
		t.Errorf("Environment property has value, but should not")
	}

	if !c.Exists("test.conf.three") {
		t.Errorf("JSON property does not have value")
	}

	val := c.Get("test.conf.three")
	if len(val) == 0 {
		t.Errorf("Property has no value")
	}

	if val != "written by json" {
		t.Errorf("Unexpected property value: %s", val)
	}

	val = c.Get("test.conf.one")
	if val != "written by yaml" {
		t.Errorf("Unexpected property value: %s", val)
	}

	val = c.Get("test.section1.name")
	if val != "section1-name" {
		t.Errorf("Unexpected property value: %s", val)
	}

	val = c.Get("test.section2.name")
	if val != "section2-name" {
		t.Errorf("Unexpected property value: %s", val)
	}

	val = c.Get("test.section2.other")
	if val != "otherValue" {
		t.Errorf("Unexpected property value: %s", val)
	}
}

func zTest_config_priority(t *testing.T) {
	endpointstr := os.Getenv(etcdEndpointEnvironmentVariable)
	if len(endpointstr) == 0 {
		endpointstr = defaultEtcdEndpoint
	}

	endpoints := strings.Split(endpointstr, ",")

	fcs, err := NewFlexConfigStore(FlexConfigStoreEtcd, endpoints, etcdTestPrefix)
	if err != nil {
		t.Errorf("Error creating store: %v", err)
		return
	}

	err = fcs.Set("test.conf.six", "fromStore")
	if err != nil {
		t.Errorf("Failed to set store property: %v", err)
	}

	defer fcs.Delete("test.conf.six")

	os.Setenv("TEST_CONF_TWO", "fromEnv")
	os.Setenv("TEST_CONF_FIVE", "fromEnv")
	os.Args = []string{"test", "--test.conf.five=fromCommandLine", "--test.conf.six=fromCommandLine"}

	c, err := NewFlexibleConfiguration(ConfigurationParameters{
		ApplicationName:             "testConfig",
		EnvironmentVariablePrefixes: []string{"TEST_CONF_"},
		ConfigurationStore:          fcs,
	})
	if err != nil {
		t.Errorf("Error calling NewFlexibleConfiguration")
	}

	if c == nil {
		t.Errorf("Initialized configuration should not be nil")
		return
	}

	val := c.Get("test.conf.two")
	if val != "fromEnv" {
		t.Errorf("Environment vars did not override file vars")
	}

	val = c.Get("test.conf.five")
	if val != "fromCommandLine" {
		t.Errorf("Command line vars did not override env vars: %s", val)
	}

	val = c.Get("test.conf.six")
	if val != "fromStore" {
		t.Errorf("Store var did not override command line var: %s", val)
	}

	val = c.Get("test.conf.one")
	if val != "written by yaml" {
		t.Errorf("File vars are not present")
	}

	val = c.Get("test.conf.three")
	if val != "written by json" {
		t.Errorf("File vars are not present")
	}

	cptr := c.(*flexibleConfiguration)
	if len(cptr.config) != 8 {
		t.Errorf("Wrong number of properties found: %d", len(cptr.config))
	}

	if !c.Exists("test.conf.one") {
		t.Errorf("Failed to find existing property")
	}

	if !c.Exists("test.conf.six") {
		t.Errorf("Failed to find existing property")
	}

	if c.Exists("test.conf.seven") {
		t.Errorf("Unexpected find of non-existing property")
	}

	c.Set("test.conf.seven", "fromCode")

	if !c.Exists("test.conf.seven") {
		t.Errorf("Failed to find existing property")
	}

	fcs.Delete("test.conf.seven")
}

func Test_config_newBadParms(t *testing.T) {
	os.Args = []string{}
	configuration = nil
	_, err := NewFlexibleConfiguration(ConfigurationParameters{
		ApplicationName: "@*^Vbd)(",
	})
	if err == nil {
		t.Errorf("Unexecpted success")
	}

	if err.Error() != "Application name not valid" {
		t.Errorf("Unexpected error text: %s", err)
	}

	if configuration != nil {
		t.Errorf("Unexpected configuration")
	}

	cfg := GetConfiguration()
	if cfg == nil {
		t.Errorf("Unexpected nil configuration")
	}

	if configuration == nil {
		t.Errorf("Unexpected nil configuration")
	}

	if len(configuration.config) > 0 {
		t.Errorf("Unexpected properties in config: %v", configuration.config)
	}

	configuration = nil
	_, err = NewFlexibleConfiguration(ConfigurationParameters{
		ApplicationName: "1abcdef",
	})
	if err == nil {
		t.Errorf("Unexecpted success")
	}

	configuration = nil
	_, err = NewFlexibleConfiguration(ConfigurationParameters{
		ApplicationName: "abc-123",
	})
	if err == nil {
		t.Errorf("Unexecpted success")
	}

	configuration = nil
	_, err = NewFlexibleConfiguration(ConfigurationParameters{
		EnvironmentVariablePrefixes: []string{"@*^Vbd)("},
	})
	if err == nil {
		t.Errorf("Unexecpted success")
	}

	if configuration != nil {
		t.Errorf("Unexpected configuration")
	}

	configuration = nil
	_, err = NewFlexibleConfiguration(ConfigurationParameters{
		EnvironmentVariablePrefixes: []string{"1abcdef"},
	})
	if err == nil {
		t.Errorf("Unexecpted success")
	}

	configuration = nil
	_, err = NewFlexibleConfiguration(ConfigurationParameters{
		EnvironmentVariablePrefixes: []string{"abc-123"},
	})
	if err == nil {
		t.Errorf("Unexecpted success")
	}
}

func Test_config_badPropertyNames(t *testing.T) {
	os.Args = []string{}
	c, err := NewFlexibleConfiguration(ConfigurationParameters{})
	if err != nil {
		t.Errorf("Error calling NewFlexibleConfiguration: %s", err)
	}

	key := "1abcdef"
	c.Set(key, "value")
	if c.Exists(key) {
		t.Errorf("Unexpected success setting bad property name")
	}

	key = "$*^%!(abc123"
	c.Set(key, "value")
	if c.Exists(key) {
		t.Errorf("Unexpected success setting bad property name")
	}

	key = "abc$123"
	c.Set(key, "value")
	if c.Exists(key) {
		t.Errorf("Unexpected success setting bad property name")
	}

	key = "abc-def"
	c.Set(key, "value")
	if !c.Exists(key) {
		t.Errorf("Error using '-' in property name")
	}

	key = "abc.def"
	c.Set(key, "value")
	if !c.Exists(key) {
		t.Errorf("Error using '.' in property name")
	}

	key = "abc_def"
	c.Set(key, "value")
	if !c.Exists(key) {
		t.Errorf("Error using '_' in property name")
	}
}

func Test_config_newDifferentSuffix(t *testing.T) {
	os.Args = []string{}
	cfg, err := NewFlexibleConfiguration(ConfigurationParameters{
		ApplicationName:      "testConfig",
		AcceptedFileSuffixes: []string{".xyz"},
	})
	if err != nil {
		t.Errorf("Error calling NewFlexibleConfiguration: %s", err)
	}

	c := cfg.(*flexibleConfiguration)
	if len(c.config) != 1 {
		t.Errorf("Unexpected number of properties: %d", len(c.config))
		fmt.Printf("config: %v\n", c.config)
	}

	val := Get("test.conf.seven")
	if val != "written by xyz" {
		t.Errorf("Unexpected value: %s", val)
	}
}
