package binding

import "net/http"

var (
	JSON     = jsonBinding{}
	XML      = xmlBinding{}
	QUERY    = queryBinding{}
	PROTOBUF = protobufBinding{}
)

// Binding bind request's data to any interface
type Binding interface {
	Name() string
	Bind(*http.Request, any) error
}

// BindingBody bind request body
type BindingBody interface {
	Binding
	BindBody([]byte, any) error
}
