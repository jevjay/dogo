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
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	a "dogo/agent"
	cfg "dogo/config"
	"dogo/logger"

	"github.com/nlopes/slack"
)

// SlackListener stores Slack API data
type SlackListener struct {
	// Slack client
	Client *slack.Client
	// Slack bot ID
	BotID string
	// Slack channel ID
	ChannelID string
}

// SlackMessage stores data used in constructing a message event
type SlackMessage struct {
	// Slack channel name
	Channel string
	// Message text
	Text string
	// Message timestamp
	Timestamp string
	// Threaded message timestamp
	ThreadTimestamp string
	// Slack message attachment
	Attachment slack.Attachment
}

// AgentOut stores data received from agents with JSON output format
type AgentOut struct {
	Time       string           `json:"time"`
	Level      string           `json:"level"`
	Message    string           `json:"message"`
	Attachment slack.Attachment `json: attachment`
}

// FormatMessage prepares final Slack message format from agents Stdout/JSON output message
func FormatMessage(msg string, format string) (SlackMessage, error) {
	switch format {
	case "stdout":
		// Resutrn stdout format directly as it is
		return SlackMessage{
			Text: msg,
		}, nil

	case "json":
		var out AgentOut
		if err := json.Unmarshal([]byte(msg), &out); err != nil {
			return SlackMessage{}, err
		}

		return SlackMessage{
			Text:       out.Message,
			Attachment: out.Attachment,
		}, nil
	}

	err := fmt.Errorf("Failed to parse agents configured output format")
	return SlackMessage{}, err
}

// SendMessage sends message to SlackChannel
func SendMessage(api *slack.Client, msg SlackMessage) error {
	params := slack.NewPostMessageParameters()

	// Respond in thread if not a direct message.
	if !strings.HasPrefix(msg.Channel, "D") {
		params.ThreadTimestamp = msg.Timestamp
	}

	// Respond in same thread if message came from a thread.
	if msg.ThreadTimestamp != "" {
		params.ThreadTimestamp = msg.ThreadTimestamp
	}

	_, _, err := api.PostMessage(
		msg.Channel,
		slack.MsgOptionPostMessageParameters(params),
		slack.MsgOptionText(msg.Text, false),
		slack.MsgOptionAttachments(msg.Attachment),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	if err != nil {
		return err
	}

	return nil
}

// ListenAndResponse establishes real time messaging (RTB) connection with Slack messenger
func (s *SlackListener) ListenAndResponse(api *slack.Client, cmds *cfg.Commands, c chan cfg.OutputConfig) {
	// Start listening slack events
	rtm := s.Client.NewRTM()
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.HandleSlackMessage(api, ev, cmds, c); err != nil {
				logger.Log(err.Error(), true)
			}
		}
	}
}

// HandleSlackMessage handles Slack message events
func (s *SlackListener) HandleSlackMessage(api *slack.Client, ev *slack.MessageEvent, cmds *cfg.Commands, c chan cfg.OutputConfig) error {
	text := ev.Text
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	c1 := make(chan *a.Container, 1)
	cmd := cmds.Configs

	for _, in := range cmd.Input {
		matched, _ := regexp.MatchString(in.Command, text)
		if matched {
			// Check command type
			switch in.Type {
			case "docker":
				// Run your long running function in it's own goroutine and pass back it's
				// response into our channel.
				go func() {
					// Call docker agent locally
					resp, err := a.Execute(in.Command, in.Image)
					c1 <- resp
					if err != nil {
						logger.Log(err.Error(), true)
					}
				}()

				// Listen on our channel AND a timeout channel - which ever happens first.
				select {
				case res := <-c1:
					msg, err := FormatMessage(res.Out, in.Format)
					if err != nil {
						logger.Log(err.Error(), true)
					}

					msg.Channel = ev.Channel
					msg.Timestamp = ev.EventTimestamp
					msg.ThreadTimestamp = ev.ThreadTimestamp

					err = SendMessage(api, msg)

					if err != nil {
						return err
					}

				case <-time.After(in.Timeout * time.Second):
					err := fmt.Errorf("Command %v timed out after %v", in.Command, in.Timeout)
					logger.Log(err.Error(), true)

					// Inform user about timedout command
					slackErr := SendMessage(api, SlackMessage{
						Channel:         ev.Channel,
						Text:            err.Error(),
						Timestamp:       ev.EventTimestamp,
						ThreadTimestamp: ev.ThreadTimestamp,
					})

					if slackErr != nil {
						return slackErr
					}
				}

			case "lambda":
				// Call Lambda execution and wait for response

			case "kubernetes":
				// Call Kubernetes job execution and wait for response

			default:
				// Construct error message
				err := fmt.Errorf("Unknown input agent value %s", in.Type)
				// Show error message
				return err
			}
		}
	}
	return nil
}
