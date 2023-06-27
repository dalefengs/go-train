package accesslog

import (
	"fmt"
	"go-train/web"
	"testing"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()

	server := web.NewHTTPServer(":9090", web.ServerWithMiddleware(mdl))
	server.Get("/a/b/*", func(ctx *web.Context) {
		err := ctx.ResponseJSONOK(map[string]string{
			"name": "213",
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
}
