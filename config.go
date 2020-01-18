package flexconfig

/*
Copyright 2018-2020 The flexconfig Authors

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
	defaultAppName                    = ""
	defaultConfigurationSuffix        = ".conf"
	flexconfigCommandlineFileLocation = "flexconfig.configuration.file.location"
	flexConfigEnvFileLocation         = "FLEXCONFIG_CONFIGURATION_FILE_LOCATION"
)

var (
	defaultEnvironmentVariablePrefixes []string
)

var (
	// ErrParmNameNotValid indicates the application name is either
	// missing or uses characters outside those accepted as property names.
	ErrParmNameNotValid = errors.New("Application name not valid")

	// ErrParmEnvPrefixNotValid indicates an environment prefix uses
	// characters outside those accepted as property names.
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
// contain configuration files. The use of ApplicationName implies that a
// default set of directories will be searched for configuration properties.
// If the property 'DirectoryList' is specified, the ApplicationName property
// will not be used. By default, the value of ApplicationName is empty. If
// neither ApplicationName not DirectoryList has a value (this is the default),
// no search for configuration files is done. If ApplicationName
// has a value and DirectoryList is empty, directories are searched in the
// following order, where directories later in the list that define properties
// will override the same property found in a directory listed earlier in the
// list:
//     /usr/local/etc/<name>
//     /opt/etc/<name>
//     /etc/opt/<name>
//     /etc/<name>
//     $HOME/.<name>
//     ./.<name>
// where <name> is the value of ApplicationName. All files having a suffix
// specified in AcceptedFileSuffixes are searched for configuration properties.
//
// DirectoryList specifies the list of directories that will be searched for
// configuration files. The directories are searched in the order given. By
// default, all files with a suffix specified in AcceptedFileSuffixes are read
// to find configuration properties. When DirectoryList has a value, the
// property ApplicationName is not used. When DirectoryList does not have a
// value, a directory list is used by default based on the value of
// ApplicationName.
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
	DirectoryList               []string
	EnvironmentVariablePrefixes []string
	AcceptedFileSuffixes        []string
	IniNamePrefix               string
	ConfigurationStore          FlexConfigStore
}

// configuration is a singleton holding the global static configuration. It
// is set when NewFlexibleConfiguration() is called. It is used when
// GetConfiguration() is called. GetConfiguration() is called by the functions
// Get(), Set(), and Exists(). GetConfiguration() is not called when the methods
// Get(), Set(), and Exists() are called as methods of an interface.
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

// NewFlexibleConfiguration initializes and returns a new Configuration. The
// default Configuration is overwritten with this new Configuration.
func NewFlexibleConfiguration(
	parameters ConfigurationParameters) (Config, error) {
	// Note: Creating a new FlexibleConfiguration will overwrite any
	// existing global configuration.

	if !nameIsValid(parameters.ApplicationName) {
		return nil, ErrParmNameNotValid
	}

	// DirectoryList will be validated when the directories are opened.

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

	// Read the static configuration
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
	readFiles := true

	// Check if environment variable specifies the location of a
	// single cconfiguration file.
	configFile := os.Getenv(flexConfigEnvFileLocation)
	if len(configFile) > 0 {
		readSingleConfigFile(vars, configFile, parameters.IniNamePrefix)
		if len(vars) > 0 {
			readFiles = false
		}
	}

	// Check if a command line argument is used to specify the location
	// of a single configuration file.
	configFile = searchArgument(os.Args, flexconfigCommandlineFileLocation)
	if len(configFile) > 0 {
		if len(vars) > 0 {
			vars = make(map[string]string)
		}

		readSingleConfigFile(vars, configFile, parameters.IniNamePrefix)
		if len(vars) > 0 {
			readFiles = false
		}
	}

	// configuration files are the lowest priority
	if readFiles {
		if len(parameters.DirectoryList) > 0 {
			readConfigFilesFromDirectoryList(vars,
				parameters.DirectoryList,
				parameters.AcceptedFileSuffixes,
				parameters.IniNamePrefix)
		} else if len(parameters.ApplicationName) > 0 {
			readConfigFiles(vars,
				parameters.ApplicationName,
				parameters.AcceptedFileSuffixes,
				parameters.IniNamePrefix)
		}
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

// readSingleConfigFile reads properties set in a single configuration file.
func readSingleConfigFile(
	vars map[string]string,
	configFile, iniNamePrefix string) {
	// Break file name into path and name and read the file at
	// that location.

	path := ""
	name := ""
	index := strings.LastIndex(configFile, "/")
	if index < 0 {
		path = "."
		name = configFile
	} else {
		path = configFile[:index]
		name = configFile[index+1:]
	}

	readConfigFile(vars, path, name, iniNamePrefix)
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
