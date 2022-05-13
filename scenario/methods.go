package scenario

import (
	"fmt"
	"strconv"

	"github.com/stretchr/testify/assert"
)

func (s *scenario) Given(args ...Argument) *scenario {
	s.arguments = args
	return s
}

func (s *scenario) GivenProvider(argumentProvider ArgumentProvider) *scenario {
	s.argumentProvider = argumentProvider
	return s
}

func (s *scenario) When(action Action) *scenario {
	s.action = action
	return s
}

func (s *scenario) Expect(responses ...Response) *scenario {
	s.expectedResponses = responses
	return s
}

func (s *scenario) ExpectTrue() *scenario {
	return s.Expect(true)
}


func (s *scenario) ExpectFalse() *scenario {
	return s.Expect(false)
}


func (s *scenario) ExpectZero() *scenario {
	return s.Expect(0)
}


func (s *scenario) ExpectNil() *scenario {
	return s.Expect(nil)
}

func (s *scenario) ExpectPanic(expectedPanic any) *scenario {
	s.expectedPanic = expectedPanic
	return s
}

func (s *scenario) ExpectWith(asserter Asserter) *scenario {
	s.asserter = asserter
	return s
}

func (s *scenario) ExpectPanicWith(panicHandler PanicHandler) *scenario {
	s.panicHandler = panicHandler
	s.expectedPanic = true
	return s
}

func (s *scenario) Run() {
	tests := s.argumentProvider()
	if len(tests) == 1 {
		s.runSingleTest(tests[0])
		return
	}
	s.runMultipleTests(tests)

}

func (s *scenario) runMultipleTests(tests []Arguments) {
	for index, args := range tests {
		s.currentTitle = s.title + " [" + strconv.Itoa(index) + "]"
		s.runSingleTest(args)
	}
}

func (s *scenario) runSingleTest(args Arguments) {
	defer s.mainPanicHandler()
	fmt.Printf("starting scenario %s\n", s.currentTitle)
	responses := s.action(args...)
	s.asserter(responses)
}

func (s *scenario) mainPanicHandler() {
	err := recover()
	if err == nil && s.expectedPanic == nil {
		return
	}

	s.panicHandler(err)
}

func (s *scenario) defaultPanicHandler(err interface{}) {
	if s.expectedPanic == nil {
		fmt.Printf("scenario %s - unexpected panic %v", s.currentTitle, err)
		s.test.FailNow()
	}

	if err == nil {
		fmt.Printf("scenario %s - expected panic but nothing happened", s.currentTitle)
		s.test.FailNow()
	}

	assert.Equalf(s.test, s.expectedPanic, err, "scenario %s - panic validation", s.currentTitle)
}

func (s *scenario) defaultAsserter(responses Responses) {
	assert.Equalf(s.test, s.expectedResponses, responses, "scenario %s - asserting response\n", s.currentTitle)
}

func (s *scenario) defaultArgumentProvider() []Arguments {
	return []Arguments{s.arguments}
}