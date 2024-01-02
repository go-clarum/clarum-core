package config

import "testing"

func TestDefaults(t *testing.T) {
	config := &Config{}
	config.setDefaults()

	if config.Profile != "dev" {
		t.Errorf("Expected <dev>")
	}
	if config.Logging.Level != "info" {
		t.Errorf("Expected <info>")
	}
	if config.Actions.TimeoutSeconds != 10 {
		t.Errorf("Expected <10>")
	}
}
