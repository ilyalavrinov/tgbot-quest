package main

import "github.com/admirallarimda/tgbotbase"
import "github.com/go-redis/redis"

type QuestRecord struct {
	questID string
	quest   Quest
}

type QuestStorage interface {
	StoreQuest(quest Quest) error
	StoreStage(stage Stage) error

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

func (s *redisQuestStorage) StoreQuest(quest Quest) error {
	return nil
}

func (s *redisQuestStorage) StoreStage(stage Stage) error {
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
