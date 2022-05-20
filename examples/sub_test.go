package examples

import (
	"testing"

	"github.com/kaiquegarcia/gopest/scenario"
	"github.com/kaiquegarcia/gopest/tablescenario"
	"github.com/stretchr/testify/assert"
)

func sub(a, b int) (int, int, int) {
	return a - b, a, b
}

func TestSingleSubWithScenario(t *testing.T) {
	scenario.New(t, "should return 0").
		Given(1, 1).
		When(func(args ...scenario.Argument) scenario.Responses {
			a := args[0].(int)
			b := args[1].(int)
			return scenario.Output(sub(a, b))
		}).
		Expect(0, 1, 1).
		Run()
}

func TestSingleSubWithoutScenario(t *testing.T) {
	// Arrange
	a := 1
	b := 1

	// Act
	result, rA, rB := sub(a, b)

	// Assert
	assert.Equal(t, 0, result)
	assert.Equal(t, 1, rA)
	assert.Equal(t, 1, rB)
}

func TestMultipleSubsWithoutScenario(t *testing.T) {
	type Case struct {
		title  string
		input  []int
		output []int
	}

	cases := []Case{
		{"1-1=0", []int{1, 1}, []int{0, 1, 1}},
		{"-1-1=-2", []int{-1, 1}, []int{-2, -1, 1}},
		{"-1-(-1)=0", []int{-1, -1}, []int{0, -1, -1}},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// Act
			result, rA, rB := sub(c.input[0], c.input[1])

			// Assert
			assert.Equal(t, c.output[0], result)
			assert.Equal(t, c.output[1], rA)
			assert.Equal(t, c.output[2], rB)
		})
	}
}

func TestMultipleSubsWithScenario(t *testing.T) {
	tablescenario.New(t).
		GivenCases(
			tablescenario.Case("1-1=0", scenario.Input(1, 1), scenario.Output(0, 1, 1)),
			tablescenario.Case("-1-1=-2", scenario.Input(-1, 1), scenario.Output(-2, -1, 1)),
			tablescenario.Case("-1-(-1)=0", scenario.Input(-1, -1), scenario.Output(0, -1, -1)),
		).
		When(func(args ...scenario.Argument) scenario.Responses {
			a := args[0].(int)
			b := args[1].(int)
			return scenario.Output(sub(a, b))
		}).
		Run()
}
