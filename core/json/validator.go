package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
)

// TODO: documentation
// The problems that we had with a basic implementation:
//   - we only know if they are equal or not, nothing more, no information about why
//   - we cannot use this to ignore fields/values, ex. timestamp values
func Compare(expected []byte, actual []byte) (string, error) {
	var expectedMap, actualMap map[string]any
	slog.Debug(fmt.Sprintf("comparing [%s] to [%s]", expectedMap, actualMap))

	expectedMap, err1 := toMap(expected)
	if err1 != nil {
		return "", err1
	}

	actualMap, err2 := toMap(actual)
	if err2 != nil {
		return "", err2
	}

	var logResult strings.Builder

	compareErrors := compareJsonMaps("$", expectedMap, actualMap,
		true, &logResult, "", []error{})

	return logResult.String(), errors.Join(compareErrors...)
}

// todo: what happens when expected and actual each have one different field - unexpected field validation
func compareJsonMaps(pathParent string, expected map[string]any, actual map[string]any,
	strictSizeCheck bool, logResult *strings.Builder, logIndent string, compareErrors []error) []error {
	currentIndentation := logIndent + "  "

	if strictSizeCheck && len(expected) != len(actual) {
		logResult.WriteString("{ <-- number of fields does not match\n")
		compareErrors = append(compareErrors,
			errors.New(fmt.Sprintf("%s - number of fields does not match", pathParent)))
	} else {
		logResult.WriteString("{\n")
	}

	// TODO: implement ignore element by jsonPath & ignore value by using @ignore@
	for key, expectedValue := range expected {
		if actualValue, exists := actual[key]; exists {
			logResult.WriteString(fmt.Sprintf("%s\"%s\": ", currentIndentation, key))

			expectedValueType := reflect.TypeOf(expectedValue)
			actualValueType := reflect.TypeOf(actualValue)

			if expectedValueType.Kind() != actualValueType.Kind() {
				compareErrors = handleTypeMismatch(getJsonPath(pathParent, key),
					expectedValueType, actualValueType, logResult, compareErrors)
			} else {
				// we only consider JSON Kinds, since the Unmarshal already parsed & checked them
				switch actualValueType.Kind() {
				case reflect.Array:
					// TODO: impl array handling
				case reflect.String:
					expectedString := expectedValue.(string)
					actualString := actualValue.(string)

					compareErrors = handleValue(getJsonPath(pathParent, key),
						expectedString != actualString,
						expectedString, actualString, logResult, compareErrors)
				case reflect.Float64:
					compareErrors = handleValue(getJsonPath(pathParent, key),
						expectedValue.(float64) != actualValue.(float64),
						formatFloat(expectedValue), formatFloat(actualValue), logResult, compareErrors)
				case reflect.Bool:
					expectedBool := expectedValue.(bool)
					actualBool := actualValue.(bool)

					compareErrors = handleValue(getJsonPath(pathParent, key),
						expectedBool != actualBool,
						strconv.FormatBool(expectedBool), strconv.FormatBool(actualBool), logResult, compareErrors)
				case reflect.Map:
					compareErrors = compareJsonMaps(getJsonPath(pathParent, key),
						expectedValue.(map[string]any), actualValue.(map[string]any),
						strictSizeCheck, logResult, currentIndentation, compareErrors)
				}
			}
		} else {
			logResult.WriteString(fmt.Sprintf("%s X-- missing field '%s'\n", currentIndentation, key))
		}
	}

	logResult.WriteString(fmt.Sprintf("%s}\n", logIndent))
	return compareErrors
}

func handleTypeMismatch(path string, expectedValueType reflect.Type, actualValueType reflect.Type,
	logResult *strings.Builder, compareErrors []error) []error {

	baseErrorMessage := fmt.Sprintf("type mismatch - expected [%s] but found [%s]",
		convertToJsonType(expectedValueType), convertToJsonType(actualValueType))

	compareErrors = append(compareErrors, errors.New(fmt.Sprintf("%s - %s", path, baseErrorMessage)))
	logResult.WriteString(fmt.Sprintf(" <-- %s\n", baseErrorMessage))

	return compareErrors
}

// When describing types we have to consider that the users will think of JSON types when reading logs & error messages.
// This is why we have to translate Go types into JSON types.
func convertToJsonType(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Float64:
		return "number"
	default:
		return goType.String()
	}
}

func handleValue(path string, mismatch bool, expectedValue string, actualValue string, logResult *strings.Builder,
	compareErrors []error) []error {
	logResult.WriteString(fmt.Sprintf("%s,", actualValue))

	if mismatch {
		compareErrors = append(compareErrors,
			errors.New(fmt.Sprintf("%s - value mismatch - expected [%s] but received [%s]", path, expectedValue, actualValue)))
		logResult.WriteString(fmt.Sprintf(" <-- value mismatch - expected [%s]", expectedValue))
	}

	logResult.WriteString("\n")
	return compareErrors
}

// We rely on json.Unmarshal to detect invalid json structures
// json.Unmarshal returns a map[string]interface{} with all the fields of the JSON object
//
// number is a reflect.Float64
// string is a reflect.String
// boolean is a reflect.Bool
// array is a reflect.Array
// struct is a reflect.Struct
func toMap(rawJson []byte) (map[string]any, error) {
	var result any
	if err := json.Unmarshal(rawJson, &result); err != nil {
		return nil, handleError("unable to parse JSON - error [%s] - from string [%s]", rawJson, err)
	}

	return result.(map[string]any), nil
}

// We have to trim trailing zeroes from the parsed float64 number before logging them.
func formatFloat(expectedValue any) string {
	return strconv.FormatFloat(expectedValue.(float64), 'f', -1, 64)
}

func getJsonPath(pathParent string, key string) string {
	return fmt.Sprintf("%s.%s", pathParent, key)
}

func handleError(format string, a ...any) error {
	errorMessage := fmt.Sprintf(format, a...)
	slog.Error(errorMessage)
	return errors.New(errorMessage)
}
