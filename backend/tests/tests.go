package tests

func ToPtr[V any](v V) *V {
	return &v
}
