package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/niskhakov/todobot-reminder/pkg/repository"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthLink(message.Chat.ID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.messages.Responses.Start, authLink))
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessTokens)
}

func (b *Bot) generateAuthLink(chatID int64) (string, error) {
	authLink := b.generateTodoistLink(chatID)

	return authLink, nil
}

func (b *Bot) generateTodoistLink(chatID int64) string {
	return b.todoistClient.GetAuthorizationRequestURL(context.Background(), fmt.Sprint(chatID))
}
