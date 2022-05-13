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

func Input(args ...Argument) Arguments {
	return args
}

func Output(responses ...Response) Responses {
	return responses
}
