package quest

import (
	"errors"
	"fmt"
	"github.com/admirallarimda/tgbotbase"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"sync"
)

type AnswerResult struct {
	Active   bool
	Correct  bool
	Finished bool
}

type QuestEngine interface {
	StartQuest(userID tgbotbase.UserID, questID string) error
	CheckAnswer(userID tgbotbase.UserID, answer string) AnswerResult
	GetCurrentQuestion(userID tgbotbase.UserID) tgbotapi.Chattable
	AddQuest(questID string, quest Quest)
}

type activeUserQuest struct {
	quest Quest
	state State
}

type questEngine struct {
	quests map[string]Quest

	activeQuests map[tgbotbase.UserID]activeUserQuest
	mutex        sync.Mutex
}

var _ QuestEngine = &questEngine{}

func NewQuestEngine(pool tgbotbase.RedisPool) QuestEngine {
	engine := &questEngine{
		quests:       make(map[string]Quest, 0),
		activeQuests: make(map[tgbotbase.UserID]activeUserQuest, 0)}
	storage := NewRedisQuestStorage(pool)
	quests, err := storage.LoadAll()
	if err != nil {
		panic(err)
	}

	for _, rec := range quests {
		log.WithFields(log.Fields{"quest": rec.questID, "stages_n": len(rec.quest.stages)}).Info("Quest loaded")
		engine.quests[rec.questID] = rec.quest
	}
	return engine
}

func (q *questEngine) StartQuest(userID tgbotbase.UserID, questID string) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
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
	q.mutex.Lock()
	questData, found := q.activeQuests[userID]
	q.mutex.Unlock()
	if !found {
		log.WithFields(log.Fields{"user": userID}).Warn("Active quest not found on checking answer")
		return AnswerResult{
			Active:   false,
			Correct:  false,
			Finished: false}
	}

	newState := questData.quest.CheckAnswer(answer, questData.state)
	if newState == nil {
		log.WithFields(log.Fields{"user": userID, "answer": answer}).Debug("Incorrect answer")
		return AnswerResult{
			Active:   true,
			Correct:  false,
			Finished: false}
	}

	log.WithFields(log.Fields{"user": userID, "answer": answer}).Debug("Correct answer")
	finished := false
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if newState.IsFinished() {
		finished = true
		delete(q.activeQuests, userID)
	} else {
		q.activeQuests[userID] = activeUserQuest{
			quest: questData.quest,
			state: *newState}
	}
	return AnswerResult{Active: true, Correct: true, Finished: finished}

}

func (q *questEngine) GetCurrentQuestion(userID tgbotbase.UserID) tgbotapi.Chattable {
	q.mutex.Lock()
	questData, found := q.activeQuests[userID]
	q.mutex.Unlock()
	if !found {
		log.WithFields(log.Fields{"user": userID}).Warn("Active quest not found on getting current question")
		return tgbotapi.NewMessage(int64(userID), "You do not have any active quest :(")
	}

	pic := questData.quest.GetPicture(questData.state)
	text := questData.quest.GetQuestion(questData.state)
	if pic != nil {
		buf := tgbotapi.FileBytes{
			Name:  "TODO_put_some_name_here",
			Bytes: pic}
		msg := tgbotapi.NewPhotoUpload(int64(userID), buf)
		msg.Caption = text
		return msg
	}
	return tgbotapi.NewMessage(int64(userID), text)
}

func (q *questEngine) AddQuest(questID string, quest Quest) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.quests[questID] = quest
}
