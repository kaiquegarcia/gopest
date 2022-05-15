package tablescenario

import "github.com/kaiquegarcia/gopest/scenario"

type tableCase struct {
	title string
	args scenario.Arguments
	expectations scenario.Responses
}