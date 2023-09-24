package config

import (
	"fmt"
	"github.com/goclarum/clarum/core/files"
	"gopkg.in/yaml.v3"
	"log/slog"
	"path"
)

var c *Config

type Config struct {
	Profile string
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
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

	configYaml, _ := yaml.Marshal(config)
	slog.Info(fmt.Sprintf("Using the following config:\n[\n%s]", configYaml))
}

func Version() string {
	return version
}

func BaseDir() string {
	return *baseDir
}

//func LoggingLevel() slog.Level {
//	return slog.
//}

func LogLevelName() string {
	return c.Logging.Level
}
