package quest

import (
	"fmt"
	"github.com/admirallarimda/tgbotbase"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"time"
)

type ResultMonitor interface {
	QuestStarted(questID string, userID tgbotbase.UserID, t time.Time)
	QuestFinished(questID string, userID tgbotbase.UserID, t time.Time)
	QuestionAnsweredCorrectly(questID string, userID tgbotbase.UserID, t time.Time)
	QuestionAnsweredIncorrectly(questID string, userID tgbotbase.UserID, t time.Time)
}

type questEvent struct {
	questID string
	userID  tgbotbase.UserID
	t       time.Time
}

type questStats struct {
	started          time.Time
	finished         time.Time
	answeredTimes    []time.Time
	incorrectAnswers int
}

type tgOwnerNotifyResultMonitor struct {
	startedCh           chan questEvent
	finishedCh          chan questEvent
	answeredCorrectCh   chan questEvent
	answeredIncorrectCh chan questEvent

	stats map[string]map[tgbotbase.UserID]questStats // questID -> userID -> stats

	tgbot  *tgbotbase.Bot
	owners []tgbotbase.UserID
}

func NewTGResultMonitor(tgbot *tgbotbase.Bot, owners []tgbotbase.UserID) ResultMonitor {
	if len(owners) == 0 {
		log.Panic("0 owners")
	}

	mon := &tgOwnerNotifyResultMonitor{
		startedCh:           make(chan questEvent, 0),
		finishedCh:          make(chan questEvent, 0),
		answeredCorrectCh:   make(chan questEvent, 0),
		answeredIncorrectCh: make(chan questEvent, 0),
		stats:               make(map[string]map[tgbotbase.UserID]questStats, 0),
		tgbot:               tgbot,
		owners:              owners}

	go mon.run()
	return mon
}

func (mon *tgOwnerNotifyResultMonitor) QuestStarted(questID string, userID tgbotbase.UserID, t time.Time) {
	mon.startedCh <- questEvent{questID, userID, t}
}

func (mon *tgOwnerNotifyResultMonitor) QuestFinished(questID string, userID tgbotbase.UserID, t time.Time) {
	mon.finishedCh <- questEvent{questID, userID, t}
}

func (mon *tgOwnerNotifyResultMonitor) QuestionAnsweredCorrectly(questID string, userID tgbotbase.UserID, t time.Time) {
	mon.answeredCorrectCh <- questEvent{questID, userID, t}
}

func (mon *tgOwnerNotifyResultMonitor) QuestionAnsweredIncorrectly(questID string, userID tgbotbase.UserID, t time.Time) {
	mon.answeredIncorrectCh <- questEvent{questID, userID, t}
}

func (mon *tgOwnerNotifyResultMonitor) run() {
	for {
		select {
		case e := <-mon.startedCh:
			mon.ensureStats(e.questID)
			stats := mon.stats[e.questID][e.userID]
			stats.started = e.t
			mon.stats[e.questID][e.userID] = stats
			log.WithFields(log.Fields{"quest": e.questID, "user": e.userID, "time": e.t}).Debug("User started a quest")
			mon.send(fmt.Sprintf("Started '%s' by ID %d at %s", e.questID, e.userID, e.t))
		case e := <-mon.finishedCh:
			mon.ensureStats(e.questID)
			stats := mon.stats[e.questID][e.userID]
			stats.finished = e.t
			mon.stats[e.questID][e.userID] = stats
			tdiff := stats.finished.Sub(stats.started)
			log.WithFields(log.Fields{"quest": e.questID, "user": e.userID, "time": e.t, "tdiff": tdiff}).Debug("User finished a quest")
			mon.send(fmt.Sprintf("Finished '%s' by ID %d at %s (spent %s, made %d mistakes)", e.questID, e.userID, e.t, tdiff, stats.incorrectAnswers))
		case e := <-mon.answeredCorrectCh:
			mon.ensureStats(e.questID)
			stats := mon.stats[e.questID][e.userID]
			stats.answeredTimes = append(stats.answeredTimes, e.t)
			mon.stats[e.questID][e.userID] = stats
			log.WithFields(log.Fields{"quest": e.questID, "user": e.userID, "time": e.t, "answerN": len(stats.answeredTimes)}).Debug("User answered correctly")
		case e := <-mon.answeredIncorrectCh:
			mon.ensureStats(e.questID)
			stats := mon.stats[e.questID][e.userID]
			stats.incorrectAnswers++
			mon.stats[e.questID][e.userID] = stats
			log.WithFields(log.Fields{"quest": e.questID, "user": e.userID, "time": e.t, "total_incorrect": stats.incorrectAnswers}).Debug("User answered incorrectly")
		}
	}
}

func (mon *tgOwnerNotifyResultMonitor) send(msg string) {
	for _, owner := range mon.owners {
		mon.tgbot.Send(tgbotapi.NewMessage(int64(owner), msg))
	}
}

func (mon *tgOwnerNotifyResultMonitor) ensureStats(questID string) {
	if _, found := mon.stats[questID]; !found {
		mon.stats[questID] = make(map[tgbotbase.UserID]questStats, 0)
	}
}
