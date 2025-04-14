package option

import (
	"reflect"
	"regexp"

	"github.com/mitchellh/mapstructure"
)

type Decodable interface {
	Decode(reflect.Type, reflect.Type, interface{}) error
}

func DecodeHookFunc() mapstructure.DecodeHookFunc {
	optionalRegex := regexp.MustCompile("Optional[?.*]?")
	optionalTRegex := regexp.MustCompile("OptionalT[?.*]?")
	someRegex := regexp.MustCompile("Some[?.*]?")
	return func(from, to reflect.Type, val interface{}) (interface{}, error) {
		// Optional
		if optionalRegex.Match([]byte(to.Name())) {
			opt := reflect.New(to).Interface().(Decodable)
			if err := opt.Decode(from, to, val); err != nil {
				return nil, err
			}
			return reflect.ValueOf(opt).Elem().Interface(), nil
		}

		if optionalTRegex.Match([]byte(to.Name())) {
			opt := reflect.New(to).Interface().(Decodable)
			if err := opt.Decode(from, to, val); err != nil {
				return nil, err
			}
			return reflect.ValueOf(opt).Elem().Interface(), nil
		}

		// Some
		if someRegex.Match([]byte(to.Name())) {
			some := reflect.New(to).Interface().(Decodable)
			if err := some.Decode(from, to, val); err != nil {
				return nil, err
			}
			return reflect.ValueOf(some).Elem().Interface(), nil
		}

		return val, nil
	}
}
