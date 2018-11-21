package main

import (
	"flag"
	"github.com/admirallarimda/tgbot-quest/internal/pkg/quest"
	log "github.com/sirupsen/logrus"
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
	flag.Parse()

	if (*argQuest == "") || (*argQuestion == "") || (*argAnswers == "") {
		log.Panic("One of mandatory arguments is not set")
	}

	if *argStage == "" {
		*argStage = time.Now().Format(timeFormat)
	}

	answers := strings.Split(*argAnswers, ";")

	q := quest.NewQuest()
	q.stages[*argStage] = quest.NewStage(*argQuestion, answers)

}
