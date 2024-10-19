package encoder

type Option interface {
	apply(*Encoder)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Encoder)

func (f optionFunc) apply(enc *Encoder) {
	f(enc)
}

type StringFilter = func(string) string

func WithNameFilter(filter StringFilter) Option {
	return optionFunc(func(enc *Encoder) {
		prev := enc.nameFilter
		combined := func(s string) string {
			return filter(prev(s))
		}
		enc.nameFilter = combined
	})
}

func WithNullableMap() Option {
	return optionFunc(func(enc *Encoder) {
		enc.nullableMap = true
	})
}

func WithNullableSlice() Option {
	return optionFunc(func(enc *Encoder) {
		enc.nullableSlice = true
	})
}
