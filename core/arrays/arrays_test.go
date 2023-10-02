package arrays

import (
	"testing"
)

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
