package durations

import (
	"testing"
	"time"
)

func TestValueOverDefault(t *testing.T) {
	result := GetDurationWithDefault(500*time.Millisecond, 2*time.Second)

	if result != 500*time.Millisecond {
		t.Errorf("expected 500ms configured time but received %v", result)
	}
}

func TestDefaultOverValue(t *testing.T) {
	result := GetDurationWithDefault(0, 2*time.Second)

	if result != 2*time.Second {
		t.Errorf("expected 2s configured time but received %v", result)
	}
}
