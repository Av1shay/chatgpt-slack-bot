package main

import (
	"context"
	"fmt"
	"github.com/Av1shay/chatgpt-slack-bot/gpt"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	defaultMaxTokens = 100
	defaultDebug     = true
)

var slackMentionReg *regexp.Regexp

func init() {
	slackMentionReg = regexp.MustCompile("(?:\\s)<@[^, ]*|(?:^)<@[^, ]*")
}

type ChatEngine interface {
	Do(ctx context.Context, prompt string) (string, error)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("failed to load .env file:", err)
	}
	openAiToken, found := os.LookupEnv("OPENAI_API_KEY")
	if !found {
		log.Fatal("OPENAI_API_KEY env is missing from config")
	}
	slackAuthToken, found := os.LookupEnv("SLACK_AUTH_TOKEN")
	if !found {
		log.Fatal("SLACK_AUTH_TOKEN env is missing from config")
	}
	slackAppToken, found := os.LookupEnv("SLACK_APP_TOKEN")
	if !found {
		log.Fatal("SLACK_APP_TOKEN env is missing from config")
	}
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = defaultDebug
	}
	maxTokens, err := strconv.Atoi(os.Getenv("CHAT_GPT_MAX_TOKENS"))
	if err != nil {
		maxTokens = defaultMaxTokens
	}

	gptService := gpt.NewService(openAiToken, gpt.WithMaxTokens(maxTokens))

	slackClient := slack.New(slackAuthToken, slack.OptionDebug(debug), slack.OptionAppLevelToken(slackAppToken))
	socketClient := socketmode.New(slackClient,
		socketmode.OptionDebug(debug),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)))

	go consumeMessages(ctx, slackClient, socketClient, gptService)

	panic(socketClient.Run())
}

func consumeMessages(ctx context.Context, slackClient *slack.Client, socketClient *socketmode.Client, chatEngine ChatEngine) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down socketmode listener")
			return
		case event := <-socketClient.Events:
			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				socketClient.Ack(*event.Request)
				eventPayload, _ := event.Data.(slackevents.EventsAPIEvent)
				if err := handleEventMessage(ctx, slackClient, chatEngine, eventPayload); err != nil {
					log.Println("HandleEventMessage error", err)
				}
			}
		}
	}
}

func handleEventMessage(ctx context.Context, slackClient *slack.Client, chatEngine ChatEngine, eventPayload slackevents.EventsAPIEvent) error {
	switch eventPayload.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventPayload.InnerEvent
		switch event := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			prompt := slackMentionReg.ReplaceAllString(event.Text, "")
			resp, err := chatEngine.Do(ctx, prompt)
			if err != nil {
				return fmt.Errorf("failed to get message from chat engine: %w", err)
			}
			attachment := slack.Attachment{}
			attachment.Text = resp
			attachment.Color = "#4af030"
			_, _, err = slackClient.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
			if err != nil {
				return fmt.Errorf("failed to post message: %w", err)
			}
		default:
			return fmt.Errorf("unsupported inner event type %T", event)
		}
	default:
		return fmt.Errorf("unsupported event type %q", eventPayload.Type)
	}
	return nil
}
