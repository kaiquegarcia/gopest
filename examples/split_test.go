package examples

import (
	"testing"

	"github.com/kaiquegarcia/gopest/scenario"
	"github.com/kaiquegarcia/gopest/tablescenario"
)

func TestSplit(t *testing.T) {
	// inspired in https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
	tablescenario.New(t).
		GivenCases(
			tablescenario.Case("'a/b/c' should split by / to ['a','b','c']", scenario.Input("a/b/c", "/"), scenario.Output([]string{"a", "b", "c"})),
			tablescenario.Case("'a/b/c' should split by , to ['a/b/c']", scenario.Input("a/b/c", ","), scenario.Output([]string{"a/b/c"})),
			tablescenario.Case("'abc' should split by / to ['abc']", scenario.Input("abc", "/"), scenario.Output([]string{"abc"})),
		).
		When(func(args ...scenario.Argument) scenario.Responses {
			input := args[0].(string)
			sep := args[1].(string)
			return scenario.Output(Split(input, sep))
		}).
		Run()
}
