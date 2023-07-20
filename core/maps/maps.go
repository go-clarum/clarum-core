package maps

import "github.com/goclarum/clarum/core/arrays"

func EqualString(m1 map[string]string, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}

	for m1Key, m1Value := range m1 {
		if m2Value, exists := m2[m1Key]; exists {
			if m1Value != m2Value {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func EqualStrings(m1 map[string][]string, m2 map[string][]string) bool {
	if len(m1) != len(m2) {
		return false
	}

	for m1Key, m1Values := range m1 {
		if m2Values, exists := m2[m1Key]; exists {
			if !arrays.Equal(m1Values, m2Values) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}
