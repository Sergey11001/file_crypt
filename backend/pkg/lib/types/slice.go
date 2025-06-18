package types

func SliceDiff[T comparable](a, b []T) []T {
	var diff []T
	m := make(map[T]bool, len(b))

	for _, v := range b {
		m[v] = true
	}

	for _, v := range a {
		if !m[v] {
			diff = append(diff, v)
		}
	}

	return diff
}
