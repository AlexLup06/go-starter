package ptr

func Ptr[T any](val T) *T {
	return &val
}

func DefaultIfNil[T any](val *T, defaultVal T) T {
	if val == nil {
		return defaultVal
	}
	return *val
}
