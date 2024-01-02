package arrays

func Contains(array []string, searched string) bool {
	if array == nil || len(searched) == 0 {
		return false
	}

	for _, item := range array {
		if item == searched {
			return true
		}
	}
	return false
}
