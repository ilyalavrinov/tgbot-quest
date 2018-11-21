package quest

import "github.com/admirallarimda/tgbotbase"
import "github.com/go-redis/redis"
import "fmt"
import "math"
import "errors"
import "strings"
import log "github.com/sirupsen/logrus"

type QuestRecord struct {
	questID string
	quest   Quest
}

func NewQuestRecord(questID string, quest Quest) *QuestRecord {
	return &QuestRecord{
		questID, quest}
}

type StageRecord struct {
	stageID string
	stage   Stage
}

type QuestStorage interface {
	StoreQuest(quest QuestRecord) error
	StoreStage(questID string, stage StageRecord) error

	LoadAll() ([]QuestRecord, error)
	LoadQuest(questID string) (*Quest, error)
	LoadStage(questID, stageID string) (*Stage, error)
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

	if rec.stage.pic != nil {
		err = s.client.HSet(redisQuestion(stageKey), "pic", rec.stage.pic).Err()
		if err != nil {
			return err
		}
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
	questKeys, err := tgbotbase.GetAllKeys(s.client, scanQuests())
	if err != nil {
		return nil, err
	}
	uniqueQuestKeys := make(map[string]bool, len(questKeys))
	for _, key := range questKeys {
		parts := strings.Split(key, ":")
		uniqueQuestKeys[parts[2]] = true
	}
	quests := make([]QuestRecord, 0, len(uniqueQuestKeys))
	for k := range uniqueQuestKeys {
		quest, err := s.LoadQuest(k)
		if err != nil {
			log.WithFields(log.Fields{"quest": k, "error": err}).Warn("Unable to load quest")
			continue
		}
		quests = append(quests, QuestRecord{k, *quest})
	}
	return quests, nil
}

func (s *redisQuestStorage) LoadQuest(questID string) (*Quest, error) {
	stageKeys, err := tgbotbase.GetAllKeys(s.client, scanStages(questID))
	if err != nil {
		return nil, err
	}
	stages := make(map[string]Stage, len(stageKeys))
	for _, key := range stageKeys {
		parts := strings.Split(key, ":")
		stageID := parts[3]
		stage, err := s.LoadStage(questID, stageID)
		if err != nil {
			log.WithFields(log.Fields{"quest": questID, "stage": stageID, "error": err}).Warn("Unable to load stage")
			return nil, err
		}
		stages[stageID] = *stage
	}
	return &Quest{stages}, nil
}

func (s *redisQuestStorage) LoadStage(questID, stageID string) (*Stage, error) {
	stageKey := redisStageKey(questID, stageID)
	fields, err := s.client.HGetAll(redisQuestion(stageKey)).Result()
	if err != nil {
		return nil, err
	}

	qtext, found := fields["text"]
	if !found {
		return nil, errors.New("Empty question text")
	}

	answers, err := s.client.LRange(redisAnswers(stageKey), 0, math.MaxInt64).Result()
	if err != nil {
		return nil, err
	}
	if len(answers) == 0 {
		return nil, errors.New("Empty list of answers")
	}

	stage := NewStage(qtext, answers)
	if pic, found := fields["pic"]; found {
		stage.AddPicture([]byte(pic))
	}

	return &stage, nil
}

func redisQuestKey(questID string) string {
	return fmt.Sprintf("tg:quest:%s", questID)
}

func scanQuests() string {
	return "tg:quest:*"
}

func scanStages(questID string) string {
	return fmt.Sprintf("tg:quest:%s:*:question", questID)
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
