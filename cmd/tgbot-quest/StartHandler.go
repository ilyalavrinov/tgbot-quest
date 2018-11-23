package main

import (
	"fmt"

	"github.com/admirallarimda/tgbot-quest/internal/pkg/quest"
	"github.com/admirallarimda/tgbotbase"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"sync"
)

type startHandler struct {
	tgbotbase.BaseHandler
	engine quest.QuestEngine
	names  *sync.Map
}

func (h *startHandler) Name() string {
	return "start"
}

func (h *startHandler) HandleOne(msg tgbotapi.Message) {
	userID := tgbotbase.UserID(msg.From.ID)
	chatID := msg.Chat.ID
	log.WithFields(log.Fields{"userID": userID, "userName": msg.From.UserName, "message": msg.Text}).Debug("Incoming start")
	h.names.Store(msg.From.ID, msg.From.UserName)
	questID := msg.CommandArguments()
	err := h.engine.StartQuest(userID, questID)
	if err != nil {
		h.OutMsgCh <- tgbotapi.NewMessage(chatID, fmt.Sprintf("Я не смог стартовать квест с именем '%s'", questID))
	} else {
		h.OutMsgCh <- h.engine.GetCurrentQuestion(userID)
	}
}

func (h *startHandler) Init(outCh chan<- tgbotapi.Chattable, srvCh chan<- tgbotbase.ServiceMsg) tgbotbase.HandlerTrigger {
	h.OutMsgCh = outCh
	return tgbotbase.NewHandlerTrigger(nil, []string{"start"})
}

func NewStartHandler(engine quest.QuestEngine, usernames *sync.Map) tgbotbase.IncomingMessageHandler {
	return &startHandler{engine: engine,
		names: usernames}
}
