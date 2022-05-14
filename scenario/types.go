package scenario

type any interface{}
type Argument any
type Arguments []Argument
type Response any
type Responses []Response
type Action func(args ...Argument) Responses
type Asserter func(responses Responses)
type PanicHandler func(err interface{})
type ArgumentProvider func() []Arguments

type ScenarioInterface interface {
	Given(...Argument) *scenario
	GivenProvider(ArgumentProvider) *scenario
	When(Action) *scenario
	Expect(...Response) *scenario
	ExpectTrue() *scenario
	ExpectFalse() *scenario
	ExpectZero() *scenario
	ExpectNil() *scenario
	ExpectPanic(any) *scenario
	ExpectWith(Asserter) *scenario
	ExpectPanicWith(PanicHandler) *scenario
	Run()
}

func Input(args ...Argument) Arguments {
	return args
}

func Output(responses ...Response) Responses {
	return responses
}
