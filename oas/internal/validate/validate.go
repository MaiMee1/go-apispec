package validate

import (
	"github.com/go-playground/validator/v10"
)

var Inst = validator.New(validator.WithRequiredStructEnabled())

func init() {
	if err := Inst.RegisterValidation("uri-reference", func(fl validator.FieldLevel) bool {
		fl.Field()
		// TODO: use subprocess to use node to validate regex
		return true
	}); err != nil {
		panic(err)
	}
	if err := Inst.RegisterValidation("regex", func(fl validator.FieldLevel) bool {
		fl.Field()
		// TODO: use subprocess to use node to validate regex
		return true
	}); err != nil {
		panic(err)
	}
}
