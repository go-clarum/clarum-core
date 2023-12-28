package json

import (
	"fmt"
	"testing"
)

func TestEmptyInExpectedJson(t *testing.T) {
	expectedError := "[$.aliases] - array size mismatch - expected [0] but received [1]"

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"" +
		"]" +
		"}")

	comparator := Builder().Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestEmptyInActualJson(t *testing.T) {
	expectedError := "[$.aliases] - array size mismatch - expected [1] but received [0]"

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"]" +
		"}")

	comparator := Builder().Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestTypeMismatchJson(t *testing.T) {
	expectedError := "[$.aliases[1]] - value type mismatch - expected [string] but found [number]\n" +
		"[$.aliases[2]] - value type mismatch - expected [string] but found [object]\n" +
		"[$.aliases[3]] - value type mismatch - expected [string] but found [array]"

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"," +
		"\"The Dark Knight\"," +
		"\"Batsy\"," +
		"\"The Gotham Guardian\"" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"," +
		"123," +
		"{" +
		"\"someStringKey\": \"someValue\"," +
		"\"someNumberKey\": 123" +
		"}," +
		"[1,2,3]" +
		"]" +
		"}")

	comparator := Builder().Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestStringValidation(t *testing.T) {
	expectedError := "[$.aliases[1]] - value mismatch - expected [The Dark Knight] but received [Robin]"

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"," +
		"\"The Dark Knight\"" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"," +
		"\"Robin\"" +
		"]" +
		"}")

	comparator := Builder().Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestNumberValidation(t *testing.T) {
	expectedError := "[$.measures[1]] - value mismatch - expected [82] but received [83]\n" +
		"[$.measures[3]] - value mismatch - expected [64.1] but received [64.2]"

	expectedValue := []byte("{" +
		"\"measures\": [" +
		"45," +
		"82," +
		"32.2," +
		"64.1" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"measures\": [" +
		"45," +
		"83," +
		"32.2," +
		"64.2" +
		"]" +
		"}")

	comparator := Builder().Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestBoolValidation(t *testing.T) {
	expectedError := "[$.someBooleanArray[1]] - value mismatch - expected [true] but received [false]"

	expectedValue := []byte("{" +
		"\"someBooleanArray\": [" +
		"true," +
		"true" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"someBooleanArray\": [" +
		"true," +
		"false" +
		"]" +
		"}")

	comparator := Builder().Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestDeepArrayValidation(t *testing.T) {
	expectedError := "[$.parent[1][1]] - value type mismatch - expected [string] but found [number]"

	expectedValue := []byte("{" +
		"\"parent\": [" +
		"[" +
		"\"child11\"," +
		"\"child12\"" +
		"]," +
		"[" +
		"\"child21\"," +
		"\"child22\"" +
		"]" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"parent\": [" +
		"[" +
		"\"child11\"," +
		"\"child12\"" +
		"]," +
		"[" +
		"\"child21\"," +
		"123" +
		"]" +
		"]" +
		"}")

	comparator := Builder().Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestObjectValidation(t *testing.T) {
	expectedError := "[$.addresses[0].number] - value mismatch - expected [1007] but received [1035]\n" +
		"[$.addresses[1].hidden] - value mismatch - expected [true] but received [false]"

	expectedValue := []byte("{" +
		"\"addresses\": [" +
		"{" +
		"\"name\": \"Home\"," +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}," +
		"{" +
		"\"name\": \"Batcave\"," +
		"\"street\": \"unknown\"," +
		"\"number\": 0," +
		"\"hidden\": true" +
		"}" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"addresses\": [" +
		"{" +
		"\"name\": \"Home\"," +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1035," +
		"\"hidden\": false" +
		"}," +
		"{" +
		"\"name\": \"Batcave\"," +
		"\"street\": \"unknown\"," +
		"\"number\": 0," +
		"\"hidden\": false" +
		"}" +
		"]" +
		"}")

	comparator := Builder().Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}
