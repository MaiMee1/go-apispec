package oas

import (
	"testing"
)

func TestReference_Validate(t *testing.T) {
	var r Reference
	r.Ref = "#/components/schemas/Pet"
	err := validate.Struct(r)
	if err != nil {
		t.Error(err)
	}
}
