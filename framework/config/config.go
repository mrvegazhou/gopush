package config

import (
	"fmt"
	"errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"json/jsonutil"
)



func processFile(config interface{}, errorOnUnmatchedKeys bool, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		if errorOnUnmatchedKeys {
			return yaml.UnmarshalStrict(data, config)
		}
		return yaml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".json"):
		return jsonutil.unmarshalJSON(data, config, errorOnUnmatchedKeys)
	}
	default:
		if err := jsonutil.unmarshalJSON(data, config, errorOnUnmatchedKeys); err == nil {
			return nil
		} else if strings.Contains(err.Error(), "json: unknown field") {
			return err
		}

		var yamlError error
		if errorOnUnmatchedKeys {
			yamlError = yaml.UnmarshalStrict(data, config)
		} else {
			yamlError = yaml.Unmarshal(data, config)
		}

		if yamlError == nil {
			return nil
		} else if yErr, ok := yamlError.(*yaml.TypeError); ok {
			return yErr
		}

		return errors.New("failed to decode config")
}

func getConfigurationFiles(files ...string) []string {
	for i := len(files) - 1; i >= 0; i-- {
		foundFile := false
		file := files[i]
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			foundFile = true
			results = append(results, file)
		}
		if !foundFile {
			fmt.Printf("Failed to find configuration %v\n", file)
		}
	}
	return results
}

func Load(config interface{}, errorOnUnmatchedKeys, files ...string) error {
	for _, file := range getConfigurationFiles(files...) {
		if err := processFile(config, errorOnUnmatchedKeys, file); err != nil {
			return err
		}
	}

}