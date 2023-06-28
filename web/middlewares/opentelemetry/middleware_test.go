package opentelemetry

import (
	"go-train/web"
	"go.opentelemetry.io/otel"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := MiddlewareBuilder{
		Tracer: tracer,
	}

	s := web.NewHTTPServer(":9090", web.ServerWithMiddleware(builder.Build()))
	s.Get("/user", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()
		c, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		c, third1 := tracer.Start(c, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		c, third2 := tracer.Start(c, "third_layer_2")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()
	})
	err := s.Start()
	if err != nil {
		t.Fatal(err)
	}
}
