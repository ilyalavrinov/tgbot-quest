package main

import (
	"regexp"

	"github.com/admirallarimda/tgbotbase"
	"gopkg.in/telegram-bot-api.v4"
)

type answerHandler struct {
	tgbotbase.BaseHandler
}

func (h *answerHandler) Name() string {
	return "answer handler"
}

func (h *answerHandler) HandleOne(msg tgbotapi.Message) {

}

func (h *answerHandler) Init(outCh chan<- tgbotapi.Chattable, srvCh chan<- tgbotbase.ServiceMsg) tgbotbase.HandlerTrigger {
	h.OutMsgCh = outCh
	return tgbotbase.NewHandlerTrigger(regexp.MustCompile("*"), nil)
}

func NewAnswerHandler() tgbotbase.IncomingMessageHandler {
	return &answerHandler{}
}
