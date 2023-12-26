package json

import (
	"fmt"
	"testing"
)

// to test
// empty object
// array everywhere + empty

func TestValidator(t *testing.T) {
	expectedValue := []byte("{" +
		"\"myBoolean\": true," +
		" \"myString\": \"string value\"," +
		" \"myNumber\": 12345," +
		" \"myDecimal\": 0.88," +
		" \"mySecondString\": \"string value2\"," +
		"\"myObject\": {" +
		"\"someNumber\": 1.23," +
		"\"myStringInObject\": \"string in object value\"" +
		"}" +
		"}")

	actualValue := []byte("{" +
		"\"myBoolean\": \"true\"," +
		" \"myString\": \"string value\"," +
		" \"myNumber\": \"12344\"," +
		" \"myDecimal\": 0.88," +
		"\"myObject\": {" +
		"\"myStringInObject\": \"string in object value\"," +
		"\"someNumber\": 123," +
		"\"secondStringInObject\": \"second string in object\"" +
		"}" +
		"}")

	logResult, err := Compare(expectedValue, actualValue)
	fmt.Println(logResult)

	if err != nil {
		t.Error(err)
	}
}
