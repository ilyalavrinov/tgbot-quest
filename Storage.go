package main

import "github.com/admirallarimda/tgbotbase"
import "github.com/go-redis/redis"
import "fmt"

type QuestRecord struct {
	questID string
	quest   Quest
}

type StageRecord struct {
	stageID string
	stage   Stage
}

type QuestStorage interface {
	StoreQuest(quest QuestRecord) error
	StoreStage(questID string, stage StageRecord) error

	LoadAll() ([]QuestRecord, error)
	LoadQuest(questID string) (Quest, error)
	LoadStage(questID, stageID string) (Stage, error)
}

type redisQuestStorage struct {
	client *redis.Client
}

func NewRedisQuestStorage(pool tgbotbase.RedisPool) QuestStorage {
	return &redisQuestStorage{client: pool.GetConnByName("quest")}
}

func (s *redisQuestStorage) StoreQuest(q QuestRecord) error {
	var err error
	for stageID, stage := range q.quest.stages {
		err = s.StoreStage(q.questID, StageRecord{stageID, stage})
		if err != nil {
			break
		}
	}
	return err
}

func (s *redisQuestStorage) StoreStage(questID string, rec StageRecord) error {
	stageKey := redisStageKey(questID, rec.stageID)
	err := s.client.HSet(redisQuestion(stageKey), "text", rec.stage.question).Err()
	if err != nil {
		return err
	}
	answers := make([]string, 0, len(rec.stage.answers))
	for a, _ := range rec.stage.answers {
		answers = append(answers, a)
	}
	err = s.client.LPush(redisAnswers(stageKey), answers).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *redisQuestStorage) LoadAll() ([]QuestRecord, error) {
	return nil, nil
}

func (s *redisQuestStorage) LoadQuest(questID string) (Quest, error) {
	return Quest{}, nil
}

func (s *redisQuestStorage) LoadStage(questID, stageID string) (Stage, error) {
	return Stage{}, nil
}

func redisQuestKey(questID string) string {
	return fmt.Sprintf("tg:quest:%s", questID)
}

func redisStageKey(questID string, stageID string) string {
	return fmt.Sprintf("%s:%s", redisQuestKey(questID), stageID)
}

func redisQuestion(stageKey string) string {
	return fmt.Sprintf("%s:question", stageKey)
}

func redisAnswers(stageKey string) string {
	return fmt.Sprintf("%s:answers", stageKey)
}
