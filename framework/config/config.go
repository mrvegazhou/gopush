package Config

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)


func UnmarshalJSON(data []byte, config interface{}, errorOnUnmatchedKeys bool) error {
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)

	if errorOnUnmatchedKeys {
		decoder.DisallowUnknownFields()
	}

	err := decoder.Decode(config)
	if err != nil && err != io.EOF {
		return err
	}
	return nil

}

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
		return UnmarshalJSON(data, config, errorOnUnmatchedKeys)
	default:
		if err := UnmarshalJSON(data, config, errorOnUnmatchedKeys); err == nil {
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
}

func getConfigurationFiles(files ...string) []string {
	var results []string
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

func processTags(config interface{}) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() == reflect.Ptr {
		configValue = configValue.Elem()
	}
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}
	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
		)

		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank {
			if value := fieldStruct.Tag.Get("default"); value != "" {
				//Addr():获取 v 值的地址，相当于 & 取地址操作。v 值必须可寻址。
				//Interface():将 v 值转换为空接口类型。v 值必须可转换为接口类型。
				if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
					return err
				}
			} else if fieldStruct.Tag.Get("required") == "true" {
				return errors.New(fieldStruct.Name + " is required, but blank")
			}
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := processTags(field.Addr().Interface()); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Slice {
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := processTags(field.Index(i).Addr().Interface()); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func Load(config interface{}, errorOnUnmatchedKeys bool, files ...string) error {
	for _, file := range getConfigurationFiles(files...) {
		if err := processFile(config, errorOnUnmatchedKeys, file); err != nil {
			return err
		}
	}
	return processTags(config)
}
