package json

import (
	"fmt"
	"testing"
)

func TestInvalidExpectedJSON(t *testing.T) {
	expectedError := "unable to parse JSON - error [invalid character '}' in literal true (expecting 'e')] - from string [{\"active\": tru}]"

	expectedValue := []byte("{" +
		"\"active\": tru" +
		"}")
	actualValue := []byte("{" +
		"\"active\": true" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestInvalidActualJSON(t *testing.T) {
	expectedError := "unable to parse JSON - error [invalid character '}' looking for beginning of value] - from string [{\"active\": true,\"aliases\": [\"Batman\",}]"

	expectedValue := []byte("{" +
		"\"active\": true" +
		"}")
	actualValue := []byte("{" +
		"\"active\": true," +
		"\"aliases\": [\"Batman\"," +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestEmptyObject(t *testing.T) {
	expectedValue := []byte("{}")
	actualValue := []byte("{}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	if err != nil {
		t.Error(err)
	}
}

func TestExpectEmptyObject(t *testing.T) {
	expectedError := "[$] - number of fields does not match"

	expectedValue := []byte("{}")
	actualValue := []byte("{" +
		"\"active\": true" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestReceiveEmptyObject(t *testing.T) {
	expectedError := "[$] - number of fields does not match\n" +
		"[$.active] - field is missing"

	expectedValue := []byte("{" +
		"\"active\": true" +
		"}")
	actualValue := []byte("{}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestDeepEmptyObject(t *testing.T) {
	expectedError := "[$.location] - number of fields does not match\n" +
		"[$.location.street] - field is missing\n" +
		"[$.location.number] - field is missing"

	expectedValue := []byte("{" +
		"\"active\": true," +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007" +
		"}" +
		"}")
	actualValue := []byte("{" +
		"\"active\": true," +
		"\"location\": {" +
		"}" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestMissingObject(t *testing.T) {
	expectedError := "[$] - number of fields does not match\n" +
		"[$.location] - field is missing"

	expectedValue := []byte("{" +
		"\"active\": true," +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007" +
		"}" +
		"}")
	actualValue := []byte("{" +
		"\"active\": true" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

// TODO: add array
func TestOKValidationAllTypes(t *testing.T) {
	expectedValue := []byte("{" +
		"\"active\": true," +
		" \"name\": \"Bruce Wayne\"," +
		" \"age\": 38," +
		" \"height\": 1.879," +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}" +
		"}")

	actualValue := []byte("{" +
		"\"active\": true," +
		" \"name\": \"Bruce Wayne\"," +
		" \"age\": 38," +
		" \"height\": 1.879," +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	if err != nil {
		t.Error(err)
	}
}

// flaky test because of the order of fields inside the JSON object
func TestErrorValidationAllTypes(t *testing.T) {
	expectedError := "[$.name] - value mismatch - expected [Bruce] but received [Bruce Wayne]\n" +
		"[$.age] - value mismatch - expected [37] but received [38]\n" +
		"[$.location.number] - value mismatch - expected [1007] but received [1008]\n" +
		"[$.location.hidden] - value mismatch - expected [false] but received [true]"

	expectedValue := []byte("{" +
		"\"active\": true," +
		" \"name\": \"Bruce\"," +
		" \"age\": 37," +
		" \"height\": 1.879," +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}" +
		"}")

	actualValue := []byte("{" +
		"\"active\": true," +
		" \"name\": \"Bruce Wayne\"," +
		" \"age\": 38," +
		" \"height\": 1.879," +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1008," +
		"\"hidden\": true" +
		"}" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestKindValidationBooleanType(t *testing.T) {
	expectedError := "[$.active] - type mismatch - expected [boolean] but found [string]"

	expectedValue := []byte("{" +
		"\"active\": true" +
		"}")
	actualValue := []byte("{" +
		"\"active\": \"true\"" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestKindValidationNumberType(t *testing.T) {
	expectedError := "[$.age] - type mismatch - expected [string] but found [number]"

	expectedValue := []byte("{" +
		" \"age\": \"38\"" +
		"}")
	actualValue := []byte("{" +
		" \"age\": 38" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestKindValidationObjectType(t *testing.T) {
	expectedError := "[$.location] - type mismatch - expected [string] but found [object]"

	expectedValue := []byte("{" +
		" \"location\": \"Mountain Drive\"" +
		"}")
	actualValue := []byte("{" +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestKindValidationArrayType(t *testing.T) {
	expectedError := "[$.aliases] - type mismatch - expected [string] but found [array]"

	expectedValue := []byte("{" +
		" \"aliases\": \"Batman\"" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"" +
		"]" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func checkError(t *testing.T, err error, expectedError string) {
	if err == nil {
		t.Error("Error expected, but there was none.")
	}
	if err.Error() != expectedError {
		t.Error(err)
	}
}
