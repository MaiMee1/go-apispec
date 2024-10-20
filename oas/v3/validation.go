package oas

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/MaiMee1/go-apispec/oas/internal/validate"
	"github.com/go-playground/validator/v10"
)

func init() {
	err := validate.RegisterValidation("url_fragment", func(fl validator.FieldLevel) bool {
		v := fl.Field()
		if v.Kind() != reflect.String {
			return true // skip
		}
		if _, err := url.Parse(fmt.Sprintf("https://example.com/%s", v.String())); err != nil {
			return false
		}
		return true
	})
	if err != nil {
		panic(err)
	}
}
