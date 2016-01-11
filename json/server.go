package json

import (
	"encoding/json"
	"github.com/sunqk13022/httprpcjson"
	"net/http"
)

var null = json.RawMessage([]byte("null"))
var Version = "2.0"

type serverRequest struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
	Id      *json.RawMessage `json:"id"`
}

type serverResponse struct {
	Version string           `json:"jsonrpc"`
	Result  interface{}      `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
	Id      *json.RawMessage `json:"id"`
}

func NewCodec() *Codec {
	return &Codec{}
}

type Codec struct {
}

func (c *Codec) NewRequest(r *http.Request) httprpcjson.CodecRequest {
	return newCodecRequest(r)
}

func newCodecRequest(r *http.Request) httprpcjson.CodecRequest {
	req := new(serverRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		err = &Error{
			Code:    1,
			Message: err.Error(),
			Data:    req,
		}
	}

	if req.Version != Version {
		err = &Error{
			Code:    1,
			Message: "jsonrpc must be " + Version,
			Data:    req,
		}
	}
	r.Body.Close()
	return &CodecRequest{request: req, err: err}
}

type CodecRequest struct {
	request *serverRequest
	err     error
}

func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.request.Method, nil
	}
	return "", c.err
}

func (c *CodecRequest) ReadRequest(args interface{}) error {
	if c.err == nil && c.request.Params != nil {
		err := json.Unmarshal(*c.request.Params, args)
		if err != nil {
			c.err = &Error{
				Code:    1,
				Message: err.Error(),
				Data:    c.request.Params,
			}
		}
	}
	return c.err
}

func (c *CodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}) {
	res := &serverResponse{
		Version: Version,
		Result:  reply,
		Id:      c.request.Id,
	}
	c.writeServerResponse(w, res)
}

func (c *CodecRequest) WriteError(w http.ResponseWriter, status int, err error) {
	jsonErr, ok := err.(*Error)
	if !ok {
		jsonErr = &Error{
			Code:    1,
			Message: err.Error(),
		}
	}

	res := &serverResponse{
		Version: Version,
		Error:   jsonErr,
		Id:      c.request.Id,
	}
	c.writeServerResponse(w, res)
}

func (c *CodecRequest) writeServerResponse(w http.ResponseWriter, res *serverResponse) {
	if c.request.Id != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		err := encoder.Encode(res)

		if err != nil {
			httprpcjson.WriteError(w, 400, err.Error())
		}
	}
}
