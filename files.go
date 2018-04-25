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
	"sort"
	"strings"
)

// readConfigFiles performs a search for config files in an ordered set of
// standard directories that may contain configuration. The specified name is
// the last field of the name of a directory in one of the standard locations.
// Configuration files found at that directory are read, creating configuration
// properties.
func readConfigFiles(vars map[string]string, name string, suffixes []string, iniPrefix string) {
	readFiles(vars, "/usr/local/etc/"+name, suffixes, iniPrefix)
	readFiles(vars, "/opt/etc/"+name, suffixes, iniPrefix)
	readFiles(vars, "/etc/opt/"+name, suffixes, iniPrefix)
	readFiles(vars, "/etc/"+name, suffixes, iniPrefix)

	homedir := os.Getenv("HOME")
	if len(homedir) > 0 {
		dir := homedir
		if !strings.HasSuffix(dir, "/") {
			dir = dir + "/"
		}

		readFiles(vars, dir+"."+name, suffixes, iniPrefix)
	}

	// If the current working directory is the same as $HOME, this will
	// read a set of config files a second time. There should be no change
	// in the resulting configuration.
	readFiles(vars, "."+name, suffixes, iniPrefix)
}

// readFiles checks for and reads configuration files in a single directory.
// If the directory exists, files with any of the specified suffixes are
// read and configuration properties created.
func readFiles(vars map[string]string, dirname string, suffixes []string, iniPrefix string) {
	dir, err := os.Open(dirname)
	if err != nil {
		return
	}

	defer dir.Close()

	filenames, err := dir.Readdirnames(0)
	if err != nil {
		return
	}

	if filenames == nil || len(filenames) == 0 {
		return
	}

	sort.Strings(filenames)

	for _, f := range filenames {
		for _, suffix := range suffixes {
			if strings.HasSuffix(f, suffix) {
				readConfigFile(vars, dirname, f, iniPrefix)
			}
		}
	}
}

// readConfigFile reads a single configuration file and creates configuration
// properties based on its contents. If file contents are json, yaml, or ini,
// properties are created. Other file types are ignored.
func readConfigFile(vars map[string]string, path string, name string, iniPrefix string) {
	fileContents, err := ioutil.ReadFile(path + "/" + name)
	if err != nil {
		return
	}

	contents := string(fileContents)

	// Parse either yaml or json
	err = parseYaml(vars, contents)
	if err != nil {
		// File contents were neither YAML nor JSON, try INI
		err = parseIniFile(vars, iniPrefix, contents)
		if err != nil {
			// Unknown file type, ignore the file
		}
	}
}
