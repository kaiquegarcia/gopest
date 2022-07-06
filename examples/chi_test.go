package examples

import (
	"net/http"
	"testing"

	"github.com/kaiquegarcia/gopest/webscenario"
)

func TestGivenChi(t *testing.T) {
	webscenario.New(t, "chi").
		GivenChi(customizer).
		Call(http.MethodGet, "/").
		ExpectHttpStatus(http.StatusOK).
		ExpectJson(`{
			"text": "Hello world!",
			"request-id": "<<PRESENCE>>"
		}`).
		Run()
}
