package oas

import (
	"testing"

	"github.com/MaiMee1/go-apispec/oas/internal/validate"
)

func TestReference_Validate(t *testing.T) {
	var r Reference
	r.Ref = "#/components/schemas/Pet"
	err := validate.Inst.Struct(r)
	if err != nil {
		t.Error(err)
	}
}
