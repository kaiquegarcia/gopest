package webscenario

import (
	"net/http"

	"gopkg.in/xmlpath.v2"
)

type any interface{}
type FormBody map[string]string
type method string
type headers map[string]string
type xmlNodes map[string]node

type node struct {
	path string
	expectedValue string
	compiler *xmlpath.Path
}
// TODO: type htmlNodes map[string]node

const (
	MethodConnect method = http.MethodConnect
	MethodDelete method = http.MethodDelete
	MethodGet method = http.MethodGet
	MethodHead method = http.MethodHead
	MethodOptions method = http.MethodOptions
	MethodPatch method = http.MethodPatch
	MethodPost method = http.MethodPost
	MethodPut method = http.MethodPut
	MethodTrace method = http.MethodTrace
)