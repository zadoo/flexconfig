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
	if len(name) == 0 {
		return
	}

	dirlist := []string{
		"/usr/local/etc/" + name,
		"/opt/etc/" + name,
		"/opt/" + name + "/etc",
		"/etc/opt/" + name,
		"/etc/" + name,
	}

	homedir := os.Getenv("HOME")
	if len(homedir) > 0 {
		dir := homedir
		if !strings.HasSuffix(dir, "/") {
			dir = dir + "/"
		}

		dirlist = append(dirlist, dir+"."+name)
	}

	dirlist = append(dirlist, "."+name)

	readConfigFilesFromDirectoryList(vars, dirlist, suffixes, iniPrefix)
}

// readConfigFilesFromDirectory reads configuration files from a list of
// directories. Config files are read from directories in the order listed.
func readConfigFilesFromDirectoryList(vars map[string]string, dirlist []string, suffixes []string, iniPrefix string) {
	// TODO: consider how this function can be tested, perhaps by passing
	//       a function to it to perform the work where for testing a
	//       test function could be passed in.
	for _, dir := range dirlist {
		readFiles(vars, dir, suffixes, iniPrefix)
	}
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
