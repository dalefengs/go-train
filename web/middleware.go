package web

// Middleware AOP 方案 责任链模式
type Middleware func(next HandleFunc) HandleFunc
