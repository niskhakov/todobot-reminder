package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/niskhakov/gotodoist"
	"github.com/niskhakov/todobot-reminder/pkg/config"
	"github.com/niskhakov/todobot-reminder/pkg/jobqueue"
	"github.com/niskhakov/todobot-reminder/pkg/repository"
)

const (
	TodoistMainRequesterJobID jobqueue.JobID = "TodoistMainRequesterJobID"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	todoistClient   *todoist.Client
	tokenRepository repository.TokenRepository
	redirectURL     string

	messages config.Messages
	jobQueue jobqueue.Scheduler
}

func NewBot(bot *tgbotapi.BotAPI, todoistClient *todoist.Client, tr repository.TokenRepository, redirectURL string, messages config.Messages, jq jobqueue.Scheduler) *Bot {
	return &Bot{bot: bot, todoistClient: todoistClient, tokenRepository: tr, redirectURL: redirectURL, messages: messages, jobQueue: jq}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.initTodoistListening()

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}

}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

// func (b *Bot) BroadcastMessage(msg string) {
// 	b.tokenRepository.ForEach(repository.AccessTokens, func(chatID int64, accessToken string, accumulator interface{}) error {
// 		tmsg := tgbotapi.NewMessage(chatID, "ReminderGO")
// 		rmsg, err := b.bot.Send(tmsg)
// 		if err != nil {
// 			log.Printf("Error while sending message: %s", err.Error())
// 		}

// 		log.Printf("ChatID: %d, AccessToken: %s, Message: %s\n", chatID, accessToken, rmsg.Text)

// 		return nil
// 	}, nil)
// }
