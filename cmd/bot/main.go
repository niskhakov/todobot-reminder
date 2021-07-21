package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/niskhakov/gotodoist"
	"github.com/niskhakov/todobot-reminder/pkg/config"
	"github.com/niskhakov/todobot-reminder/pkg/jobqueue/simplejq"
	"github.com/niskhakov/todobot-reminder/pkg/repository"
	"github.com/niskhakov/todobot-reminder/pkg/repository/boltdb"
	"github.com/niskhakov/todobot-reminder/pkg/server"
	"github.com/niskhakov/todobot-reminder/pkg/telegram"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(fmt.Errorf("can't init configuration: %w", err))
	}

	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	todoistClient, err := todoist.NewClient(cfg.TodoistClientID, cfg.TodoistClientSecret)
	if err != nil {
		log.Fatal(fmt.Errorf("can't create new todoist client: %w", err))
	}

	jq := simplejq.NewScheduler()

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("can't init database: %w", err))
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, todoistClient, tokenRepository, cfg.AuthServerURL, cfg.Messages, jq)

	authorizationServer := server.NewAuthorizationServer(todoistClient, tokenRepository, cfg.TelegramBotURL)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(fmt.Errorf("can't start telegram bot: %w", err))
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(fmt.Errorf("can't start auth server: %w", err))
	}

}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		// TODO: Выпилить request tokens
		_, err = tx.CreateBucketIfNotExists([]byte(repository.CodeTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
