package security

import "github.com/MaiMee1/go-apispec/oas/v3"

func None() oas.SecurityRequirement {
	return oas.SecurityRequirement{}
}

type All = oas.SecurityRequirement
