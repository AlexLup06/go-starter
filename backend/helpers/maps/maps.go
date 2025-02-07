package maps

func Values[K comparable, S any](m map[K]S) []S {
	result := make([]S, 0)
	if m == nil {
		return result
	}
	for key := range m {
		result = append(result, m[key])
	}
	return result
}
