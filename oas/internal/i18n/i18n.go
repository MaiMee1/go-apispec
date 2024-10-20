package i18n

import (
	"io"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

var (
	uni *ut.UniversalTranslator
)

func FindTranslator(locales ...string) (trans ut.Translator, found bool) {
	return uni.FindTranslator(locales...)
}
func GetTranslator(locale string) (trans ut.Translator, found bool) {
	return uni.GetTranslator(locale)
}
func GetFallback() ut.Translator {
	return uni.GetFallback()
}
func VerifyTranslations() (err error) {
	return uni.VerifyTranslations()
}

func init() {
	en := en.New()
	uni = ut.New(en, en)
}

type UniversalTranslator interface {
	Export(format ut.ImportExportFormat, dirname string) error
	Import(format ut.ImportExportFormat, dirnameOrFilename string) error
	ImportByReader(format ut.ImportExportFormat, reader io.Reader) error
	FindTranslator(locales ...string) (trans ut.Translator, found bool)
	GetTranslator(locale string) (trans ut.Translator, found bool)
	GetFallback() ut.Translator
	AddTranslator(translator locales.Translator, override bool) error
	VerifyTranslations() (err error)
}
