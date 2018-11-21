package main

import log "github.com/sirupsen/logrus"
import "gopkg.in/gcfg.v1"
import "github.com/admirallarimda/tgbotbase"

type config struct {
	tgbotbase.Config
	Redis tgbotbase.RedisConfig
}

func readGcfg(filename string) config {
	log.WithFields(log.Fields{"file": filename}).Info("Reading configuration")

	var cfg config

	err := gcfg.ReadFileInto(&cfg, filename)
	if err != nil {
		log.WithFields(log.Fields{"file": filename, "error": err}).Error("Could not correctly parse configuration")
		panic(err)
	}

	log.WithFields(log.Fields{"file": filename, "cfg": cfg}).Info("Configuration has been successfully read")
	return cfg
}

func main() {
	log.Info("Starting daily budget bot")

	cfg := readGcfg("bot.cfg")
	botCfg := tgbotbase.Config{TGBot: cfg.TGBot, Proxy_SOCKS5: cfg.Proxy_SOCKS5}
	tgbot := tgbotbase.NewBot(botCfg)

	//pool := tgbotbase.NewRedisPool(cfg.Redis)

	engine := NewQuestEngine()

	q1 := NewQuest()
	q1.stages["a"] = NewStage("Как зовут моего хозяина?", []string{"илья", "ilya"})
	q1.stages["b"] = NewStage("What year is it?", []string{"2018"})
	engine.AddQuest("test", q1)

	q2 := NewQuest()
	q2.stages["a"] = NewStage("Быть или не быть?", []string{"быть", "не быть"})
	engine.AddQuest("test2", q2)

	tgbot.AddHandler(tgbotbase.NewIncomingMessageDealer(NewStartHandler(engine)))
	tgbot.AddHandler(tgbotbase.NewIncomingMessageDealer(NewAnswerHandler(engine)))

	tgbot.Start()

	log.Info("Daily budget bot has stopped")
}
