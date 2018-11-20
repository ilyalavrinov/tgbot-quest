package main

import (
	"strings"
)

type Stage struct {
	question string
	answers  map[string]bool
}

func NewStage(question string, answers []string) Stage {
	s := Stage{question: question}
	for _, a := range answers {
		s.answers[a] = true
	}
	return s
}

type Quest struct {
	stages []Stage
}

type State struct {
	stageIx  int
	stageLen int
}

func (s State) IsFinished() bool {
	return s.stageIx >= s.stageLen
}

func (q Quest) CheckAnswer(answer string, state State) (newState *State) {
	stage := q.stages[state.stageIx]
	answer = strings.ToLower(answer)

	if _, found := stage.answers[answer]; found {
		newState = &State{stageIx: state.stageIx + 1}
	}
	return
}

func (q Quest) CreateInitialState() State {
	return State{stageIx: 0}
}

func (q Quest) GetQuestion(state State) string {
	return q.stages[state.stageIx].question
}
