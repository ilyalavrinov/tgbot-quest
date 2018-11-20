package main

import (
	"errors"
	"fmt"

	"github.com/admirallarimda/tgbotbase"
	"gopkg.in/telegram-bot-api.v4"
)

type AnswerResult struct {
	correct  bool
	finished bool
}

type QuestEngine interface {
	StartQuest(userID tgbotbase.UserID, questID string) error
	CheckAnswer(userID tgbotbase.UserID, answer string) AnswerResult
	GetCurrentQuestion(userID tgbotbase.UserID) tgbotapi.Chattable
}

type activeUserQuest struct {
	quest Quest
	state State
}

type questEngine struct {
	quests map[string]Quest

	activeQuests map[tgbotbase.UserID]activeUserQuest
}

var _ QuestEngine = &questEngine{}

func (q *questEngine) StartQuest(userID tgbotbase.UserID, questID string) error {
	quest, found := q.quests[questID]
	if !found {
		return errors.New(fmt.Sprintf("Quest '%s' is not registered", questID))
	}

	q.activeQuests[userID] = activeUserQuest{
		quest: quest,
		state: quest.CreateInitialState()}

	return nil
}

func (q *questEngine) CheckAnswer(userID tgbotbase.UserID, answer string) AnswerResult {
	return AnswerResult{}
}

func (q *questEngine) GetCurrentQuestion(userID tgbotbase.UserID) tgbotapi.Chattable {
	return tgbotapi.NewMessage(int64(userID), "hello")
}
