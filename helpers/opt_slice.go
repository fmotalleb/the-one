package helpers

import (
	"github.com/fmotalleb/the-one/types/option"
)

func OptToSlice[T any](slice []option.OptionalT[T], def ...T) []T {
	out := make([]T, len(slice))
	var defValue T
	if len(def) != 0 {
		defValue = def[0]
	}
	for i, v := range slice {
		out[i] = v.UnwrapOr(defValue)
	}
	return out
}

func SomeToSlice[T any](slice []option.Some[T]) []T {
	out := make([]T, len(slice))
	for i, v := range slice {
		out[i] = *v.Unwrap()
	}
	return out
}
