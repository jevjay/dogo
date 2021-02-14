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

package main

import (
	cl "dogo/client"
	cfg "dogo/config"

	"log"
	"os"
	"strconv"

	"github.com/nlopes/slack"
)

func main() {
	// Retrieve values
	token := cfg.GetEnv("DOGO_SLACKTOKEN")
	config := cfg.GetEnv("DOGO_CONFIG")
	botID := cfg.GetEnv("DOGO_BOTID")
	channelID := cfg.GetEnv("DOGO_CHANNELID")
	debugMode := cfg.GetEnv("DOGO_DEBUG")
	// Convert debug mode value from
	debug, _ := strconv.ParseBool(debugMode)

	// Use default config path if no custom path provided
	if config == "" {
		config = "config.yaml"
	}

	commands := cfg.Commands{}
	commands.ParseConfig(config, debug)

	api := slack.New(token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)))

	listener := cl.SlackListener{
		Client:    api,
		BotID:     botID,
		ChannelID: channelID,
	}

	// Set up agent output channel
	c := make(chan cfg.OutputConfig)

	listener.ListenAndResponse(api, &commands, c)
}
