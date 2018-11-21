package main

import (
	"flag"
	"github.com/admirallarimda/tgbot-quest/internal/pkg/quest"
	"github.com/admirallarimda/tgbotbase"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var argQuest = flag.String("quest", "", "ID of the target quest")
var argStage = flag.String("stage", "", "ID of the target stage (optional)")
var argPic = flag.String("pic", "", "Path/URL of the picture which will be attached to a question (optional)")
var argQuestion = flag.String("question", "", "Question itself")
var argAnswers = flag.String("answers", "", "Semicolon (;)-split list of answers")

const timeFormat = "20060102150405.000"

func main() {
	log.SetLevel(log.DebugLevel)
	flag.Parse()

	if (*argQuest == "") || (*argQuestion == "") || (*argAnswers == "") {
		flag.PrintDefaults()
		log.Panic("One of mandatory arguments is not set")
	}

	if *argStage == "" {
		*argStage = time.Now().Format(timeFormat)
	}

	answers := strings.Split(*argAnswers, ";")

	q := quest.NewQuest()
	stage := quest.NewStage(*argQuestion, answers)
	if *argPic != "" {
		if (*argPic)[:4] == "http" {
			log.WithField("pic", *argPic).Panic("HTTP will be handled later")
		}
		f, err := os.Open(*argPic)
		if err != nil {
			log.WithFields(log.Fields{"pic": *argPic, "error": err}).Panic("Error on file open")
		}
		b := make([]byte, 1024*1024)
		n, err := f.Read(b)
		if err != nil {
			log.WithFields(log.Fields{"pic": *argPic, "error": err}).Panic("Error on file reading")
		}
		log.WithFields(log.Fields{"pic": *argPic, "bytes_read": n}).Debug("File has been read")
		stage.AddPicture(b[:n])
	}
	q.AddStage(*argStage, stage)

	cfg := tgbotbase.RedisConfig{"127.0.0.1:6379", ""}
	pool := tgbotbase.NewRedisPool(cfg)
	storage := quest.NewRedisQuestStorage(pool)
	storage.StoreQuest(*quest.NewQuestRecord(*argQuest, q))
}
