package option

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/fmotalleb/the-one/template"
	"github.com/fmotalleb/the-one/types/decodable"
)

type OptionalT[T any] struct {
	Optional[T]
}

func (o *OptionalT[T]) Decode(_, to reflect.Type, template interface{}) error {
	if template == nil {
		o.Option = &None[T]{}
		return nil
	}

	raw, err := applyTemplate(template)
	if err != nil {
		return err
	}
	parsed, err := transform[T](to, raw)
	if err != nil {
		return err
	}
	var target T

	if err := decodable.Decode(parsed, &target); err != nil {
		return err
	}

	o.Option = New(&target)
	return nil
}

func applyTemplate(val interface{}) (string, error) {
	strVal := fmt.Sprintf("%v", val)
	strVal, err := template.EvaluateTemplate(strVal, map[string]any{})
	if err != nil {
		return "", err
	}

	// If T implements TextUnmarshaler
	return strVal, nil
}

func transform[T any](to reflect.Type, strVal string) (T, error) {
	if to.Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
		v := reflect.New(to).Interface().(encoding.TextUnmarshaler)
		if err := v.UnmarshalText([]byte(strVal)); err != nil {
			var zero T
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
		var zero T
		return zero, err
	}
	return target, nil
}

func (o OptionalT[T]) IsSome() bool {
	if o.Option == nil {
		return false
	}
	return o.Option.IsSome()
}

func (o OptionalT[T]) IsNone() bool {
	if o.Option == nil {
		return true
	}
	return o.Option.IsNone()
}

func (o OptionalT[T]) Unwrap() *T {
	if o.Option == nil {
		panic("called Unwrap on a None value")
	}
	return o.Option.Unwrap()
}

func (o OptionalT[T]) UnwrapOr(def T) T {
	if o.Option == nil {
		return def
	}
	return o.Option.UnwrapOr(def)
}
