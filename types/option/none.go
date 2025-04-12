package option

type None[T any] struct{}

func (n *None[T]) IsSome() bool { return false }
func (n *None[T]) IsNone() bool { return true }
func (n *None[T]) Unwrap() *T {
	return new(T)
}
func (n *None[T]) UnwrapOr(def T) *T { return &def }
