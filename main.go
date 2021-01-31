package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	a "dogo/agents"

	"github.com/nlopes/slack"
)

// SlackListener stores Slack API data
type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

// ListenAndResponse establishes real time messaging (RTB) connection with Slack messenger
func (s *SlackListener) ListenAndResponse(cmds *Commands, c chan OutputConfig) {
	// Start listening slack events
	rtm := s.client.NewRTM()
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.HandleSlackMessage(ev, rtm, cmds, c); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// HandleSlackMessage handles Slack message events
func (s *SlackListener) HandleSlackMessage(ev *slack.MessageEvent, rtm *slack.RTM, cmds *Commands, c chan OutputConfig) error {
	text := ev.Text
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	cmd := cmds.Configs
	for _, in := range cmd.Input {
		matched, _ := regexp.MatchString(in.Command, text)
		if matched {
			// Check command type
			switch in.Type {
			case "Docker":
				// Call docker agent locally
				a.CreateNewContainer(in.Command, in.Image)
			case "Lambda":
				// Call Lambda execution and wait for response

			case "Kubernetes":
				// Call Kubernetes job execution and wait for response

			default:
				// Show error message

			}
		}
	}
	return nil
}

func main() {
	// Retrieve values
	token := getenv("DOGO_SLACKTOKEN")
	config := getenv("DOGO_CONFIG")
	botID := getenv("DOGO_BOTID")
	channelID := getenv("DOGO_CHANNELID")
	debugMode := getenv("DOGO_DEBUG")
	debug, _ := strconv.ParseBool(debugMode)

	// Use default config path if no custom path provided
	if config == "" {
		config = "config.yaml"
	}

	commands := Commands{}
	commands.ParseConfig(config, debug)

	api := slack.New(token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)))

	var listener = &SlackListener{
		client:    api,
		botID:     botID,
		channelID: channelID,
	}
	// Set up agent output channel
	c := make(chan OutputConfig)

	listener.ListenAndResponse(&commands, c)
}
