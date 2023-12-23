package config

import (
	"fmt"
	"github.com/goclarum/clarum/core/files"
	"gopkg.in/yaml.v3"
	"log/slog"
	"path"
	"strings"
	"time"
)

var c *Config

type Config struct {
	Profile string
	Actions struct {
		TimeoutSeconds uint
	}
	Logging struct {
		Level string `yaml:"level"`
	}
}

func init() {
	configFilePath := path.Join(*baseDir, *configFile)
	config, err := files.ReadYamlFileToStruct[Config](configFilePath)
	if err != nil {
		slog.Info("Failed to load config file. Default values will be used instead")
		config = &Config{}
	}

	config.setDefaults()
	config.overwriteWithCliFlags()
	c = config

	configYaml, _ := yaml.Marshal(config)
	slog.Info(fmt.Sprintf("Using the following config:\n[\n%s]", configYaml))
}

func Version() string {
	return version
}

func BaseDir() string {
	return *baseDir
}

func LoggingLevel() slog.Level {
	return parseLevel(c.Logging.Level)
}

func ActionTimeout() time.Duration {
	return time.Duration(c.Actions.TimeoutSeconds) * time.Second
}

func parseLevel(level string) slog.Level {
	lcLevel := strings.ToLower(level)
	var result slog.Level

	switch lcLevel {
	case "error":
		result = slog.LevelError
	case "warn":
		result = slog.LevelWarn
	case "debug":
		result = slog.LevelDebug
	default:
		result = slog.LevelInfo
	}

	return result
}
