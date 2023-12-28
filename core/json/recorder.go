package json

import (
	"fmt"
	"reflect"
	"strings"
)

type Recorder interface {
	AppendFieldName(indent string, fieldName string) Recorder
	AppendValue(indent string, value any, kind reflect.Kind) Recorder
	AppendValidationErrorSignal(message string) Recorder
	AppendMissingFieldErrorSignal(indent string, path string) Recorder
	AppendStartObject() Recorder
	AppendEndObject(indent string) Recorder
	AppendStartArray() Recorder
	AppendEndArray(indent string) Recorder
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

func (recorder *defaultRecorder) AppendValue(indent string, value any, kind reflect.Kind) Recorder {
	if kind == reflect.Map {
		recorder.logResult.WriteString(fmt.Sprintf("%sobject,", indent))
	} else if kind == reflect.Slice {
		recorder.logResult.WriteString(fmt.Sprintf("%sarray,", indent))
	} else if kind != reflect.Invalid {
		recorder.logResult.WriteString(fmt.Sprintf("%s%v,", indent, value))
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

func (recorder *defaultRecorder) AppendStartObject() Recorder {
	recorder.logResult.WriteString("{")
	return recorder
}

func (recorder *defaultRecorder) AppendEndObject(indent string) Recorder {
	recorder.logResult.WriteString(fmt.Sprintf("%s}\n", indent))
	return recorder
}

func (recorder *defaultRecorder) AppendStartArray() Recorder {
	recorder.logResult.WriteString("[")
	return recorder
}

func (recorder *defaultRecorder) AppendEndArray(indent string) Recorder {
	recorder.logResult.WriteString(fmt.Sprintf("%s]\n", indent))
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

func (recorder *noopRecorder) AppendValue(indent string, value any, kind reflect.Kind) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendValidationErrorSignal(message string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendMissingFieldErrorSignal(indent string, path string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendStartObject() Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendEndObject(indent string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendStartArray() Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendEndArray(indent string) Recorder {
	return recorder
}

func (recorder *noopRecorder) AppendNewLine() Recorder {
	return recorder
}
