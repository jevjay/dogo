package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Commands contains a list of commands configs
type Commands struct {
	Configs CommandConfig `json:"config" yaml:"config"`
}

// CommandConfig stores command input/output configurations
type CommandConfig struct {
	Input []InputConfig `json:"input" yaml:"input"`
}

// InputConfig stores configuration of a command input
type InputConfig struct {
	Args    string       `json:"args" yaml:"args"`
	Command string       `json:"command" yaml:"command"`
	Desc    string       `json:"desc" yaml:"desc"`
	Error   string       `json:"error" yaml:"error"`
	Image   string       `json:"image" yaml:"image"`
	Type    string       `json:"type" yaml:"type"`
	Timeout string       `json:"timeout" yaml:"timeout"`
	Output  OutputConfig `json:"output" yaml:"output"`
}

// OutputConfig stores configuration of a command output
type OutputConfig struct {
	Message    string `json:"message" yaml:"message"`
	Attachment string `json:"attachment" yaml:"attachment"`
	Color      string `json:"color" yaml:"color"`
}

// ParseConfig parses configuration files command(s)
func (c *Commands) ParseConfig(path string, debug bool) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse command config file: %s", err)
		return err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse command config file: %s", err)
		return err
	}
	// Debug output
	if debug {
		log.Println("Commands config: \n", string(data))
	}

	return nil
}
