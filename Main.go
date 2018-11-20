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

	tgbot.AddHandler(tgbotbase.NewIncomingMessageDealer(NewAnswerHandler()))

	tgbot.Start()

	log.Info("Daily budget bot has stopped")
}
