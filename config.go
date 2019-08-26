package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// CommandsConfig stores command data list
type CommandsConfig struct {
	Command []Command `json:"command" yaml:"command"`
}

// Command stores configuration of a command retrieved from a config file
type Command struct {
	Answer string `json:"answer" yaml:"answer"`
	Name   string `json:"name" yaml:"name"`
	Desc   string `json:"desc" yaml:"desc"`
	Image  string `json:"image" yaml:"image"`
}

// ParseConfig parses configuration files command(s)
func (c *CommandsConfig) ParseConfig(path string) error {
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
	return nil
}
