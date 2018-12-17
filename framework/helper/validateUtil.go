package helper

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"sync"
)

// validate v8 升级 v9

type ValidatorV9 struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &ValidatorV9{}

func (v *ValidatorV9) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *ValidatorV9) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *ValidatorV9) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
