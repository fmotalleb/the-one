package option

import (
	"encoding/json"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Optional[T any] struct {
	Option[T]
}

// func (o Optional[T]) IsSome() bool {
// 	return o.opt != nil && o.opt.IsSome()
// }

// func (o Optional[T]) IsNone() bool {
// 	return o.opt == nil || o.opt.IsNone()
// }

// func (o Optional[T]) Unwrap() *T {
// 	if o.opt == nil {
// 		return new(T)
// 	}
// 	return o.opt.Unwrap()
// }

// func (o Optional[T]) UnwrapOr(def T) *T {
// 	if o.opt == nil {
// 		return &def
// 	}
// 	return o.opt.UnwrapOr(def)
// }

func (o *Optional[T]) Decode(_, _ reflect.Type, val interface{}) error {
	if val == nil {
		o.Option = &None[T]{}
		return nil
	}
	var target T
	if err := mapstructure.Decode(val, &target); err != nil {
		return err
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
