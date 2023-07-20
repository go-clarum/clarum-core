package arrays

import (
	"testing"
)

func TestEqualNilAndEmpty(t *testing.T) {
	var a1 []string = nil
	var a2 = make([]string, 0)

	if !Equal(a1, a2) {
		t.Errorf("Expected <true>")
	}
}

func TestEqual(t *testing.T) {
	a1 := []string{"s1", "s2", "s3"}
	a2 := []string{"s3", "s1", "s2"}

	if !Equal(a1, a2) {
		t.Errorf("Expected <true>")
	}
}

func TestNotEqualDiffValue(t *testing.T) {
	a1 := []string{"s1", "s2", "s3"}
	a2 := []string{"s3", "s4", "s2"}

	if Equal(a1, a2) {
		t.Errorf("Expected <false>")
	}
}

func TestNotEqualDiffCount(t *testing.T) {
	a1 := []string{"s1", "s2"}
	a2 := []string{"s3", "s4", "s2"}

	if Equal(a1, a2) {
		t.Errorf("Expected <false>")
	}
}

func TestContainsNilAndEmpty(t *testing.T) {
	var a1 []string = nil
	s1 := "s1"

	if Contains(a1, s1) {
		t.Errorf("Expected <false>")
	}

	a2 := []string{"s1", "s2"}
	s2 := ""

	if Contains(a2, s2) {
		t.Errorf("Expected <false>")
	}
}

func TestContains(t *testing.T) {
	a := []string{"s1", "s2"}
	s := "s1"

	if !Contains(a, s) {
		t.Errorf("Expected <true>")
	}
}

func TestNotContains(t *testing.T) {
	a := []string{"s1", "s2"}
	s := "s21"

	if Contains(a, s) {
		t.Errorf("Expected <false>")
	}
}
