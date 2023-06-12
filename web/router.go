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

func NewRouter() *router {
	return &router{
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
}

type node struct {
	path string
	// 子节点的映射
	children map[string]*node
	// 用户注册的业务逻辑
	handler HandleFunc
}

// 递归查找路由
func (r *router) findRoute(method, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return root, true
	}
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		child, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		// 找到了继续深入
		root = child
	}
	// 只返回 true 只代表我有这个节点，但不一定有 handlerFunc
	// root.handler != nil 代表我既有节点，也有业务逻辑 handlerFunc
	return root, true
}

// 查找路由
func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}
	child, ok := n.children[path]
	return child, ok
}

// 没有找到则创建路由
func (n *node) childOrCreate(seg string) *node {
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
