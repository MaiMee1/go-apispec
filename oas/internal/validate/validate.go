package validate

import (
	"context"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func Struct(s interface{}) error {
	return validate.Struct(s)
}
func StructCtx(ctx context.Context, s interface{}) (err error) {
	return validate.StructCtx(ctx, s)
}
func StructFiltered(s interface{}, fn validator.FilterFunc) error {
	return validate.StructFiltered(s, fn)
}
func StructFilteredCtx(ctx context.Context, s interface{}, fn validator.FilterFunc) (err error) {
	return validate.StructFilteredCtx(ctx, s, fn)
}
func StructPartial(s interface{}, fields ...string) error {
	return validate.StructPartial(s, fields...)
}
func StructPartialCtx(ctx context.Context, s interface{}, fields ...string) (err error) {
	return validate.StructPartialCtx(ctx, s, fields...)
}
func StructExcept(s interface{}, fields ...string) error {
	return validate.StructExcept(s, fields...)
}
func StructExceptCtx(ctx context.Context, s interface{}, fields ...string) (err error) {
	return validate.StructExceptCtx(ctx, s, fields...)
}
func Var(field interface{}, tag string) error {
	return validate.Var(field, tag)
}
func VarCtx(ctx context.Context, field interface{}, tag string) (err error) {
	return validate.VarCtx(ctx, field, tag)
}
func VarWithValue(field interface{}, other interface{}, tag string) error {
	return validate.VarWithValue(field, other, tag)
}
func VarWithValueCtx(ctx context.Context, field interface{}, other interface{}, tag string) (err error) {
	return validate.VarWithValueCtx(ctx, field, other, tag)
}

type Validate interface {
	SetTagName(name string)
	ValidateMapCtx(ctx context.Context, data map[string]interface{}, rules map[string]interface{}) map[string]interface{}
	ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{}
	RegisterTagNameFunc(fn validator.TagNameFunc)
	RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error
	RegisterValidationCtx(tag string, fn validator.FuncCtx, callValidationEvenIfNull ...bool) error
	RegisterAlias(alias, tags string)
	RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{})
	RegisterStructValidationCtx(fn validator.StructLevelFuncCtx, types ...interface{})
	RegisterStructValidationMapRules(rules map[string]string, types ...interface{})
	RegisterCustomTypeFunc(fn validator.CustomTypeFunc, types ...interface{})
	RegisterTranslation(tag string, trans ut.Translator, registerFn validator.RegisterTranslationsFunc, translationFn validator.TranslationFunc) (err error)
	Struct(s interface{}) error
	StructCtx(ctx context.Context, s interface{}) (err error)
	StructFiltered(s interface{}, fn validator.FilterFunc) error
	StructFilteredCtx(ctx context.Context, s interface{}, fn validator.FilterFunc) (err error)
	StructPartial(s interface{}, fields ...string) error
	StructPartialCtx(ctx context.Context, s interface{}, fields ...string) (err error)
	StructExcept(s interface{}, fields ...string) error
	StructExceptCtx(ctx context.Context, s interface{}, fields ...string) (err error)
	Var(field interface{}, tag string) error
	VarCtx(ctx context.Context, field interface{}, tag string) (err error)
	VarWithValue(field interface{}, other interface{}, tag string) error
	VarWithValueCtx(ctx context.Context, field interface{}, other interface{}, tag string) (err error)
}
