package json

import (
	"fmt"
	"reflect"
	"strings"
)

type Recorder interface {
	AppendFieldName(indent string, fieldName string) Recorder
	AppendValue(indent string, path string, value any, kind reflect.Kind) Recorder
	AppendValidationErrorSignal(message string) Recorder
	AppendMissingFieldErrorSignal(indent string, path string) Recorder
	AppendStartObject(indent string, path string) Recorder
	AppendEndObject(indent string, path string) Recorder
	AppendStartArray(indent string, path string) Recorder
	AppendEndArray(indent string, path string) Recorder
	AppendNewLine() Recorder
	GetLog() string
}

type defaultRecorder struct {
	logResult strings.Builder
}

type noopRecorder struct {
}

func NewDefaultRecorder() Recorder {
	return &defaultRecorder{}
}

func (recorder *defaultRecorder) GetLog() string {
	return recorder.logResult.String()
}

func (recorder *defaultRecorder) AppendFieldName(indent string, fieldName string) Recorder {
	recorder.logResult.WriteString(fmt.Sprintf("%s\"%s\": ", indent, fieldName))
	return recorder
}

func (recorder *defaultRecorder) AppendValue(indent string, path string, value any, kind reflect.Kind) Recorder {
	childOfArray := pathIsArrayChild(path)

	var indentToSet string
	if childOfArray {
		indentToSet = indent
	} else {
		indentToSet = ""
	}

	if kind == reflect.Map {
		recorder.logResult.WriteString(fmt.Sprintf("%sobject,", indentToSet))
	} else if kind == reflect.Slice {
		recorder.logResult.WriteString(fmt.Sprintf("%sarray,", indentToSet))
	} else if kind != reflect.Invalid {
		recorder.logResult.WriteString(fmt.Sprintf("%s%v,", indentToSet, value))
	}
	return recorder
}

func (recorder *defaultRecorder) AppendValidationErrorSignal(message string) Recorder {
	recorder.logResult.WriteString(fmt.Sprintf(" <-- %s\n", message))
	return recorder
}

func (recorder *defaultRecorder) AppendMissingFieldErrorSignal(indent string, path string) Recorder {
	recorder.logResult.WriteString(fmt.Sprintf("%s X-- missing field [%s]\n", indent, path))
	return recorder
}

func (recorder *defaultRecorder) AppendStartObject(indent string, path string) Recorder {
	childOfArray := pathIsArrayChild(path)

	if childOfArray {
		recorder.logResult.WriteString(fmt.Sprintf("%s{", indent))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s{", ""))
	}
	return recorder
}

func (recorder *defaultRecorder) AppendEndObject(indent string, path string) Recorder {
	root := pathIsRoot(path)

	if root {
		recorder.logResult.WriteString(fmt.Sprintf("%s}\n", ""))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s},\n", indent))
	}
	return recorder
}

func (recorder *defaultRecorder) AppendStartArray(indent string, path string) Recorder {
	childOfArray := pathIsArrayChild(path)

	if childOfArray {
		recorder.logResult.WriteString(fmt.Sprintf("%s[", indent))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s[", ""))
	}
	return recorder
}

func (recorder *defaultRecorder) AppendEndArray(indent string, path string) Recorder {
	root := pathIsRoot(path)

	if root {
		recorder.logResult.WriteString(fmt.Sprintf("%s]\n", indent))
	} else {
		recorder.logResult.WriteString(fmt.Sprintf("%s],\n", indent))
	}
	return recorder
}

func (recorder *defaultRecorder) AppendNewLine() Recorder {
	recorder.logResult.WriteString("\n")
	return recorder
}

func (recorder *noopRecorder) GetLog() string {
	return ""
}

func (recorder *noopRecorder) AppendFieldName(indent string, fieldName string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendValue(indent string, path string, value any, kind reflect.Kind) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendValidationErrorSignal(message string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendMissingFieldErrorSignal(indent string, path string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendStartObject(indent string, path string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendEndObject(indent string, path string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendStartArray(indent string, path string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendEndArray(indent string, path string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendNewLine() Recorder {
	return recorder
}
