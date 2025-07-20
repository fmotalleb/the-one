package option

import (
	"encoding/json"
	"reflect"

	"github.com/fmotalleb/the-one/types/decodable"
)

type Optional[T any] struct {
	Option[T]
}

func (o *Optional[T]) Decode(_ reflect.Type, val interface{}) (any, error) {
	if val == nil {
		o.Option = &None[T]{}
		return val, nil
	}
	var target T
	result, err := decodable.UsingParserOf[T](val)
	if result != nil {
		if err != nil {
			return nil, err
		}
		target = *result
	} else {
		err := decodable.Decode(val, &target)
		if err != nil {
			return nil, err
		}
	}

	o.Option = New(&target)
	return o, nil
}

func (o *Optional[T]) Encode() (interface{}, error) {
	if o.IsNone() {
		return nil, nil
	}
	return o.Unwrap(), nil
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(o.Unwrap())
}

func (o Optional[T]) MarshalYAML() (interface{}, error) {
	if o.IsNone() {
		return nil, nil
	}
	return o.Unwrap(), nil
}

func (o Optional[T]) IsSome() bool {
	if o.Option == nil {
		return false
	}
	return o.Option.IsSome()
}

func (o Optional[T]) IsNone() bool {
	if o.Option == nil {
		return true
	}
	return o.Option.IsNone()
}

func (o Optional[T]) Unwrap() *T {
	if o.Option == nil {
		panic("called Unwrap on a None value")
	}
	return o.Option.Unwrap()
}

func (o Optional[T]) UnwrapOr(def T) T {
	if o.Option == nil {
		return def
	}
	return o.Option.UnwrapOr(def)
}
