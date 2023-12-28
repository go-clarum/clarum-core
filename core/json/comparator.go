package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
)

type options struct {
	strictObjectSizeCheck bool
	pathsToIgnore         []string
	logger                *slog.Logger
	recorder              Recorder
}

type Comparator struct {
	options
}

// TODO: documentation
// goroutine safe
// recorder setting
// The problems that we had with a basic implementation:
//   - we only know if they are equal or not, nothing more, no information about why
//   - we cannot use this to ignore fields/values, ex. timestamp values
func (comparator *Comparator) Compare(expected []byte, actual []byte) (string, error) {
	var expectedMap, actualMap map[string]any
	comparator.logger.Debug(fmt.Sprintf("json comparator - comparing [%s] to [%s]", expectedMap, actualMap))

	expectedMap, err1 := toMap(expected)
	if err1 != nil {
		return "", err1
	}

	actualMap, err2 := toMap(actual)
	if err2 != nil {
		return "", err2
	}

	compareErrors := comparator.compareJsonMaps("$", expectedMap, actualMap,
		"", []error{})

	if len(compareErrors) > 0 {
		comparator.logger.Debug(fmt.Sprintf("json comparator - JSON structures do not match"))
	} else {
		comparator.logger.Debug(fmt.Sprintf("json comparator - JSON structures match"))
	}

	return comparator.recorder.GetLog(), errors.Join(compareErrors...)
}

// TODO: implement strictObjectSizeCheck
// todo: what happens when expected and actual each have one different field (size is the same) - actual map has unexpected fields
func (comparator *Comparator) compareJsonMaps(pathParent string, expected map[string]any, actual map[string]any,
	logIndent string, compareErrors []error) []error {
	currIndent := logIndent + "  "

	compareErrors = handleFieldsCheck(pathParent, expected, actual, comparator.strictObjectSizeCheck, comparator.recorder, compareErrors)

	// TODO: implement ignore element by jsonPath & ignore value by using @ignore@
	for key, expectedValue := range expected {
		if actualValue, exists := actual[key]; exists {
			comparator.recorder.AppendFieldName(currIndent, key)

			expectedValueType := reflect.TypeOf(expectedValue)
			actualValueType := reflect.TypeOf(actualValue)

			if expectedValueType.Kind() != actualValueType.Kind() {
				compareErrors = handleTypeMismatch(getJsonPath(pathParent, key),
					expectedValueType, actualValueType, comparator.recorder, compareErrors)
			} else {
				// we only consider JSON Kinds, since the Unmarshal already parsed & checked them
				switch actualValueType.Kind() {
				case reflect.String:
					expectedString := expectedValue.(string)
					actualString := actualValue.(string)

					compareErrors = compareValue(getJsonPath(pathParent, key),
						expectedString != actualString,
						expectedString, actualString, comparator.recorder, "", compareErrors)
				case reflect.Float64:
					compareErrors = compareValue(getJsonPath(pathParent, key),
						expectedValue.(float64) != actualValue.(float64),
						formatFloat(expectedValue), formatFloat(actualValue), comparator.recorder, "", compareErrors)
				case reflect.Bool:
					expectedBool := expectedValue.(bool)
					actualBool := actualValue.(bool)

					compareErrors = compareValue(getJsonPath(pathParent, key),
						expectedBool != actualBool,
						strconv.FormatBool(expectedBool), strconv.FormatBool(actualBool), comparator.recorder, "", compareErrors)
				case reflect.Slice:
					compareErrors = comparator.compareSlices(getJsonPath(pathParent, key),
						expectedValue.([]interface{}), actualValue.([]interface{}),
						comparator.recorder, currIndent, compareErrors)
				case reflect.Map:
					compareErrors = comparator.compareJsonMaps(getJsonPath(pathParent, key),
						expectedValue.(map[string]any), actualValue.(map[string]any),
						currIndent, compareErrors)
				}
			}
		} else {
			compareErrors = handleMissingField(getJsonPath(pathParent, key),
				currIndent, comparator.recorder, compareErrors)
		}
	}

	comparator.recorder.AppendEndObject(logIndent)
	return compareErrors
}

