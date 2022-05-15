package examples

import (
	"testing"

	"github.com/kaiquegarcia/gopest/scenario"
	"github.com/kaiquegarcia/gopest/tablescenario"
	"github.com/stretchr/testify/assert"
)

func TestShouldExpectTrue(t *testing.T) {
	scenario.New(t, "should assert true").
		Given(123, 27, 150).
		When(func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output(Sum(args[0].(int), args[1].(int)) == args[2].(int))
		}).
		ExpectTrue().
		Run()
}

func TestShouldExpectFalse(t *testing.T) {
	scenario.New(t, "should assert false").
		Given(123, 27, 151).
		When(func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output(Sum(args[0].(int), args[1].(int)) == args[2].(int))
		}).
		ExpectFalse().
		Run()
}

func TestShouldExpectZero(t *testing.T) {
	scenario.New(t, "should assert zero").
		Given(123, -123).
		When(func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output(Sum(args[0].(int), args[1].(int)))
		}).
		ExpectZero().
		Run()
}

func TestShouldExpectNil(t *testing.T) {
	scenario.New(t, "should assert nil").
		When(func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output(nil)
		}).
		ExpectNil().
		Run()
}

func TestShouldExpectEqual(t *testing.T) {
	scenario.New(t, "should assert equal").
		Given(123, 27).
		When(func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output(Sum(args[0].(int), args[1].(int)))
		}).
		Expect(150).
		Run()
}

func TestShouldExpectPanic(t *testing.T) {
	scenario.New(t, "should assert panic").
		When(func(args ...scenario.Argument) scenario.Responses {
			panic("please don't do that")
		}).
		ExpectPanic("please don't do that").
		Run()
}

func TestShouldExpectWith(t *testing.T) {
	scenario.New(t, "should assert with closure").
		Given(123, 27).
		When(func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output(Sum(args[0].(int), args[1].(int)))
		}).
		ExpectWith(func(t *testing.T, responses scenario.Responses) {
			assert.Equal(t, 150, responses[0].(int), "asserting from closure")
		}).
		Run()
}

func TestShouldExpectPanicWith(t *testing.T) {
	scenario.New(t, "should assert panic with closure").
		When(func(args ...scenario.Argument) scenario.Responses {
			panic("please don't do that")
		}).
		ExpectPanicWith(func(t *testing.T, err interface{}) {
			assert.Equal(t, "please don't do that", err, "asserting from closure")
		}).
		Run()
}

func TestShouldExpectTrueWithArgumentProvider(t *testing.T) {
	tablescenario.New(t).
		GivenCases(
			tablescenario.Case("Sum(123,27) should result 150", scenario.Input(123, 27), scenario.Output(150)),
			tablescenario.Case("Sum(1, 1) should result 2", scenario.Input(1, 1), scenario.Output(2)),
			tablescenario.Case("Sum(3, 7) should result 10", scenario.Input(3, 7), scenario.Output(10)),
		).
		When(func(args ...scenario.Argument) scenario.Responses {
			a := args[0].(int)
			b := args[1].(int)
			return scenario.Output(Sum(a, b))
		}).
		Run()
}
