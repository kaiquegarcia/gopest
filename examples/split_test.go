package examples

import (
	"testing"

	"github.com/kaiquegarcia/gopest/scenario"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	// inspired in https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
	provider := func () []scenario.Arguments {
		return []scenario.Arguments{
			{"a/b/c", "/", []string{"a", "b", "c"}},
			{"a/b/c", ",", []string{"a/b/c"}},
			{"abc", "/", []string{"abc"}},
		}
	}

	scenario.New(t, "test split").
		GivenProvider(provider).
		When(func(args ...scenario.Argument) scenario.Responses {
			input := args[0].(string)
			sep := args[1].(string)
			expect := args[2]
			result := Split(input, sep)
			return scenario.Output(expect, result)
		}).
		ExpectWith(func(responses scenario.Responses) {
			assert.Equal(t, responses[0], responses[1])
		}).
		Run()
}