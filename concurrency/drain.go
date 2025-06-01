package concurrency

func Drain[T any](ch <-chan T) []T {
	result := make([]T, 0)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return result
			}
			result = append(result, v)
		default:
			return result
		}
	}
}
