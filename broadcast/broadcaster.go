package broadcast

type Subscription[T any] interface {
	Subscribe() <-chan T
}

type Broadcaster[T any] struct {
	subscribers []chan T
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		subscribers: make([]chan T, 0),
	}
}

func (b *Broadcaster[T]) Subscribe() <-chan T {
	// This method has an overhead where it does not release nor reuse its subscribers
	ch := make(chan T, 1)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Broadcaster[T]) Broadcast(input <-chan T) {
	for val := range input {
		for _, sub := range b.subscribers {
			sub <- val
		}
	}
}

func (b *Broadcaster[T]) BindTo(input <-chan T) {
	for val := range input {
		for _, sub := range b.subscribers {
			sub <- val
		}
	}
	for _, sub := range b.subscribers {
		close(sub)
	}
	b.subscribers = make([]chan T, 0)
}
