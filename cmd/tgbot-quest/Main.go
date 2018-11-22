package main

import (
	"github.com/admirallarimda/tgbot-quest/internal/pkg/quest"
	"github.com/admirallarimda/tgbotbase"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gcfg.v1"
	"math/rand"
	"time"
)

type config struct {
	tgbotbase.Config
	Redis tgbotbase.RedisConfig
	Owner struct {
		ID int
	}
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
	log.SetLevel(log.DebugLevel)
	log.Info("Starting daily budget bot")

	cfg := readGcfg("bot.cfg")
	botCfg := tgbotbase.Config{TGBot: cfg.TGBot, Proxy_SOCKS5: cfg.Proxy_SOCKS5}
	tgbot := tgbotbase.NewBot(botCfg)

	rand.Seed(int64(time.Now().Second()))

	pool := tgbotbase.NewRedisPool(cfg.Redis)
	resmon := quest.NewTGResultMonitor(tgbot, []tgbotbase.UserID{tgbotbase.UserID(cfg.Owner.ID)})
	engine := quest.NewQuestEngine(pool, resmon)

	tgbot.AddHandler(tgbotbase.NewIncomingMessageDealer(NewStartHandler(engine)))
	tgbot.AddHandler(tgbotbase.NewIncomingMessageDealer(NewAnswerHandler(engine)))

	tgbot.Start()

	log.Info("Daily budget bot has stopped")
}
