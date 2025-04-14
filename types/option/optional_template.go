package option

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/fmotalleb/the-one/template"
)

type OptionalT[T any] struct {
	opt Option[T]
}

func (o OptionalT[T]) IsSome() bool {
	return o.opt != nil && o.opt.IsSome()
}

func (o OptionalT[T]) IsNone() bool {
	return o.opt == nil || o.opt.IsNone()
}

func (o OptionalT[T]) Unwrap() *T {
	if o.opt == nil {
		return new(T)
	}
	return o.opt.Unwrap()
}

func (o OptionalT[T]) UnwrapOr(def T) *T {
	if o.opt == nil {
		return &def
	}
	return o.opt.UnwrapOr(def)
}

func (o *OptionalT[T]) Decode(_, to reflect.Type, val interface{}) error {
	if val == nil {
		o.opt = &None[T]{}
		return nil
	}
	parsed, err := mapper[T](to, val)
	if err != nil {
		return err
	}
	var target T
	if err := mapstructure.Decode(parsed, &target); err != nil {
		return err
	}

	o.opt = New(&target)
	return nil
}

func mapper[T any](to reflect.Type, val interface{}) (T, error) {
	var zero T

	strVal := fmt.Sprintf("%v", val)
	strVal, err := template.EvaluateTemplate(strVal, map[string]any{})
	if err != nil {
		return zero, err
	}

	// If T implements TextUnmarshaler
	if to.Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
		v := reflect.New(to).Interface().(encoding.TextUnmarshaler)
		if err := v.UnmarshalText([]byte(strVal)); err != nil {
			return zero, err
		}
		return reflect.ValueOf(v).Elem().Interface().(T), nil
	}

	// Try decoding the string using JSON (handles bool/int/float/etc.)
	var target T
	if err := json.Unmarshal([]byte(strVal), &target); err == nil {
		return target, nil
	}

	// Fallback: try decoding wrapped string
	if err := mapstructure.Decode(strVal, &target); err != nil {
		return zero, err
	}
	return target, nil
}

func (o OptionalT[T]) Encode() (interface{}, error) {
	if o.IsNone() {
		return nil, nil
	}
	return o.Unwrap(), nil
}

func (o OptionalT[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(o.Unwrap())
}

func (o OptionalT[T]) MarshalYAML() (interface{}, error) {
	if o.IsNone() {
		return nil, nil
	}
	return o.Unwrap(), nil
}
