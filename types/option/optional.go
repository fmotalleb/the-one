package option

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type Optional[T any] struct {
	opt Option[T]
}

func (o Optional[T]) IsSome() bool      { return o.opt.IsSome() }
func (o Optional[T]) IsNone() bool      { return o.opt.IsNone() }
func (o Optional[T]) Unwrap() *T        { return o.opt.Unwrap() }
func (o Optional[T]) UnwrapOr(def T) *T { return o.opt.UnwrapOr(def) }

func (o *Optional[T]) Decode(val interface{}) error {
	if val == nil {
		o.opt = &None[T]{}
		return nil
	}

	var target T
	if err := mapstructure.Decode(val, &target); err != nil {
		return err
	}

	o.opt = New(&target)
	return nil
}

func (o Optional[T]) Encode() (interface{}, error) {
	if o.opt == nil || o.opt.IsNone() {
		return nil, nil
	}
	return o.opt.Unwrap(), nil
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.opt == nil || o.opt.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(o.opt.Unwrap())
}

func (o Optional[T]) MarshalYAML() (interface{}, error) {
	if o.opt == nil || o.opt.IsNone() {
		return nil, nil // Represents YAML null
	}
	return o.opt.Unwrap(), nil
}
