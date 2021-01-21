package main

import "testing"

func TestParseConfig(t *testing.T) {
	// Load
	commands := Commands{}
	commands.ParseConfig("./test/config.yml", false)

	for _, in := range commands.Configs.Input {
		if in.Command != "joke" {
			t.Errorf("Failed parsing commands config file. Command value expected to be %v, got %v", "joke", in.Command)
		}

		if in.Desc != "Tell a random joke from the list" {
			t.Errorf("Failed parsing commands config file. Command value expected to be %v, got %v", "Tell a random joke from the list", in.Desc)
		}

		if in.Timeout != "2" {
			t.Errorf("Failed parsing commands config file. Command value expected to be %v, got %v", "2", in.Timeout)
		}
	}
}
