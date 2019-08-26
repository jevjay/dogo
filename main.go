package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
	"regexp"
	"strings"
)

// SlackListener stores Slack API data
type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

// ListenAndResponse establishes real time messaging (RTB) connection with Slack messenger
func (s *SlackListener) ListenAndResponse(cfg *CommandsConfig) {
	// Start listening slack events
	rtm := s.client.NewRTM()
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.HandleSlackMessage(ev, rtm, cfg); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// HandleSlackMessage handles Slack message events
func (s *SlackListener) HandleSlackMessage(ev *slack.MessageEvent, rtm *slack.RTM, cfg *CommandsConfig) error {
	text := ev.Text
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	for _, c := range cfg.Command {
		matched, _ := regexp.MatchString(c.Name, text)
		if matched {
			if c.Answer != "" {
				rtm.SendMessage(rtm.NewOutgoingMessage(c.Answer, ev.Channel))
			}

			if c.Image != "" {
				msg, _ := CreateNewContainer(c.Image)
				rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
			}
			break
		}
	}
	return nil
}

func main() {
	// Retrieve values
	token := getenv("SLACKTOKEN")
	config := getenv("CONFIG")
	botID := getenv("BOTID")
	channelID := getenv("CHANNELID")
	// Use default config path if no custom path provided
	if config == "" {
		config = "config.yaml"
	}

	commands := CommandsConfig{}
	commands.ParseConfig(config)

	api := slack.New(token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)))

	var listener = &SlackListener{
		client:    api,
		botID:     botID,
		channelID: channelID,
	}
	listener.ListenAndResponse(&commands)
}
