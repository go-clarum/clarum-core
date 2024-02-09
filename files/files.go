package files

import (
	"errors"
	"github.com/go-clarum/clarum-core/logging"
	"github.com/go-clarum/clarum-core/validators/strings"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadYamlFileToStruct[S any](filePath string) (*S, error) {
	if strings.IsBlank(filePath) {
		logging.Error("Unable to read file. File path is empty")
		return nil, errors.New("file path is empty")
	}

	buf, err := os.ReadFile(filePath)
	if err != nil {
		logging.Errorf("Failed to load file: %s", err.Error())
		return nil, err
	}

	out := new(S)

	if err := yaml.Unmarshal(buf, out); err != nil {
		logging.Errorf("Failed to unmarshal yaml file %s: %s", filePath, err.Error())
		return nil, err
	}

	return out, err
}
