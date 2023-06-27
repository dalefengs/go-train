package web

import (
	"strings"
)

// 路由树（森林 ）
type router struct {
	// GET、POST、PUT... 都应该有一棵树
	// http method => 路由树根节点
	trees map[string]*node
}

func NewRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {
	if !strings.HasPrefix(path, "/") {
		panic("路由必须是 / 开头")
	}
	if handleFunc == nil {
		panic("HandleFunc is nil")
	}
	// 查找树
	root, ok := r.trees[method]
	if !ok {
		// 说明没有根结点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	// 根节点特殊处理
	if path == "/" {
		if root.handler != nil {
			panic("路由重复注册")
		}
		root.handler = handleFunc
		root.route = "/"
		return
	}
	path = strings.Trim(path, "/") // 去掉左右 /
	segs := strings.Split(path, "/")
	for _, p := range segs {
		if p == "" {
			panic("不能有连续的 /")
		}
		children := root.childOrCreate(p)
		root = children
	}
	// 递归结束 找到最后的节点
	if root.handler != nil {
		panic("路由重复注册")
	}
	root.handler = handleFunc
	root.route = path
}

// 递归查找路由
func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		// 命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			// path = :id  去掉冒号
			pathParams[child.path[1:]] = seg
		}
		// 找到了继续深入
		root = child
	}
	// 只返回 true 只代表我有这个节点，但不一定有 handlerFunc
	// root.handler != nil 代表我既有节点，也有业务逻辑 handlerFunc
	m := &matchInfo{
		n:          root,
		pathParams: pathParams,
	}
	return m, true
}

// childOf 查找路由,优先匹配静态匹配。
// 子节点
// 标记是否是路径参数
// 标记是否命中
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChildren != nil {
			return n.paramChildren, true, true
		}
		return n.starChildren, false, n.starChildren != nil
	}
	child, ok := n.children[path]
	if !ok { // 静态匹配匹配失败，
		if n.paramChildren != nil { // 匹配路径参数
			return n.paramChildren, true, true
		}
		return n.starChildren, false, n.starChildren != nil
	}
	return child, false, ok
}

// 没有找到则创建路由
func (n *node) childOrCreate(seg string) *node {
	if seg[0] == ':' {
		if n.starChildren != nil {
			panic("web：不允许同时注册路径参数和通配符匹配，已有通配符匹配")
		}
		n.paramChildren = &node{
			path: seg,
		}
		return n.paramChildren
	}
	if seg == "*" {
		if n.paramChildren != nil {
			panic("web：不允许同时注册路径参数和通配符匹配, 已有路径参数")
		}
		n.starChildren = &node{
			path: seg,
		}
		return n.starChildren
	}
	if n.children == nil {
		// 不存在就新建
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		// 不存在就新建
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

// 路由节点
type node struct {
	route string

	path string
	// 子节点的映射 - 静态映射
	children map[string]*node
	// 通配符映射， 不允许 /user/*/home
	starChildren *node

	// 路径参数
	paramChildren *node

	// 用户注册的业务逻辑
	handler HandleFunc
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
