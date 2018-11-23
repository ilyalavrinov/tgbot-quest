package main

import (
	"github.com/admirallarimda/tgbot-quest/internal/pkg/quest"
	"github.com/admirallarimda/tgbotbase"
	"gopkg.in/telegram-bot-api.v4"
)

type statHandler struct {
	tgbotbase.BaseHandler
	resmon quest.ResultMonitor
}

func (h *statHandler) Name() string {
	return "stat handler"
}

func (h *statHandler) HandleOne(msg tgbotapi.Message) {
	h.resmon.SendStats(msg.CommandArguments())
}

func (h *statHandler) Init(outCh chan<- tgbotapi.Chattable, srvCh chan<- tgbotbase.ServiceMsg) tgbotbase.HandlerTrigger {
	h.OutMsgCh = outCh
	return tgbotbase.NewHandlerTrigger(nil, []string{"stats"})
}

func newStatsHandler(monitor quest.ResultMonitor) tgbotbase.IncomingMessageHandler {
	return &statHandler{resmon: monitor}
}
