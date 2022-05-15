package scenario

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (s *scenario) Given(args ...Argument) *scenario {
	s.arguments = args
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
	s.test.Run(s.title, func(t *testing.T) {
		defer s.mainPanicHandler(t)
		responses := s.action(s.arguments...)
		s.asserter(t, responses)
	})
}

func (s *scenario) mainPanicHandler(t *testing.T) {
	err := recover()
	if err == nil && s.expectedPanic == nil {
		return
	}

	s.panicHandler(t, err)
}

func (s *scenario) defaultPanicHandler(t *testing.T, err interface{}) {
	if s.expectedPanic == nil {
		t.Fatalf("scenario %s - unexpected panic %v", s.title, err)
		t.FailNow()
	}

	if err == nil {
		t.Fatalf("scenario %s - expected panic but nothing happened", s.title)
		t.FailNow()
	}

	assert.Equalf(t, s.expectedPanic, err, "scenario %s - panic validation", s.title)
}

func (s *scenario) defaultAsserter(t *testing.T, responses Responses) {
	assert.Equalf(t, s.expectedResponses, responses, "scenario %s - asserting response\n", s.title)
}
