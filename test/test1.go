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

func main() {
	g := gee.New()
	g.Use(gee.Logger())
	g.SetFunMap(template.FuncMap{
		"FormatAsData": FormatAsDate,
	})
	g.LoadHTMLGlob("templates/*")
	g.Static("/assets", "./static")

	{
		stu1 := &student{Name: "Geektutu", Age: 20}
		stu2 := &student{Name: "Jack", Age: 22}

		g.GET("/students", func(c *gee.Context) {
			c.HTML(http.StatusOK, "arr.tmpl", gee.H{
				"title":  "gee",
				"stuArr": [2]*student{stu1, stu2},
			})
		})

		g.GET("/date", func(c *gee.Context) {
			c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
				"title": "gee",
				"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
			})
		})
	}

	{
		g.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "css.tmpl", nil)
		})
		g.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
		g.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"username": c.PostForm("username"), "password": c.PostForm("password")})
		})
		g.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		g.GET("/assets/*filepath", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
		})
	}

	v1 := g.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "css.tmpl", nil)
		})

		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := g.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"username": c.PostForm("username"), "password": c.PostForm("password")})
		})
	}
	g.Run(":9999")
}

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
