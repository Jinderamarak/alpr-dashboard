package util

func MapPtr[U, V any](ptr *U, mapper func(U) V) *V {
	if ptr == nil {
		return nil
	}
	mapped := mapper(*ptr)
	return &mapped
}
