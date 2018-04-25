package flexconfig

import (
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

// parseYaml parses the specified content expecting json or yaml format,
// creating configuration properties based on the content.
//
func parseYaml(vars map[string]string, contents string) error {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(contents), &m)
	if err != nil {
		return err
	}

	setYamlStruct(vars, "", m)

	return nil
}

// setYamlStruct accepts a parsed yaml structure (map of interfaces) and
// creates configuration properties representing the content.
func setYamlStruct(
	vars map[string]string,
	prefix string,
	m map[interface{}]interface{}) {
	for k, v := range m {
		key := prefix + k.(string)
		setYamlVar(vars, key, v)
	}
}

// setYamlVar accepts a parsed yaml key and value and creates configuration
// property (or properties if it is a struct or array) representing the value.
func setYamlVar(vars map[string]string, key string, v interface{}) {
	switch reflect.TypeOf(v).Name() {
	case "string":
		vars[key] = v.(string)
	case "int":
		vars[key] = strconv.FormatInt(int64(v.(int)), 10)
	case "bool":
		vars[key] = strconv.FormatBool(v.(bool))
	case "float64":
		vars[key] = strconv.FormatFloat(v.(float64), 'g', -1, 64)
	default:
		if reflect.TypeOf(v).Kind() == reflect.Array ||
			reflect.TypeOf(v).Kind() == reflect.Slice {
			a := v.([]interface{})
			for i, av := range a {
				setYamlVar(vars, key+"."+strconv.Itoa(i), av)
			}

			return
		}

		setYamlStruct(vars, key+".", v.(map[interface{}]interface{}))
	}
}
