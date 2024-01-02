package strings

import "testing"

func TestIsBlank(t *testing.T) {
	if !IsBlank("") {
		t.Errorf("Expected <true>")
	}
	if !IsBlank(" ") {
		t.Errorf("Expected <true>")
	}
	if IsBlank("batman") {
		t.Errorf("Expected <false>")
	}
	if IsBlank("   batman") {
		t.Errorf("Expected <false>")
	}
}

func TestIsNotBlank(t *testing.T) {
	if IsNotBlank("") {
		t.Errorf("Expected <false>")
	}
	if IsNotBlank(" ") {
		t.Errorf("Expected <false>")
	}
	if !IsNotBlank("batman") {
		t.Errorf("Expected <true>")
	}
	if !IsNotBlank("   batman") {
		t.Errorf("Expected <true>")
	}
}
