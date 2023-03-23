package binding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonBinding struct{}

func (j jsonBinding) Name() string {
	return "json"
}

func (j jsonBinding) Bind(request *http.Request, a any) error {
	if request.Body == nil || a == nil {
		return errors.New("invalid data")
	}
	return decodeJSON(request.Body, a)
}

func (jsonBinding) BindBody(body []byte, obj any) error {
	return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
