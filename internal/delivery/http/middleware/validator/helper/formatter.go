package validatorhelper

import (
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// FieldJSONFormatter function to create field name same with json name
func FieldJSONFormatter(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
