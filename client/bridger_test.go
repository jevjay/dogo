// Copyright 2021 tappythumbz development
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"reflect"
	"testing"

	"github.com/nlopes/slack"
)

func TestFormatMessageStdout(t *testing.T) {
	msg := "Test message"
	fmt := "stdout"
	result := SlackMessage{Text: msg}

	out, _ := FormatMessage(msg, fmt)

	if !reflect.DeepEqual(out, result) {
		t.Errorf("Input message of %v type does not match expected value after parsing.", fmt)
	}
}

func TestFormatMessageJson(t *testing.T) {
	msg := "{\"time\":\"2021-02-20T18:25:46-04:00\",\"level\":\"info\",\"message\":\"Test message\"}"
	fmt := "json"
	result := SlackMessage{
		Text:       "Test message",
		Attachment: slack.Attachment{},
	}

	out, err := FormatMessage(msg, fmt)

	if !reflect.DeepEqual(out, result) {
		t.Errorf("Input message of %v type does not match expected value after parsing", fmt)
	}

	if err != nil {
		t.Errorf("Failed parsing %v format message", fmt)
	}
}
