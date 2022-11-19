package gee

import (
	"log"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc //支持中间件
	engine      *Engine      //访问router
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) addRoute(method string, comp string, handler HandleFunc) {
	pattern := rg.prefix + comp
	log.Printf("Route %4s - %s", method, comp)
	rg.engine.router.addRoute(method, pattern, handler)
}

func (rg *RouterGroup) GET(pattern string, handler HandleFunc) {
	rg.addRoute("GET", pattern, handler)
}

func (rg *RouterGroup) POST(pattern string, handler HandleFunc) {
	rg.addRoute("POST", pattern, handler)
}
