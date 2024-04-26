package utils

func Min(a, b uintptr) uintptr {
	if a < b {
		return a
	}
	return b
}
