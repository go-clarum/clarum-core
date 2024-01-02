package config

import (
	"github.com/goclarum/clarum/core/validators/strings"
	"log/slog"
)

const (
	version              string = "v1"
	defaultBaseDir       string = "."
	defaultConfigFile    string = "clarum-properties.yaml"
	defaultLogLevel      string = "info"
	defaultProfile       string = "dev"
	defaultActionTimeout uint   = 10
)

// replace missing attributes from the configuration with their default values
func (config *Config) setDefaults() {
	slog.Debug("Replacing missing values with defaults")

	if strings.IsBlank(config.Profile) {
		config.Profile = defaultProfile
	}
	if strings.IsBlank(config.Logging.Level) {
		config.Logging.Level = defaultLogLevel
	}
	if config.Actions.TimeoutSeconds == 0 {
		config.Actions.TimeoutSeconds = defaultActionTimeout
	}
}
