package option

import (
	"encoding/json"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/fmotalleb/the-one/types/decodable"
)

type Optional[T any] struct {
	Option[T]
}

func (o *Optional[T]) Decode(from, to reflect.Type, val interface{}) error {
	if val == nil {
		o.Option = &None[T]{}
		return nil
	}
	var target T
	result, err := decodable.UsingParserOf[T](val)
	if result != nil {
		if err != nil {
			return err
		}
		target = *result
	} else {
		err := mapstructure.Decode(val, &target)
		if err != nil {
			return err
		}
	}

	o.Option = New(&target)
	return nil
}

func (o *Optional[T]) Encode() (interface{}, error) {
	if o.IsNone() {
		return nil, nil
	}
	return o.Unwrap(), nil
}

func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(o.Unwrap())
}

func (o *Optional[T]) MarshalYAML() (interface{}, error) {
	if o.IsNone() {
		return nil, nil
	}
	return o.Unwrap(), nil
}
