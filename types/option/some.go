package option

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Some[T any] struct {
	data *T
}

func (s Some[T]) IsSome() bool   { return true }
func (s Some[T]) IsNone() bool   { return false }
func (s Some[T]) Unwrap() *T     { return s.data }
func (s Some[T]) UnwrapOr(_ T) T { return *s.data }

func (s *Some[T]) Decode(_ reflect.Type, val interface{}) (any, error) {
	if val == nil {
		return nil, errors.New("required input is missing")
	}

	var target T
	if err := mapstructure.Decode(val, &target); err != nil {
		return nil, err
	}

	s.data = &target
	return s, nil
}

func (s Some[T]) Encode() (interface{}, error) {
	return any(s.data), nil
}

func (s Some[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

func (s Some[T]) MarshalYAML() (interface{}, error) {
	return s.Unwrap(), nil
}
