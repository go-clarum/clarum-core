package files

import (
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/validators/strings"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

func ReadYamlFileToStruct[S any](filePath string) (*S, error) {
	if strings.IsBlank(filePath) {
		slog.Error("Unable to read file. File path is empty")
		return nil, errors.New("file path is empty")
	}

	buf, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load file: %s", err.Error()))
		return nil, err
	}

	out := new(S)

	if err := yaml.Unmarshal(buf, out); err != nil {
		slog.Error(fmt.Sprintf("Failed to unmarshal yaml file %s: %s", filePath, err.Error()))
		return nil, err
	}

	return out, err
}
