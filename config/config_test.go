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

package config

import "testing"

func TestParseConfig(t *testing.T) {
	// Load
	commands := Commands{}
	commands.ParseConfig("../test/config.yml", false)

	for _, in := range commands.Configs.Input {
		if in.Command != "hello-stdout" && in.Command != "hello-json" {
			t.Errorf("Failed parsing commands config file. Command value expected to be %v, got %v", "joke", in.Command)
		}

		if in.Desc != "Say hello" {
			t.Errorf("Failed parsing commands config file. Command value expected to be %v, got %v", "Tell a random joke from the list", in.Desc)
		}

		if in.Timeout != 10 && in.Timeout != 3 {
			t.Errorf("Failed parsing commands config file. Command value expected to be %v, got %v", "2", in.Timeout)
		}
	}
}
