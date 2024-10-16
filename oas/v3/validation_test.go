package oas

import (
	"testing"
)

func TestReference_Validate(t *testing.T) {
	var r Reference
	r.Ref = "#/components/schemas/repo.blockfint.com.thinker_core.generic_authentication.src.app.GetRolesByPolicyParams"
	err := validate.Struct(r)
	if err != nil {
		t.Error(err)
	}
}
