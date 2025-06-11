package utils

func IntPtrAddOne(p **int) {
	if *p == nil {
		val := 1
		*p = &val
	} else {
		**p++
	}
}

func Ptr[T any](v T) *T {
	return &v
}
