package scenario

import (
	"testing"
)

type scenario struct {
	test              *testing.T
	title             string
	currentTitle 	  string
	arguments         Arguments
	argumentProvider  ArgumentProvider
	action            Action
	asserter 		  Asserter
	panicHandler 	  PanicHandler
	expectedResponses Responses
	expectedPanic 	  any
}

func New(test *testing.T, title string) *scenario {
	s := &scenario{
		test:      test,
		title:     title,
		currentTitle: title,
		arguments: make(Arguments, 0),
		action: func(args ...Argument) Responses {
			return Output()
		},
		expectedResponses: make(Responses, 0),
		expectedPanic: nil,
	}
	s.argumentProvider = s.defaultArgumentProvider
	s.asserter = s.defaultAsserter
	s.panicHandler = s.defaultPanicHandler
	return s
}
