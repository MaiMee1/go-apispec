// Package flag defines functions useful with flag enums of any underlying type
package flag

import (
	"iter"

	"golang.org/x/exp/constraints"
)

func Has[F constraints.Unsigned](flag F, possible iter.Seq[F], ands ...F) bool {
	for _, and := range ands {
		var ok = false
		for or := range Range(and, possible) {
			if flag&or == or {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func Range[F constraints.Unsigned](flag F, possible iter.Seq[F]) iter.Seq[F] {
	return func(yield func(F) bool) {
		for typ := range possible {
			if flag&typ == typ {
				if !yield(typ) {
					return
				}
			}
		}
	}
}
