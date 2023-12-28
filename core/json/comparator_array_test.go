package json

import (
	"fmt"
	"testing"
)

func TestStrictEmptyInExpectedJson(t *testing.T) {
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

func TestStrictEmptyInActualJson(t *testing.T) {
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

func TestStrictTypeMismatchJson(t *testing.T) {
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

func TestStrictStringValidation(t *testing.T) {
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

func TestStrictNumberValidation(t *testing.T) {
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
func TestStrictBoolOKValidation(t *testing.T)   {}
func TestStrictArrayOKValidation(t *testing.T)  {}
func TestStrictObjectOKValidation(t *testing.T) {}

func TestNotStrictEmptyInExpectedJson(t *testing.T) {
	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"" +
		"]" +
		"}")

	comparator := Builder().
		StrictArrayCheck(false).
		Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	if err != nil {
		t.Error(err)
	}
}

func TestNotStrictEmptyInActualJson(t *testing.T) {
	expectedError := "[$.aliases] - value [Batman] is missing"

	expectedValue := []byte("{" +
		"\"aliases\": [" +
		"\"Batman\"" +
		"]" +
		"}")
	actualValue := []byte("{" +
		"\"aliases\": [" +
		"]" +
		"}")

	comparator := Builder().
		StrictArrayCheck(false).
		Comparator()
	logResult, err := comparator.Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	checkError(t, err, expectedError)
}

func TestNotStrictStringOKValidation(t *testing.T)    {}
func TestNotStrictNumberOKValidation(t *testing.T)    {}
func TestNotStrictBoolOKValidation(t *testing.T)      {}
func TestNotStrictArrayOKValidation(t *testing.T)     {}
func TestNotStrictObjectOKValidation(t *testing.T)    {}
func TestNotStrictObjectErrorValidation(t *testing.T) {}
