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
	"errors"
	"os"
	"strings"
	"unicode"
)

const (
	defaultAppName             = ""
	defaultConfigurationSuffix = ".conf"
)

var (
	defaultEnvironmentVariablePrefixes []string
)

var (
	ErrParmNameNotValid      = errors.New("Application name not valid")
	ErrParmEnvPrefixNotValid = errors.New("Environment variable prefix not valid")
)

// Config is the interface used to interact with a FlexibleConfiguration and
// is returned by NewFlexibleConfiguration and GetConfiguration.
type Config interface {
	// Exists returns whether the specified key has a non empty value in
	// the configuration.
	Exists(key string) bool

	// Get returns the value of the specified property from the
	// configuration.
	Get(key string) string

	// Set creates or modifies the specified property with the specified
	// value.
	Set(key, val string)
}

// flexibleConfiguration is the handle used to interact with a configuration.
type flexibleConfiguration struct {
	appName string
	store   FlexConfigStore
	config  map[string]string
}

// ConfigurationParameters specifies how a Config should be initialized.
//
// ApplicationName specifies the name used to find directories that may
// contain configuration files. If ApplicationName is empty (this is the
// default), no search for configuration files is done. If ApplicationName
// has a value, directories are searched in the following order, where
// directories later in the list that define properties will override the
// same property found in a directory listed earlier in the list:
//     /usr/local/etc/<name>
//     /opt/etc/<name>
//     /etc/opt/<name>
//     /etc/<name>
//     $HOME/.<name>
//     ./.<name>
// where <name> is the value of ApplicationName.
//
// EnvironmentVariablePrefixes is a list of prefixes used to determine
// which environment variables should be added to the configuration. Environment
// variable names have a canonical conversion to configuration property names:
// convert to all lower case and change all '_' to '.'. The default value
// for this field is nil, which means environment variables will not be included
// in the configuration.
//
// AcceptedFileSuffixes indicates which files in the directories mentioned
// for ApplicationName will be read to find configuration properties. If
// AcceptedFileSuffixes is nil or empty, files with the suffix ".conf" will
// be read.
//
// IniNamePrefix indicates the property name prefix that should be used for
// properties that are created from reading an INI file. The default is an
// empty string which will result in property names consisting of only the
// section and key names in the INI file.
//
// ConfigurationStore is an interface to a configuration store. When it is
// non-nil all interactions with the configuration will consult with the
// configuration store before asking the in-memory store resulting from
// reading configuration files, environment variables, and command line
// parameters.
type ConfigurationParameters struct {
	ApplicationName             string
	EnvironmentVariablePrefixes []string
	AcceptedFileSuffixes        []string
	IniNamePrefix               string
	ConfigurationStore          FlexConfigStore
}

// configuration is a singleton holding the current static configuration.
var configuration *flexibleConfiguration

// GetConfiguration returns a previously initialized configuration. If there
// is no existing configuration, a new configuration will be created by calling
// NewFlexibleConfiguration with an empty ConfigurationParameters structure.
func GetConfiguration() Config {
	if configuration == nil {
		_, err := NewFlexibleConfiguration(ConfigurationParameters{})
		if err != nil {
			configuration = new(flexibleConfiguration)
		}
	}

	return configuration
}

func NewFlexibleConfiguration(
	parameters ConfigurationParameters) (Config, error) {
	// Note: Creating a new FlexibleConfiguration will overwrite any
	// existing global configuration.

	if !nameIsValid(parameters.ApplicationName) {
		return nil, ErrParmNameNotValid
	}

	for _, ep := range parameters.EnvironmentVariablePrefixes {
		if !nameIsValid(ep) {
			return nil, ErrParmEnvPrefixNotValid
		}
	}

	if parameters.AcceptedFileSuffixes == nil ||
		len(parameters.AcceptedFileSuffixes) == 0 {
		parameters.AcceptedFileSuffixes = []string{defaultConfigurationSuffix}
	}

	configuration = new(flexibleConfiguration)
	configuration.appName = parameters.ApplicationName
	configuration.store = parameters.ConfigurationStore
	configuration.config = configuration.readConfig(parameters)

	return configuration, nil
}

// Exists returns whether the specified key is present in the global
// configuration. If the global configuration does not exist (no call has
// been made to NewFlexibleConfiguration), an empty configuration is created.
// The configuration store, if set, is checked first. If not found in the
// configuration store or the store was not set, the key is retrieved from
// the memory store created from files, environment variables, and arguments.
func Exists(key string) bool {
	cfg := GetConfiguration()
	return cfg.Exists(key)
}

