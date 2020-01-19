package goRoute

import (
	"fmt"
	"net/http"
	"strings"
)

const GET Method = "GET"
const POST Method = "POST"
const PUT Method = "PUT"
const PATCH Method = "PATCH"
const DELETE Method = "DELETE"

type Method string

type RouteParams map[string]interface{}

type RouteGroup struct {
	Prefix     string
	name       string
	middleware []string
}

type Route struct {
	Address    string
	name       string
	Controller func(requestInterface RequestInterface) interface{}
	group      *RouteGroup
	Method     Method
	middleware []string
}

func NewRouteTree() *RouteNode {
	return &RouteNode{
		key:   "",
		route: &Route{},
		next:  nil,
		child: nil,
	}
}

func NewRouteGroup(prefix string) RouteGroup {

	return RouteGroup{Prefix: prefix, middleware: make([]string, 0)}
}

func NewRoute() Route {

	return Route{middleware: make([]string, 0)}
}

func (rg *RouteGroup) Route() Route {
	r := NewRoute()
	r.group = rg

	if rg.middleware != nil {
		r.middleware = rg.middleware
	}

	return r
}

func (rg RouteGroup) Group(prefix string) RouteGroup {
	prefix = fmt.Sprintf("%s/%s", checkPrefix(rg.Prefix), prefix)
	newRg := NewRouteGroup(prefix)

	if rg.middleware != nil {
		newRg.middleware = rg.middleware
	}

	if len(rg.name) > 0 {
		newRg.name = rg.name
	}

	return newRg
}

func (rg *RouteGroup) Middleware(m string) *RouteGroup {
	rg.middleware = append(rg.middleware, m)

	return rg
}

func (rg *RouteGroup) Name(name string) *RouteGroup {
	if len(rg.name) > 0 {
		name = fmt.Sprintf("%s.%s", rg.name, name)
	}

	rg.name = name

	return rg
}

func (r Route) Group(prefix string) *RouteGroup {
	rg := NewRouteGroup(prefix)

	if r.group != nil && len(r.group.middleware) > 0 {

		rg.middleware = r.group.middleware
	}

	return &rg
}

func (r Route) Get(address string, controller func(requestInterface RequestInterface) interface{}) *Route {

	return addNewRoute(r, address, controller, GET)
}

func (r Route) Post(address string, controller func(requestInterface RequestInterface) interface{}) *Route {

	return addNewRoute(r, address, controller, POST)
}

func (r Route) Put(address string, controller func(requestInterface RequestInterface) interface{}) *Route {

	return addNewRoute(r, address, controller, PUT)
}

func (r Route) Patch(address string, controller func(requestInterface RequestInterface) interface{}) *Route {

	return addNewRoute(r, address, controller, PATCH)
}

func (r Route) Delete(address string, controller func(requestInterface RequestInterface) interface{}) *Route {

	return addNewRoute(r, address, controller, DELETE)
}

func (r *Route) Name(name string) *Route {
	if r.group != nil && len(r.group.name) > 0 {
		name = fmt.Sprintf("%s.%s", r.group.name, name)
	}

	r.name = name

	GetRouteNames().Add(r.name, r.Address)

	return r
}

func (r *Route) Middleware(m string) *Route {
	r.middleware = append(r.middleware, m)

	return r
}

func (r *Route) GetMiddleware() []string {
	return r.middleware
}

func checkPrefix(address string) string {
	if strings.HasPrefix(address, "/") {
		address = strings.TrimPrefix(address, "/")
	}

	if strings.HasSuffix(address, "/") {
		address = strings.TrimSuffix(address, "/")
	}

	return address
}

func checkParams(path string) string {
	if strings.HasPrefix(path, "$") {
		return strings.TrimPrefix(path, "$")
	}

	return ""
}

func getPath(url string, method Method) []string {
	url = fmt.Sprintf("%s/%s", method, url)

	return strings.Split(url, "/")
}

func addNewRoute(r Route, address string, controller func(requestInterface RequestInterface) interface{}, method Method) *Route {
	if len(r.group.Prefix) > 0 {
		address = checkPrefix(r.group.Prefix) + "/" + checkPrefix(address)
	}

	if r.group.middleware != nil {
		r.middleware = r.group.middleware
	}

	r.Address = checkPrefix(address)
	r.Controller = controller
	r.Method = method

	path := getPath(r.Address, method)
	GetRouteTree().AddFromPath(path, &r)

	return &r
}

func Match(request *http.Request) (*Route, bool, RouteParams) {
	url := request.URL.Path
	url = checkPrefix(url)
	path := getPath(url, Method(request.Method))

	node, params := GetRouteTree().FindFromPath(path)

	if node == nil || node.GetRoute() == nil {
		return nil, false, nil
	}

	return node.route, true, params
}
