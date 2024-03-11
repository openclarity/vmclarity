package scanner

func Ptr[T any](input T) *T {
	return &input
}

func ValOrEmpty[T any](input *T) T {
	if input == nil {
		var ret T
		return ret
	}
	return *input
}
