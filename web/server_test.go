package web

import (
	"testing"
)

func TestStart(t *testing.T) {
	s := NewHTTPServer(":9090")
	s.Get("/user/home", func(ctx *Context) {
		ctx.Resp.Write([]byte("home"))
		return
	})
	s.Get("/userInfo/*", func(ctx *Context) {
		ctx.Resp.Write([]byte((ctx.Req.URL.Path)))
		return
	})
	s.Post("/name/:name", func(ctx *Context) {
		ctx.Resp.Write([]byte((ctx.PathParams["name"])))
		return
	})
	s.Start()
}
