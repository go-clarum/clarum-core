package comparator

import (
	"fmt"
	"strings"
	"testing"
)

func TestInvalidExpectedJson(t *testing.T) {
	expectedError := "unable to parse JSON - error [invalid character '}' in literal true (expecting 'e')] - from string [{\"active\": tru}]"

	expectedValue := []byte("{" +
		"\"active\": tru" +
		"}")
	actualValue := []byte("{" +
		"\"active\": true" +
		"}")

	testComparator(t, expectedValue, actualValue, expectedError, "")
}

func TestInvalidActualJson(t *testing.T) {
	expectedError := "unable to parse JSON - error [invalid character '}' looking for beginning of value] - from string [{\"active\": true,\"aliases\": [\"Batman\",}]"

	expectedValue := []byte("{" +
		"\"active\": true" +
		"}")
	actualValue := []byte("{" +
		"\"active\": true," +
		"\"aliases\": [\"Batman\"," +
		"}")

	testComparator(t, expectedValue, actualValue, expectedError, "")
}

func TestEmptyObject(t *testing.T) {
	expectedValue := []byte("{}")
	actualValue := []byte("{}")
	expectedRecorderLog := "{\n}\n"

	testComparator(t, expectedValue, actualValue, "", expectedRecorderLog)
}

func TestExpectEmptyObject(t *testing.T) {
	expectedError := "[$] - number of fields does not match"

	expectedValue := []byte("{}")
	actualValue := []byte("{" +
		"\"active\": true" +
		"}")

	expectedRecorderLog := "{ <-- number of fields does not match\n}\n"

	testComparator(t, expectedValue, actualValue, expectedError, expectedRecorderLog)
}

func TestReceiveEmptyObject(t *testing.T) {
	expectedError := "[$] - number of fields does not match\n" +
		"[$.active] - field is missing"

	expectedValue := []byte("{" +
		"\"active\": true" +
		"}")
	actualValue := []byte("{}")

	expectedRecorderLog := "{ <-- number of fields does not match\n   X-- missing field [active]\n}\n"

	testComparator(t, expectedValue, actualValue, expectedError, expectedRecorderLog)
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

	recorderLog := testComparator(t, expectedValue, actualValue, expectedError, "")

	if !strings.Contains(recorderLog, "  \"location\": { <-- number of fields does not match\n") {
		t.Error("missing: number of fields does not match errors")
	}
	if !strings.Contains(recorderLog, "     X-- missing field [street]\n") {
		t.Error("missing: missing field [street]")
	}
	if !strings.Contains(recorderLog, "     X-- missing field [number]\n") {
		t.Error("missing: missing field [number]")
	}
}

func TestNotStrictFieldCheck(t *testing.T) {
	expectedError := "[$.location.street] - field is missing\n" +
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

	expectedRecorderLog := "{\n" +
		"  \"active\": true,\n" +
		"  \"location\": {\n" +
		"     X-- missing field [street]\n" +
		"     X-- missing field [number]\n" +
		"  },\n" +
		"}\n"

	comparator := Builder().
		StrictObjectSizeCheck(false).
		Recorder(NewDefaultRecorder()).
		Comparator()
	recorderResult, err := comparator.Compare(expectedValue, actualValue)

	checkError(t, err, expectedError)
	checkRecorderLog(t, expectedRecorderLog, recorderResult)
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

	expectedRecorderLog := "{ <-- number of fields does not match\n" +
		"  \"active\": true,\n" +
		"   X-- missing field [location]" +
		"\n}\n"

	testComparator(t, expectedValue, actualValue, expectedError, expectedRecorderLog)
}

func TestOKValidationAllTypes(t *testing.T) {
	expectedValue := []byte("{" +
		"\"active\": true," +
		" \"name\": \"Bruce Wayne\"," +
		" \"age\": 38," +
		" \"height\": 1.879," +
		"\"aliases\": [" +
		"\"Batman\"," +
		"\"The Dark Knight\"" +
		"]," +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}" +
		"}")

	actualValue := []byte("{" +
		"\"active\": true," +
		" \"name\": \"Bruce Wayne\"," +
		"\"aliases\": [" +
		"\"Batman\"," +
		"\"The Dark Knight\"" +
		"]," +
		" \"age\": 38," +
		" \"height\": 1.879," +
		"\"location\": {" +
		"\"street\": \"Mountain Drive\"," +
		"\"number\": 1007," +
		"\"hidden\": false" +
		"}" +
		"}")

	// we ignore the recorder log because the order of the elements is always different
	testComparator(t, expectedValue, actualValue, "", "")
}

