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

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

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
	Args    string        `json:"args" yaml:"args"`
	Command string        `json:"command" yaml:"command"`
	Desc    string        `json:"desc" yaml:"desc"`
	Error   string        `json:"error" yaml:"error"`
	Image   string        `json:"image" yaml:"image"`
	Format  string        `json:"format" yaml:"format"`
	Type    string        `json:"type" yaml:"type"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	Output  OutputConfig  `json:"output" yaml:"output"`
}

// OutputConfig stores configuration of a command output
type OutputConfig struct {
	Message    string `json:"message" yaml:"message"`
	Attachment string `json:"attachment" yaml:"attachment"`
	Color      string `json:"color" yaml:"color"`
}

// GetEnv retrieves the string of the environment variable
func GetEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		fmt.Println("environment variable is not set: " + name)
	}
	return v
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
