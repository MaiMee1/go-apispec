package validate

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func init() {
	validate = validator.New()

	if err := validate.RegisterValidation("url_fragment", func(fl validator.FieldLevel) bool {
		v := fl.Field()
		if v.Kind() != reflect.String {
			return true // skip
		}
		if _, err := url.Parse(fmt.Sprintf("https://example.com/%s", v.String())); err != nil {
			return false
		}
		return true
	}); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("uri-reference", func(fl validator.FieldLevel) bool {
		fl.Field()
		return true
	}); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("regex", func(fl validator.FieldLevel) bool {
		fl.Field()
		// TODO: use subprocess to use node to validate regex
		return true
	}); err != nil {
		panic(err)
	}
}
