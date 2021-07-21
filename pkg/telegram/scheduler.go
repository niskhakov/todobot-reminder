package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/niskhakov/gotodoist"
	"github.com/niskhakov/todobot-reminder/pkg/jobqueue"
	"github.com/niskhakov/todobot-reminder/pkg/repository"
)

func (b *Bot) initTodoistListening() {
	var every time.Duration = 30 * time.Second
	job := b.createTodoistListenerJob()
	b.jobQueue.Add(context.Background(), job, every)
}

func (b *Bot) createTodoistListenerJob() jobqueue.Job {
	return jobqueue.Job{
		F: func(ctx context.Context) {
			fnc := b.forEveryEntryInDBCallback(ctx)
			b.tokenRepository.ForEach(repository.AccessTokens, fnc, nil)
		},
		ID: TodoistMainRequesterJobID,
	}
}

func (b *Bot) forEveryEntryInDBCallback(ctx context.Context) func(chatID int64, accessToken string, accumulator interface{}) error {
	return func(chatID int64, accessToken string, accumulator interface{}) error {
		projects, err := b.todoistClient.GetProjects(ctx, accessToken)
		if err != nil {
			return err
		}

		inboxProject, err := findInboxProject(projects)
		if err != nil {
			return err
		}

		tasks, err := b.todoistClient.GetTasksByProject(ctx, accessToken, inboxProject.ID)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if t.Due.Datetime == "" {
				continue
			}

			ddt, err := time.Parse(time.RFC3339, t.Due.Datetime)
			if err != nil {
				return errors.New("can't case to time ")
			}

			timeDiff := time.Until(ddt)
			if timeDiff > 0 && timeDiff < 5*time.Minute {
				log.Printf("Found on project %s - task %s\n", inboxProject.Name, t.Content)
				jID := jobqueue.JobID(fmt.Sprint(t.ID))
				b.jobQueue.Stop(jID)
				msgString := fmt.Sprintf("%s - %s", t.Content, timeDiff.String())
				taskjob := jobqueue.Job{
					F: func(ctx context.Context) {
						msg := tgbotapi.NewMessage(chatID, msgString)
						b.bot.Send(msg)
					},
					ID: jID,
				}

				b.jobQueue.AddAt(ctx, taskjob, ddt)
			}
		}
		return nil
	}
}

func findInboxProject(projects []todoist.Project) (*todoist.Project, error) {
	// Find inbox project
	inboxID := -1
	var inboxProject todoist.Project

	for _, p := range projects {
		if p.InboxProject {
			inboxID = p.ID
			inboxProject = p
		}
	}

	if inboxID == -1 {
		return nil, errors.New("inbox project was not found")
	}

	return &inboxProject, nil
}
