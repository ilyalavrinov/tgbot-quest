package main

import (
	"strings"
)

type Stage struct {
	question string
	answers  map[string]bool
}

func NewStage(question string, answers []string) Stage {
	s := Stage{
		question: question,
		answers:  make(map[string]bool, len(answers))}
	for _, a := range answers {
		s.answers[strings.ToLower(a)] = true
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
		newState = &State{}
		*newState = state
		newState.stageIx++
	}
	return
}

func (q Quest) CreateInitialState() State {
	return State{
		stageIx:  0,
		stageLen: len(q.stages)}
}

func (q Quest) GetQuestion(state State) string {
	return q.stages[state.stageIx].question
}
