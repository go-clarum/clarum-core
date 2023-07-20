package maps

import "testing"

func TestEqualStringNilAndEmpty(t *testing.T) {
	var m1 map[string]string = nil
	var m2 = make(map[string]string)

	if !EqualString(m1, m2) {
		t.Errorf("Expected <true>")
	}
}

func TestEqualString(t *testing.T) {
	m1 := map[string]string{"k1": "v1", "k2": "v2"}
	m2 := map[string]string{"k2": "v2", "k1": "v1"}

	if !EqualString(m1, m2) {
		t.Errorf("Expected <true>")
	}
}

func TestNotEqualStringDiffValue(t *testing.T) {
	m1 := map[string]string{"k1": "v1", "k2": "v2"}
	m2 := map[string]string{"k2": "v2", "k1": "v3"}

	if EqualString(m1, m2) {
		t.Errorf("Expected <false>")
	}
}

func TestNotEqualStringDiffKey(t *testing.T) {
	m1 := map[string]string{"k1": "v1", "k2": "v2"}
	m2 := map[string]string{"k2": "v2", "k3": "v1"}

	if EqualString(m1, m2) {
		t.Errorf("Expected <false>")
	}
}

func TestNotEqualStringDiffCount(t *testing.T) {
	m1 := map[string]string{"k1": "v1", "k2": "v2"}
	m2 := map[string]string{"k2": "v2", "k3": "v1", "k4": "v3"}

	if EqualString(m1, m2) {
		t.Errorf("Expected <false>")
	}
}

func TestEqualStringsNilAndEmpty(t *testing.T) {
	var m1 map[string][]string = nil
	var m2 = make(map[string][]string)

	if !EqualStrings(m1, m2) {
		t.Errorf("Expected <true>")
	}
}

func TestEqualStrings(t *testing.T) {
	m1 := map[string][]string{"k1": {"v3", "v4", "v2"}, "k2": {}}
	m2 := map[string][]string{"k2": {}, "k1": {"v2", "v3", "v4"}}

	if !EqualStrings(m1, m2) {
		t.Errorf("Expected <true>")
	}
}

func TestNotEqualStringsDiffValue(t *testing.T) {
	m1 := map[string][]string{"k1": {"v3", "v4", "v1"}, "k2": {}}
	m2 := map[string][]string{"k2": {}, "k1": {"v3", "v4", "v2"}}

	if EqualStrings(m1, m2) {
		t.Errorf("Expected <false>")
	}
}

func TestNotEqualStringsDiffKey(t *testing.T) {
	m1 := map[string][]string{"k1": {"v3", "v4", "v1"}, "k2": {}}
	m2 := map[string][]string{"k2": {}, "k3": {"v3", "v4", "v2"}}

	if EqualStrings(m1, m2) {
		t.Errorf("Expected <false>")
	}
}

func TestNotEqualStringsDiffCountMap(t *testing.T) {
	m1 := map[string][]string{"k1": {"v3", "v4", "v1"}, "k2": {}}
	m2 := map[string][]string{"k3": {"v3", "v4", "v2"}}

	if EqualStrings(m1, m2) {
		t.Errorf("Expected <false>")
	}
}

func TestNotEqualStringsDiffCountSlice(t *testing.T) {
	m1 := map[string][]string{"k1": {"v3", "v4", "v1"}, "k2": {}}
	m2 := map[string][]string{"k2": {}, "k3": {"v3", "v4"}}

	if EqualStrings(m1, m2) {
		t.Errorf("Expected <false>")
	}
}
