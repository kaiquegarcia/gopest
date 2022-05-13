package examples

import (
	"testing"

	"github.com/kaiquegarcia/gopest/scenario"
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
	scenario.New(t, "should assert panic").
		Given(123, 27).
		When(func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output(Sum(args[0].(int), args[1].(int)))
		}).
		ExpectWith(func (responses scenario.Responses) {
			assert.Equal(t, 150, responses[0].(int), "asserting from closure")
		}).
		Run()
}

func TestShouldExpectPanicWith(t *testing.T) {
	scenario.New(t, "should assert panic").
		When(func(args ...scenario.Argument) scenario.Responses {
			panic("please don't do that")
		}).
		ExpectPanicWith(func (err interface{}) {
			assert.Equal(t, "please don't do that", err, "asserting from closure")
		}).
		Run()
}

func TestShouldExpectTrueWithArgumentProvider(t *testing.T) {
	dataProvider := func() []scenario.Arguments {
		return []scenario.Arguments{
			{123, 27, 150},
			{1, 1, 2},
			{3, 7, 10},
		}
	}

	scenario.New(t, "should assert true").
		GivenProvider(dataProvider).
		When(func(args ...scenario.Argument) scenario.Responses {
			return scenario.Output(Sum(args[0].(int), args[1].(int)) == args[2].(int))
		}).
		ExpectTrue().
		Run()
}