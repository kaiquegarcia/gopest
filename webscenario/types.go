package webscenario

import "net/http"

type any interface{}
type FormBody map[string]string
type method string
type headers map[string]string

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