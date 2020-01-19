package goRoute

import "net/http"

type ControllerFunc func(requestInterface RequestInterface) interface{}

type RequestInterface interface {
	GetRequest() *http.Request
}
