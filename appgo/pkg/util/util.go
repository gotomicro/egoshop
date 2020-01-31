package util

// InIntArray 判断是否在int数组中
func InIntArray(v int, sl []int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
