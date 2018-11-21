package main

import (
	"sort"
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
	stages map[string]Stage
}

func NewQuest() Quest {
	return Quest{stages: make(map[string]Stage, 0)}
}

type State struct {
	stageIx    int
	stageOrder []string
}

func (s State) IsFinished() bool {
	return s.stageIx >= len(s.stageOrder)
}

func (s State) GetStageID() string {
	return s.stageOrder[s.stageIx]
}

func (s State) Next() *State {
	s2 := s
	s2.stageIx++
	return &s2
}

func (q Quest) CheckAnswer(answer string, state State) (newState *State) {
	stage := q.stages[state.GetStageID()]
	answer = strings.ToLower(answer)

	if _, found := stage.answers[answer]; found {
		newState = state.Next()
	}
	return
}

func (q Quest) CreateInitialState() State {
	order := make([]string, 0, len(q.stages))
	for k, _ := range q.stages {
		order = append(order, k)
	}
	sort.Strings(order)
	return State{
		stageIx:    0,
		stageOrder: order}
}

func (q Quest) GetQuestion(state State) string {
	return q.stages[state.GetStageID()].question
}
