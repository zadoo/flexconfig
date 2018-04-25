package flexconfig

import (
	"strings"

	"gopkg.in/ini.v1"
)

// parseIniFile parses the specified file contents expecting it to use the
// format of an INI file.
func parseIniFile(vars map[string]string, prefix, content string) error {
	cfg, err := ini.Load([]byte(content))
	if err != nil {
		return err
	}

	if len(prefix) > 0 && prefix[len(prefix)-1:] != "." {
		prefix = prefix + "."
	}

	return parseIni(vars, prefix, cfg)
}

// parseIni accepts a loaded ini structure and iterates through the sections
// creating configruation properties from the definitions in each section.
func parseIni(vars map[string]string, prefix string, cfg *ini.File) error {
	for _, sect := range cfg.SectionStrings() {
		if sect == "DEFAULT" {
			continue
		}

		keys := cfg.Section(sect).Keys()

		for _, k := range keys {
			key := iniToConfigKey(prefix, sect, k.Name())
			vars[key] = k.Value()
		}
	}

	return nil
}

// iniToConfigKey returns the configuration property name equivalent for
// an INI section and property name, using the specified prefix as the
// property name prefix.
func iniToConfigKey(prefix, section, key string) string {
	return prefix + strings.ToLower(section) + "." + strings.ToLower(key)
}
