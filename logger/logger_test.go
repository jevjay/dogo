package logger

import "testing"

func TestLog(t *testing.T) {
	msg := "Test log message"
	isError := false

	err := Log(msg, isError)

	if err != nil {
		t.Errorf("Failed constructing log message")
	}
}
