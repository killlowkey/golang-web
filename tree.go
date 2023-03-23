package web

import (
	"net/http"
	"strings"
)

// Param 用于表示参数路由匹配参数
// 参数路由为 /user/:id，真实 http 请求路径为 /user/1
// 那么 Param key 为 id，value 为 1
type Param struct {
	key   string
	value string
}

type Params []Param

// Get 搜索 Params 返回 key 与之对应的 value
func (p Params) Get(key string) (string, bool) {
	for _, param := range p {
		if param.key == key {
			return param.value, true
		}
	}
	return "", false
}

// Trees 每个方法对应一棵树
type Trees struct {
	trees map[string]*Tree
}

func (t *Trees) findRoute(r *http.Request) (*nodeInfo, bool) {
	tree := t.trees[r.Method]
	if tree == nil {
		return &nodeInfo{}, false
	}

	return tree.root.findRoute(r.URL.Path)
}

func (t *Trees) addRouter(method, path string, middlewares []Middleware, handler HandleFunc) {
	tree, ok := t.trees[method]
	if !ok {
		tree = &Tree{
			method: method,
			root:   &node{path: "/"},
		}
		t.trees[method] = tree
	}

	tree.root.addRoute(path, middlewares, handler)
}

// Tree 每个 http 方法都有自己的路由树
type Tree struct {
	method string // 树绑定的 http 方法
	root   *node  // 树的根节点
}

// node 路由树节点
// 路由树是 web 框架的核心设计，直接决定了框架的性能，树的高度影响路由树查找性能。
// Gin 实现中孩子节点使用 *node 切片进行表示 ，这里进行优化，使用 map 进行查找
// 如果注册大量前缀不同的路由，可以显著的提高性能（树宽度问题）。
type node struct {
	path        string           // 路由绑定路径
	children    map[string]*node // 普通的孩子节点，使用 map 快速查找
	starChild   *node            // 通配符匹配
	paramChild  *node            // 参数匹配
	middlewares []Middleware     // 路由局部 Middleware，例如 Group 方法添加的 Middleware
	handler     HandleFunc       // 业务处理器
	fullPath    string           // 注册路由绑定的路径
}

// addRoute 添加路由，支持如下几种路由：
//  1. 静态路由：/a/b/c
//  2. 通配符路由：/a/*
//  3. 参数路由：/a/:name
//
// 通配符与参数路由是互斥的，要么存在通配符路由，要么存在参数路由
//
// 首先需要解决边界问题，禁用如下几种边界方便后续处理
//  1. path 路径为空
//  2. path 路径非 / 开头，例如 a/b/c
//  3. path 路径以 / 结尾，例如 /a/b/
func (n *node) addRoute(path string, middlewares []Middleware, handler HandleFunc) {
	if path == "" {
		panic("web：path 路径不允许为空")
	}

	if path[0] != '/' {
		panic("web：path 必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web：path 路径不允许为 / 结尾")
	}

	if path == "/" {
		if n.handler != nil {
			panic("web：不允许重复注册路由[/]")
		}
		n.path = path
		n.middlewares = middlewares
		n.handler = handler
		return
	}

	cur := n
	segments := strings.Split(path[1:], "/")
	for _, seg := range segments {
		if seg == "" {
			panic("web：拒绝 /a//b/c 形式的路由")
		}
		cur = cur.insert(seg)
	}

	if cur.handler != nil {
		panic("web：不允许重复注册路由[" + path + "]")
	}
	cur.fullPath = path
	cur.middlewares = middlewares
	cur.handler = handler
}

// insert 插入节点
// 通配符和参数匹配只允许存在一个
func (n *node) insert(path string) *node {
	// 通配符路由：/a/*/b
	if path == "*" {
		if n.paramChild != nil {
			panic("web：通配符与参数匹配冲突，只允许存在一个")
		}
		if n.starChild == nil {
			n.starChild = &node{path: "*"}
		}
		return n.starChild
	}

	// 参数路由：/a/:name
	if path[0] == ':' {
		if n.starChild != nil {
			panic("web：通配符与参数匹配冲突，只允许存在一个")
		}
		if n.paramChild == nil {
			n.paramChild = &node{path: path}
		}
		return n.paramChild
	}

	// 节点初始化
	if n.children == nil {
		n.children = map[string]*node{}
	}
	cur, ok := n.children[path]
	if !ok {
		cur = &node{path: path}
		n.children[path] = cur
	}
	return cur
}

// findRoute 查找路由
func (n *node) findRoute(path string) (*nodeInfo, bool) {
	if path == "/" {
		return n.toNodeInfo(nil)
	}

	// 去除前缀的 /，并进行分割。例如 /a/b/c => [a, b, c]
	segments := strings.Split(strings.Trim(path, "/"), "/")
	cur := n
	params := Params{}
	for _, seg := range segments {
		// 禁止 /a//b 路由场景
		if seg == "" {
			return &nodeInfo{}, false
		}
		res, matchParam, ok := cur.child(seg)
		if !ok {
			return &nodeInfo{}, false
		}
		// 匹配到参数
		if matchParam {
			params = append(params, Param{key: res.path[1:], value: seg})
		}
		cur = res
	}

	return cur.toNodeInfo(params)
}

// child 搜索节点，不支持路由回溯（往回查找）
// 节点匹配优先级：静态路由 > 参数路由 > 通配符路由
func (n *node) child(path string) (*node, bool, bool) {
	// 先判断是否有静态路由
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}

	res, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return res, false, ok
}

// toNodeInfo helper 函数，将 node 转为 nodeInfo
//  1. 参数 params 中无数据，设置为 nil，方便 TDD 测试
//  2. node 中 handler 为 nil，无需进行转换，说明未绑定路由
//  3. 不符合如上两个场景，进行转换
func (n *node) toNodeInfo(params Params) (*nodeInfo, bool) {
	if len(params) == 0 {
		params = nil
	}

	if n.handler == nil {
		return &nodeInfo{}, false
	}

	return &nodeInfo{
		fullPath:    n.fullPath,
		params:      params,
		middlewares: n.middlewares,
		handler:     n.handler,
	}, true
}

type nodeInfo struct {
	fullPath    string
	middlewares []Middleware
	handler     HandleFunc
	params      Params
}