// Exists returns whether the specified key is present in the configuration.
// The configuration store, if set, is checked first. If not found in the
// configuration store or the store was not set, the key is retrieved from
// the memory store created from files, environment variables, and arguments.
func (fc *flexibleConfiguration) Exists(key string) bool {
	k := strings.TrimSpace(key)
	if len(k) == 0 {
		return false
	}

	if fc.store != nil {
		val, err := fc.store.Get(k)
		if err == nil && len(val) > 0 {
			return true
		}
	}

	val, exists := fc.config[k]
	if !exists {
		return false
	}

	return len(val) > 0
}

// Get returns the value for the specified key from the global configuration.
// If the global configuration does not exist (no call has been made to
// NewFlexibleConfiguration), an empty configuration is created. The
// configuration store, if set, is checked first. If not found in the
// configuration store or the store was not set, the property value is
// retrieved from the memory store created from files, environment variables,
// and arguments.
func Get(key string) string {
	cfg := GetConfiguration()
	return cfg.Get(key)
}

// Get returns the value for the specified key from the configuration.
// The configuration store, if set, is checked first. If not found in the
// configuration store or the store was not set, the key is retrieved from
// the memory store created from files, environment variables, and arguments.
func (fc *flexibleConfiguration) Get(key string) string {
	k := strings.TrimSpace(key)
	if len(k) == 0 {
		return ""
	}

	if fc.store != nil {
		val, err := fc.store.Get(k)
		if err == nil && len(val) > 0 {
			return val
		}
	}

	return fc.config[k]
}

// Set stores the key with value in the global configuration. If the global
// configuration does not exist (no call has been made to
// NewFlexibleConfiguration), an empty configuration is created. If they key
// already exists, its value will be overwritten. If the configuration store
// is set, the key with value will be stored in both the configuration store
// as well as the memory store. Otherwise, it will be stored only in the memory
// store.
func Set(key string, val string) {
	cfg := GetConfiguration()
	cfg.Set(key, val)
}

// Set stores the key with value in the configuration. If they key already
// exists, its value will be overwritten. If the configuration store is set,
// the key with value will be stored in both the configuration store as well
// as the memory store. Otherwise, it will be stored only in the memory
// store.
func (fc *flexibleConfiguration) Set(key, val string) {
	k := strings.TrimSpace(key)
	if !propertyNameIsValid(k) {
		return
	}

	if fc.store != nil {
		fc.store.Set(key, val)
		// even if the store saves the property, save it in memory
	}

	fc.config[key] = val
}

// readConfig uses the configuration parameters to read various aspects of
// the local configuration.
func (fc *flexibleConfiguration) readConfig(
	parameters ConfigurationParameters) map[string]string {
	// Read configuration in reverse priority order (lowest priority
	// first) so that a property from a higher priority source will
	// override a previous definition.

	vars := make(map[string]string)

	// configuration files are the lowest priority
	if len(parameters.ApplicationName) > 0 {
		readConfigFiles(vars,
			parameters.ApplicationName,
			parameters.AcceptedFileSuffixes,
			parameters.IniNamePrefix)
	}

	// environment variables override file property definitions
	if parameters.EnvironmentVariablePrefixes != nil &&
		len(parameters.EnvironmentVariablePrefixes) > 0 {
		readEnvVars(vars, parameters.EnvironmentVariablePrefixes)
	}

	// command line arguments override all other local configuration
	readCommandLineArgs(vars, os.Args)

	return vars
}

// nameIsValid returns whether the specified application name or environment
// variable name is valid to be used as a property name.
func nameIsValid(name string) bool {
	if len(name) == 0 {
		return true
	}

	if !unicode.IsLetter([]rune(name)[0]) {
		return false
	}

	for _, c := range name {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
			return false
		}
	}

	return true
}

// propertyNameIsValid returns whether the specified property name consists
// only of accepted characters.
func propertyNameIsValid(name string) bool {
	if len(name) == 0 {
		return false
	}

	if !unicode.IsLetter([]rune(name)[0]) && name[0:0] != "_" {
		return false
	}

	for _, c := range name {
		if !unicode.IsLetter(c) &&
			!unicode.IsDigit(c) &&
			c != '_' && c != '.' && c != '-' {
			return false
		}
	}

	return true
}
