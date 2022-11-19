package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*Node
	handlers map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{
		roots:    map[string]*Node{},
		handlers: map[string]HandlerFunc{},
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := []string{}
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &Node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) getRoute(method string, pattern string) (*Node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	searchParts := parsePattern(pattern) // 待查pattern
	params := map[string]string{}
	node := root.search(searchParts, 0)
	if node == nil {
		return nil, nil
	}

	parts := parsePattern(node.pattern) // Trie树的匹配 pattern
	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[index]
		}
		if part[0] == '*' {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}
	return node, params
}

func (r *Router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params //获得路由之后， 返回全路径上的通配参数
		key := c.Method + "-" + node.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
