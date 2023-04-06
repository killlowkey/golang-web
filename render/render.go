package render

import "net/http"

type Render interface {
	Render(http.ResponseWriter) error
	WriteContentType(w http.ResponseWriter)
}

// 校验是否符合该接口
var (
	_ Render = JSON{}
	_ Render = String{}
	_ Render = HTML{}
	_ Render = ProtoBuf{}
)

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
