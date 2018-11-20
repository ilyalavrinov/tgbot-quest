package main

import (
	"errors"
	"fmt"

	"github.com/admirallarimda/tgbotbase"
	log "github.com/sirupsen/logrus"
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
	questData, found := q.activeQuests[userID]
	if !found {
		log.WithFields(log.Fields{"user": userID}).Warn("Active quest not found on checking answer")
		return AnswerResult{
			correct:  false,
			finished: false}
	}

	newState := questData.quest.CheckAnswer(answer, questData.state)
	if newState == nil {
		log.WithFields(log.Fields{"user": userID, "answer": answer}).Debug("Incorrect answer")
		return AnswerResult{
			correct:  false,
			finished: false}
	}

	log.WithFields(log.Fields{"user": userID, "answer": answer}).Debug("Correct answer")
	correct := true
	finished := false
	if newState.IsFinished() {
		finished = true
		delete(q.activeQuests, userID)
	} else {
		q.activeQuests[userID] = activeUserQuest{
			quest: questData.quest,
			state: *newState}
	}
	return AnswerResult{correct: correct, finished: finished}

}

func (q *questEngine) GetCurrentQuestion(userID tgbotbase.UserID) tgbotapi.Chattable {
	questData, found := q.activeQuests[userID]
	if !found {
		log.WithFields(log.Fields{"user": userID}).Warn("Active quest not found on getting current question")
		return tgbotapi.NewMessage(int64(userID), "You do not have any active quest :(")
	}
	return tgbotapi.NewMessage(int64(userID), questData.quest.GetQuestion(questData.state))
}
