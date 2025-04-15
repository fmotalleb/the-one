package option

func UnwrapAll[T any, O Option[T]](items []O) []T {
	result := make([]T, len(items))
	for index, item := range items {
		result[index] = *item.Unwrap()
	}
	return result
}