// Arrays in json are represented as slices of type interface because they can contain anything.
// Each item in the slice can be of any valid JSON type.
func (comparator *Comparator) compareSlices(path string, expected []interface{}, actual []interface{},
	recorder Recorder, currIndent string, compareErrors []error) []error {
	recorder.AppendStartArray()

	expectedLen := len(expected)
	actualLen := len(actual)
	if expectedLen != actualLen {
		recorder.AppendValidationErrorSignal(fmt.Sprintf("size mismatch - expected [%d]", expectedLen))
		return append(compareErrors,
			errors.New(fmt.Sprintf("[%s] - array size mismatch - expected [%d] but received [%d]", path, expectedLen, actualLen)))
	} else {
		recorder.AppendNewLine()
	}

	valIdent := currIndent + "  "
	for i, expectedValue := range expected {
		expectedValueType := reflect.TypeOf(expectedValue)
		actualValue := actual[i]
		actualValueType := reflect.TypeOf(actualValue)

		if expectedValueType.Kind() != actualValueType.Kind() {
			recorder.AppendValue(valIdent, actualValue, actualValueType.Kind())
			baseErrorMessage := fmt.Sprintf("value type mismatch - expected [%s] but found [%s]",
				convertToJsonType(expectedValueType), convertToJsonType(actualValueType))

			compareErrors = append(compareErrors, errors.New(fmt.Sprintf("[%s] - %s", getJsonPathArray(path, i),
				baseErrorMessage)))
			recorder.AppendValidationErrorSignal(baseErrorMessage)
		} else {
			switch actualValueType.Kind() {
			case reflect.String:
				expectedString := expectedValue.(string)
				actualString := actualValue.(string)

				compareErrors = compareValue(getJsonPathArray(path, i),
					expectedString != actualString,
					expectedString, actualString, recorder, valIdent, compareErrors)
			case reflect.Float64:
				compareErrors = compareValue(getJsonPathArray(path, i),
					expectedValue.(float64) != actualValue.(float64),
					formatFloat(expectedValue), formatFloat(actualValue), recorder, valIdent, compareErrors)
			case reflect.Bool:
				expectedBool := expectedValue.(bool)
				actualBool := actualValue.(bool)

				compareErrors = compareValue(getJsonPathArray(path, i),
					expectedBool != actualBool,
					strconv.FormatBool(expectedBool), strconv.FormatBool(actualBool), recorder, valIdent, compareErrors)
			case reflect.Slice:
				compareErrors = comparator.compareSlices(getJsonPathArray(path, i),
					expectedValue.([]interface{}), actualValue.([]interface{}),
					recorder, currIndent, compareErrors)
			case reflect.Map:
				compareErrors = comparator.compareJsonMaps(getJsonPathArray(path, i),
					expectedValue.(map[string]any), actualValue.(map[string]any),
					currIndent, compareErrors)
			}
		}
	}
	recorder.AppendEndArray(currIndent)
	return compareErrors
}

// When describing types we have to consider that the users will think of JSON types when reading logs & error messages.
// This is why we have to translate Go types into JSON types.
//
// json.Unmarshal returns a map[string]interface{} with all the fields of the JSON object:
// - number is a reflect.Float64
// - string is a reflect.String
// - boolean is a reflect.Bool
// - array is a reflect.Slice
// - struct is a reflect.Map
func convertToJsonType(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Float64:
		return "number"
	case reflect.Map:
		return "object"
	case reflect.Slice:
		return "array"
	default:
		return goType.String()
	}
}

func handleFieldsCheck(pathParent string, expected map[string]any, actual map[string]any, strictSizeCheck bool,
	recorder Recorder, compareErrors []error) []error {
	if strictSizeCheck && len(expected) != len(actual) {
		recorder.AppendStartObject().
			AppendValidationErrorSignal("number of fields does not match")

		compareErrors = append(compareErrors,
			errors.New(fmt.Sprintf("[%s] - number of fields does not match", pathParent)))
	} else {
		recorder.AppendStartObject().AppendNewLine()
	}
	return compareErrors
}

func handleTypeMismatch(path string, expectedValueType reflect.Type, actualValueType reflect.Type,
	recorder Recorder, compareErrors []error) []error {

	baseErrorMessage := fmt.Sprintf("type mismatch - expected [%s] but found [%s]",
		convertToJsonType(expectedValueType), convertToJsonType(actualValueType))

	compareErrors = append(compareErrors, errors.New(fmt.Sprintf("[%s] - %s", path, baseErrorMessage)))
	recorder.AppendValidationErrorSignal(baseErrorMessage)

	return compareErrors
}

func compareValue(path string, mismatch bool, expectedValue string, actualValue string, recorder Recorder,
	indent string, compareErrors []error) []error {
	recorder.AppendValue(indent, actualValue, reflect.String)

	if mismatch {
		compareErrors = append(compareErrors,
			errors.New(fmt.Sprintf("[%s] - value mismatch - expected [%s] but received [%s]", path, expectedValue, actualValue)))
		recorder.AppendValidationErrorSignal(fmt.Sprintf("value mismatch - expected [%s]", expectedValue))
	} else {
		recorder.AppendNewLine()
	}

	return compareErrors
}

func handleMissingField(path string, indent string, recorder Recorder, compareErrors []error) []error {
	compareErrors = append(compareErrors, errors.New(fmt.Sprintf("[%s] - field is missing", path)))
	recorder.AppendMissingFieldErrorSignal(indent, path)

	return compareErrors
}

// We rely on json.Unmarshal to detect invalid json structures here.
func toMap(rawJson []byte) (map[string]any, error) {
	var result any
	if err := json.Unmarshal(rawJson, &result); err != nil {
		return nil, handleError("unable to parse JSON - error [%s] - from string [%s]", err, rawJson)
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

func getJsonPathArray(pathParent string, index int) string {
	return fmt.Sprintf("%s[%d]", pathParent, index)
}

func handleError(format string, a ...any) error {
	errorMessage := fmt.Sprintf(format, a...)
	return errors.New(errorMessage)
}
