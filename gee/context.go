package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Json map[string]interface{}

type Context struct {
	Res        http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	handlers   []HandlerFunc // middleware
	index      int
	engine     *Engine
}

func newContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Res:    res,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (ctx *Context) Next() {
	ctx.index++
	size := len(ctx.handlers)

	for ; ctx.index < size; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

func (ctx *Context) Fail(code int, err string) {
	ctx.index = len(ctx.handlers)
	ctx.Json(code, Json{"message": err})
}

func (ctx *Context) Param(key string) string {
	value, _ := ctx.Params[key]
	return value
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Res.WriteHeader(code)
}

func (ctx *Context) setHeader(key string, value string) {
	ctx.Res.Header().Set(key, value)
}

func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.setHeader("Content-Type", "text/plain; charset=utf-8")
	ctx.Status(code)
	_, _ = ctx.Res.Write([]byte(fmt.Sprintf(format, values...)))
}

func (ctx *Context) Json(code int, data interface{}) {
	ctx.setHeader("Content-Type", "application/json; charset=utf-8")
	ctx.Status(code)

	encoder := json.NewEncoder(ctx.Res)

	if err := encoder.Encode(data); err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
	}
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	_, _ = ctx.Res.Write(data)
}

func (ctx *Context) HTML(code int, name string, data interface{}) {
	ctx.setHeader("Content-Type", "text/html")
	ctx.Status(code)

	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.Res, name, data); err != nil {
		ctx.Fail(http.StatusInternalServerError, err.Error())
	}
}
