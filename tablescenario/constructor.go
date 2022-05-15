package tablescenario

import (
	"testing"

	"github.com/kaiquegarcia/gopest/scenario"
)

// Package inspired by two awesome guys: Joubert RedRat (@joubertredrat) and Guilherme Rodrigues (@guil95)

type tableScenario struct {
	test   *testing.T
	cases  []tableCase
	action scenario.Action
}

func New(t *testing.T) *tableScenario {
	return &tableScenario{
		test: t,
		cases: make([]tableCase, 0),
		action: func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output()
		},
	}
}

func Case(title string, args scenario.Arguments, expectations scenario.Responses) tableCase {
	return tableCase{
		title: title,
		args: args,
		expectations: expectations,
	}
}