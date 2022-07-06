# GoPest

GoPest is a semantic testing library for Go, specially inspired on **Pest** (*Jest for PHP*) with [pest-plugin-gwt](https://github.com/milroyfraser/pest-plugin-gwt).

## Installation

Install the package you need using the commands:

```bash
go get github.com/kaiquegarcia/gopest/scenario
go get github.com/kaiquegarcia/gopest/tablescenario
go get github.com/kaiquegarcia/gopest/webscenario
```

## Usage

### Scenario pkg

#### Resume

This is the core package. The other pkgs use it to run your scenarios with other perspective. If you need to run multiple scenarios, use the [Table-Driven Scenario](#tablescenario-pkg). If you need to run end to end tests calling web requests to your web framework, use the [Web Scenario](#webscenario-pkg). That's all for now.

#### Guide

1. The first step is to create the test scenario: you need to call the constructor function `scenario.New`. It requires two arguments: the `*testing.T` you receive in every test function you create and a `string` describing the test you're attempting to create;
2. Then you need to inform the test input calling the method `Given`. You can inject as many arguments you need;
3. The third step is to register the action calling the method `When`. Pay atention to the closure signature: `func (args ...scenario.Argument) scenario.Responses`. This means you'll receive all the inputs in the slice `args` and you need to return a list of responses with the method `scenario.Output()`. The same way: you can output as many responses you need;
4. Now you can register all the expectations. You have a few options:
    - use `Expect` to register each output you expect from the `action` registered before;
    - use `ExpectWith` to register a closure responsible to do all the assertions. Signature: `scenario.Asserter` or `func (t *testing.T, responses ...scenario.Response)`;
    - use `ExpectTrue` to register you want a single output resulting `true`;
    - use `ExpectFalse` to register you want a single output resulting `false`;
    - use `ExpectNil` to register you want a single output resulting `nil`;
    - use `ExpectZero` to register you want a single output resulting `0` (int);
    - use `ExpectPanic` to register you want to assert a specific `panic` is triggered;
    - use `ExpectPanicWith` to register a closure responsible to do the panic assertion. Signature: `scenario.PanicHandler` or `func(t *testing.T, err interface{})`;
5. And finally, you call `scenario.Run()` to run the tests.

#### Example

```go
import (
    "testing"

    "github.com/kaiquegarcia/gopest/scenario"
)

// Sum is the function we want to test
func Sum(a, b int) int {
    return a+b
}

// TestSum is a single test
func TestSum(t *testing.T) {
    scenario.New(t, "testing sum 1+1=2").
        Given(1, 1).
        When(func (args ...scenario.Argument) scenario.Responses {
            a := args[0].(int)
            b := args[1].(int)
            return scenario.Output(Sum(a, b))
        }).
        Expect(2).
        Run()
}
```

You can look for further examples in the file `examples/sum_test.go`.

### TableScenario pkg

#### Resume

This package is specially made to run Table-Driven Tests, thanks to [Joubert RedRat (@joubertredrat)](https://github.com/joubertredrat) and [Guilherme Rodrigues (@guil95)](https://github.com/guil95) feedbacks. If you want to check more than one case in the same test, this is for you.

#### Guide

1. The first step looks like the previous one: call `tablescenario.New()`. Here we have only one argument, the `*testing.T` you receive in each test function;
2. Then you call the method `GivenCases`, where you're supposed to put each case you want to test. PS.: each case must be instanced by another constructor calling the function `tablescenario.Case()` with 3 arguments: a `string` (the title of the case), a `scenario.Arguments` created by `scenario.Input()` and `scenario.Responses` defined by `scenario.Output()`;
3. After setting up the cases, call the method `When`, which means the same mentioned in the previous section. Signature: `func (args ...scenario.Argument) scenario.Responses`;
4. And go! I-I mean, `Run()`.

#### Example

```go
import (
    "testing"

    "github.com/kaiquegarcia/gopest/scenario"
    "github.com/kaiquegarcia/gopest/tablescenario"
)

// Sub is the function we want to test
func Sub(a, b int) int {
    return a-b
}

// TestSub is a table-driven test
func TestSub(t *testing.T) {
    tablescenario.New(t).
        GivenCases(
            tablescenario.Case("sub 1-1=0", scenario.Input(1,1), scenario.Output(0)),
            tablescenario.Case("sub (-1)-(-1)=0", scenario.Input(-1,-1), scenario.Output(0)),
            tablescenario.Case("sub (-1)-1=-2", scenario.Input(-1,1), scenario.Output(-2)),
        ).
        When(func (args ...scenario.Argument) scenario.Responses {
            a := args[0].(int)
            b := args[1].(int)
            return scenario.Output(Sub(a, b))
        }).
        Run()
}
```

### WebScenario pkg

#### Resume

This package is specially made to run End-to-End Tests, making test requests to your web framework.

#### Guide

1. As all the previous guides, start the scenario calling `webscenario.New()` with 2 arguments: the `*testing.T` you receive in every test function and a `string` to describe which scenario you're testing;
2. The second step is to call the method `Given{webframework}`, where `{webframework}` should be the framework you use in your webserve. Currently we support `Fiber` and `Chi` only (please open an issue or start your own implementation to add more options!), so you can call `GivenFiber` or `GivenChi`.

    PS.: in case of `GivenFiber`, you need to inject an instance of `*fiber.App` with the desired route implemented;

    PS.2: in case of `GivenChi`, you need to inject a customizer function to inject the desired routes to be implemented. You can also inject middlewares.
3. Then you need to call the method `Call()` with two arguments: a `string` with the desired request method (`http.MethodPost`, for example) and another `string` with the desired route (`/my-route`);
4. After that, you can inject values to your requests with the following options:
    - use `Header` to set a `string` to a specific key, that'll compose the request headers (`Content-type=text/html`);
    - use `Query` to set a `string` value to a specific key, that'll compose the request query (`?key=value`);
    - use `XmlBody` to define a XML `string` as the body of the request. That will call `Header("Content-Type", "application/xml")` automatically. You can override it by calling the previous `Header` option after this;
    - use `JsonBody` to define a JSON `string` as the body of the request. That will call `Header("Content-Type", "application/json; charset=utf-8")` automatically. You can override it by calling the previous `Header` option after this;
    - use `FormBody` to define a complete Form values as the body of the request. That will call `Header("Content-Type", "application/x-www-form-urlencoded")` automatically. You can override it by calling the previous `Header` option after this;
5. Now you can register all your expectations. Here's your options:
    - use `ExpectHeader` to register each key and expected value you expect to receive in the request Header;
    - use `ExpectHttpStatus` to register which HTTP Status you expect to receive in the request;
    - use `ExpectJson` to register which JSON body you expect to receive in the request Body. This assertion is made with `jsonassert.Assertf()` from [kinbiko/jsonassert@v1.1.0](https://github.com/kinbiko/jsonassert) module. You can use all the features described here;
    - use `ExpectXmlNode` to register each XML Path and expected value you expect to receive in the request Body. This assertion is made with `xmlpath.MustCompile` from [gopkg.in/xmlpath.v2](https://gopkg.in/xmlpath.v2) module. You can use all the features described here;
6. And finally call `Run()`!

#### Example

```go
import (
    "testing"
    "net/http"

    "github.com/gofiber/fiber/v2"
    "github.com/kaiquegarcia/gopest/webscenario"
)

func TestGetMyRoute(t *testing.T) {
    app := fiber.New()
    app.Get("/my-route", function (ctx *fiber.Ctx) {
        return ctx.JSON(`{"message": "ok"}`)
    })
    webscenario.New(t, "test route GET my-route").
        GivenFiber(app).
        Call(http.MethodGet, "/my-route").
		ExpectHttpStatus(http.StatusOK).
		ExpectJson(`{"message": "ok"}`).
		Run()
}
```

You can look for further examples in the file `examples/simple_router_test.go`.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.


## Roadmap

- Add `ExpectHtmlNode` to `webscenario` pkg;
- Add `ExpectPlainText` to `webscenario` pkg;
- Add `ExpectPermanentRedirect` to `webscenario` pkg;
- Add `ExpectTemporaryRedirect` to `webscenario` pkg;
- Add `GivenGin` to `webscenario` pkg;
- Add `GivenHttpServer` to `webscenario` pkg;
- Add more examples running each case described in the docs.

## License
[MIT](https://choosealicense.com/licenses/mit/)