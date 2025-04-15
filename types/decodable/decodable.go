package decodable

import (
	"reflect"
)

type Decodable interface {
	Decode(reflect.Type, reflect.Type, interface{}) error
}
