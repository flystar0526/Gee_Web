package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePath(path string) []string {
	token := strings.Split(path, "/")
	parts := make([]string, 0)

	for _, part := range token {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (router *router) addRoute(method string, path string, handler HandlerFunc) {
	parts := parsePath(path)
	key := method + "-" + path
	_, ok := router.roots[method]

	if !ok {
		router.roots[method] = &node{}
	}

	router.roots[method].insert(path, parts, 0)
	router.handlers[key] = handler
}

func (router *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)

	root, ok := router.roots[method]

	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)

	if node == nil {
		return nil, nil
	}

	parts := parsePath(node.path)

	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[index]
		} else if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}

	return node, params
}

func (router *router) handle(ctx *Context) {
	node, params := router.getRoute(ctx.Method, ctx.Path)

	if node != nil {
		key := ctx.Method + "-" + node.path
		ctx.Params = params
		ctx.handlers = append(ctx.handlers, router.handlers[key])
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		})
	}

	ctx.Next()
}
