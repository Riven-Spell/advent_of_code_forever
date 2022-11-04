package util

func Ternary[T any](opt bool, tVal, fVal T) T {
	if opt {
		return tVal
	}

	return fVal
}
