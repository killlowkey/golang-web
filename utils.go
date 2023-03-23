package web

import (
	"net/http"
	"path"
)

var (
	httpMethods = []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
)

type H map[string]any

// validHttpMethod 校验 HTTP 请求方法
// 正常的 HTTP 请求方法，返回 true，否则返回 false
func validHttpMethod(method string) bool {
	for _, m := range httpMethods {
		if m == method {
			return true
		}
	}

	return false
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
