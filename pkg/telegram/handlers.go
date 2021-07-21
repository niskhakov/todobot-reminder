package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/niskhakov/todobot-reminder/pkg/repository"
)

const (
	commandStart   = "start"
	commandChatIDs = "debug_chats"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandChatIDs:
		return b.handleChatIDsCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}

}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	log.Printf("[%s] %s", message.From.UserName, message.Text)

	// msg := tgbotapi.NewMessage(message.Chat.ID, message.Text+" "+fmt.Sprint(message.Chat.ID))
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.NoMessage)

	// accessToken, err := b.getAccessToken(message.Chat.ID)
	// if err != nil {
	// 	return errUnauthorized
	// }

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		// User is not authorized, he has no access token
		return b.initAuthorizationProcess(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.AlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleChatIDsCommand(message *tgbotapi.Message) error {
	chatIDs := make([]string, 0)
	chatIDs = append(chatIDs, "Registered Chat IDs:")

	var fnc repository.IterateFunc = func(chatID int64, accessToken string, accumulator interface{}) error {
		chatIDs = append(chatIDs, fmt.Sprint(chatID))

		return nil
	}
	b.tokenRepository.ForEach(repository.AccessTokens, fnc, nil)

	chatIDsStr := strings.Join(chatIDs, "\n")

	msg := tgbotapi.NewMessage(message.Chat.ID, chatIDsStr)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}
