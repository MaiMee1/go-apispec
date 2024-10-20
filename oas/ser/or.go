// Package ser provides serialization helpers
package ser

import (
	"encoding/json"
	"fmt"
)

// Or allows two types to be serializable or deserialized to one value.
//
// Type A is tried first before type B for both serialization and deserialization.
type Or[A, B comparable] struct {
	X A
	Y B
}

//goland:noinspection GoMixedReceiverTypes
func (o *Or[A, B]) UnmarshalJSON(b []byte) error {
	var (
		x A
		y B
	)
	err1 := json.Unmarshal(b, &x)
	if err1 == nil {
		o.X = x
	}

	err2 := json.Unmarshal(b, &y)
	if err2 == nil {
		o.Y = y
	}

	if err2 == nil || err1 == nil {
		return nil
	}
	return fmt.Errorf("or[%T, %T]: %w: %w", x, y, err2, err1)
}

//goland:noinspection GoMixedReceiverTypes
func (o Or[A, B]) MarshalJSON() ([]byte, error) {
	var x A
	if o.X != x {
		return json.Marshal(o.X)
	}
	return json.Marshal(o.Y)
}
