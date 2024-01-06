package ptr

func ToPtr[T any](value T) *T {
	return &value
}
