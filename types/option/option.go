package option

type Option[T any] interface {
	IsSome() bool
	IsNone() bool
	Unwrap() *T
	UnwrapOr(def T) *T
}

func New[T any](data *T) Option[T] {
	if data == nil {
		return &None[T]{}
	}
	return &Some[T]{data: data}
}

func NewOptional[T any](data *T) Optional[T] {
	opt := New(data)
	return Optional[T]{
		opt: opt,
	}
}