// flaky test because the order of fields inside the JSON object changes on unmarshalling
func TestErrorValidationAllTypes(t *testing.T) {
	expectedError := "[$.name] - value mismatch - expected [Bruce] but received [Bruce Wayne]\n" +
		"[$.age] - value mismatch - expected [37] but received [38]\n" +
		"[$.location.street] - field is missing\n" +
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
		"\"address\": \"Mountain Drive\"," +
		"\"number\": 1008," +
		"\"hidden\": true" +
		"}" +
		"}")

	recorderLog := testComparator(t, expectedValue, actualValue, expectedError, "")

	if !strings.Contains(recorderLog, "  \"name\": Bruce Wayne, <-- value mismatch - expected [Bruce]\n") {
		t.Error("missing: expected [Bruce]")
	}
	if !strings.Contains(recorderLog, "  \"age\": 38, <-- value mismatch - expected [37]\n") {
		t.Error("missing: expected [37]")
	}
	if !strings.Contains(recorderLog, "    \"number\": 1008, <-- value mismatch - expected [1007]\n") {
		t.Error("missing: expected [1007]")
	}
	if !strings.Contains(recorderLog, "    \"hidden\": true, <-- value mismatch - expected [false]\n") {
		t.Error("missing: expected [false]")
	}
	if !strings.Contains(recorderLog, "     X-- missing field [street]\n") {
		t.Error("missing: expected [false]")
	}
}

func TestKindValidationBooleanType(t *testing.T) {
	expectedError := "[$.active] - type mismatch - expected [boolean] but found [string]"

	expectedValue := []byte("{" +
		"\"active\": true" +
		"}")
	actualValue := []byte("{" +
		"\"active\": \"true\"" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"active\":  <-- type mismatch - expected [boolean] but found [string]\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedError, expectedRecorderLog)
}

func TestKindValidationNumberType(t *testing.T) {
	expectedError := "[$.age] - type mismatch - expected [string] but found [number]"

	expectedValue := []byte("{" +
		" \"age\": \"38\"" +
		"}")
	actualValue := []byte("{" +
		" \"age\": 38" +
		"}")

	expectedRecorderLog := "{\n" +
		"  \"age\":  <-- type mismatch - expected [string] but found [number]\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedError, expectedRecorderLog)
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

	expectedRecorderLog := "{\n" +
		"  \"location\":  <-- type mismatch - expected [string] but found [object]\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedError, expectedRecorderLog)
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

	expectedRecorderLog := "{\n" +
		"  \"aliases\":  <-- type mismatch - expected [string] but found [array]\n" +
		"}\n"

	testComparator(t, expectedValue, actualValue, expectedError, expectedRecorderLog)

}

func checkError(t *testing.T, err error, expectedError string) {
	if len(expectedError) == 0 && err != nil { // no error expected
		t.Error(err)
	} else if len(expectedError) > 0 && err == nil {
		t.Error("Error expected, but there was none.")
	} else if len(expectedError) > 0 && err.Error() != expectedError {
		t.Error(err)
	}
}

func testComparator(t *testing.T, expectedValue []byte, actualValue []byte, expectedError string,
	expectedRecorderLog string) string {
	comparator := Builder().Recorder(NewDefaultRecorder()).Comparator()
	recorderResult, err := comparator.Compare(expectedValue, actualValue)

	checkError(t, err, expectedError)
	checkRecorderLog(t, expectedRecorderLog, recorderResult)

	return recorderResult
}

func checkRecorderLog(t *testing.T, expected string, actual string) {
	fmt.Println(actual)
	if len(expected) > 0 && expected != actual {
		t.Error("Recorder log does not match")
	}
}
