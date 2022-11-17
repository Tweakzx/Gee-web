package main

import (
	"gee"
	"net/http"
)

func main() {
	g := gee.New()
	g.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Tweakzx</h1>")
	})
	g.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	g.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"username": c.PostForm("username"), "password": c.PostForm("password")})
	})
	g.Run(":9999")
}
