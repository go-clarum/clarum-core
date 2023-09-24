package config

import (
	"github.com/goclarum/clarum/core/validators/strings"
	"log/slog"
)

const (
	version           string = "v1"
	defaultBaseDir    string = "."
	defaultConfigFile string = "clarum-properties.yaml"
	defaultLogOutput  string = "clarum-tests-logs"
	defaultLogLevel   string = "info"
	defaultProfile    string = "default"
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
	if strings.IsBlank(config.Logging.Output) {
		config.Logging.Output = defaultLogOutput
	}
}
