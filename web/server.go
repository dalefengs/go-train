package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// 确保结构体一定实现了 Server 接口
var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start() error

	AddRoute(method string, path string, handleFunc HandleFunc) // 路由注册
}

type HTTPServer struct {
	Addr string
	router
	Middlewares []Middleware
}

type HTTPServerOption func(server *HTTPServer)

func NewHTTPServer(addr string, opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: NewRouter(),
		Addr:   addr,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.Middlewares = mdls
	}
}

// ServeHTTP 处理请求入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 构建 Context
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 从后往前回溯
	// 把后一个作为前一个的 next 构造好链条
	root := h.Serve
	for i := len(h.Middlewares) - 1; i >= 0; i-- {
		root = h.Middlewares[i](root)
	}
	root(ctx)
}

func (h *HTTPServer) Serve(ctx *Context) {
	// 查找路由，并命中逻辑
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok {
		// 路由没有命中 404
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
	ctx.PathParams = info.pathParams
	ctx.MatchedRoute = info.n.route
	info.n.handler(ctx)
}

// Get 注册路由
func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodPost, path, handleFunc)
}

func (h *HTTPServer) Put(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodPut, path, handleFunc)
}

func (h *HTTPServer) Delete(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodDelete, path, handleFunc)
}

func (h *HTTPServer) Start() error {
	listen, err := net.Listen("tcp", h.Addr)
	if err != nil {
		return err
	}
	return http.Serve(listen, h)
}
