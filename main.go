package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func onlyForV2() gee.HandlerFunc {
	return func(ctx *gee.Context) {
		startTime := time.Now()
		log.Printf("[%d] %s in %v for group v2", ctx.StatusCode, ctx.Req.RequestURI, time.Since(startTime))
	}
}

func main() {
	router := gee.Default()
	router.Use(gee.Logger()) // global middleware
	router.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./static")

	student1 := &student{
		Name: "flystar0526",
		Age:  21,
	}

	student2 := &student{
		Name: "jack521",
		Age:  21,
	}

	router.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "css.tmpl", nil)
	})

	router.GET("/students", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "arr.tmpl", gee.Json{
			"title":    "gee",
			"content":  "hello, gee",
			"students": [2]*student{student1, student2},
		})
	})

	router.GET("/date", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "custom_func.tmpl", gee.Json{
			"title": "gee",
			"now":   time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC),
		})
	})

	router.GET("/panic", func(ctx *gee.Context) {
		names := []string{"flystar0526"}
		ctx.String(http.StatusOK, names[100])
	})

	v1 := router.Group("/v1")

	v1.GET("/hello", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "Hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	v1.GET("/hello/:name", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "Hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	})

	v2 := router.Group("/v2")
	v2.Use(onlyForV2())

	v2.POST("/login", func(ctx *gee.Context) {
		ctx.Json(http.StatusOK, gee.Json{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	_ = router.Run(":8080")
}
