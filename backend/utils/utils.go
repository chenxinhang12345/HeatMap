package utils

func Min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
