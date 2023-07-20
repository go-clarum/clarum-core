package arrays

func Equal(a1 []string, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}

out:
	for _, v1 := range a1 {
		for _, v2 := range a2 {
			if v1 == v2 {
				continue out
			}
		}
		return false
	}

	return true
}

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
