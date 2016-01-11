package json

import (
	"bytes"
	"errors"
	"github.com/sunqk13022/httprpcjson"
	"net/http"
	"testing"
)

type ResponseRecorder struct {
	Code      int
	HeaderMap http.Header
	Body      *bytes.Buffer
	Flushed   bool
}

func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
	}
}

const DefaultRemoteAddr = "1.2.3.4"

func (rw *ResponseRecorder) Header() http.Header {
	return rw.HeaderMap
}

func (rw *ResponseRecorder) Write(buf []byte) (int, error) {
	if rw.Body != nil {
		rw.Body.Write(buf)
	}

	if rw.Code == 0 {
		rw.Code = http.StatusOK
	}
	return len(buf), nil
}

func (rw *ResponseRecorder) WriteHeader(code int) {
	rw.Code = code
}

func (rw *ResponseRecorder) Flush() {
	rw.Flushed = true
}

var ErrResponseError = errors.New("response error")

type Service1Request struct {
	A, B int
}

type Service1Response struct {
	Result int
}

type Service1 struct {
}

func (t *Service1) Mul(r *http.Request, req *Service1Request, res *Service1Response) error {
	res.Result = req.A * req.B
	return nil
}

func (t *Service1) ResponseError(r *http.Request, req *Service1Request, res *Service1Response) error {
	return ErrResponseError
}

func execute(t *testing.T, s *httprpcjson.Server, method string, req interface{}, res interface{}) error {
	if !s.HasMethod(method) {
		t.Fatal("Expected to be registered:", method)
	}

	buf, _ := EncodeClientRequest(method, req)
	body := bytes.NewBuffer(buf)

	r, _ := http.NewRequest("POST", "http://localhost:8080/", body)
	r.Header.Set("Content-Type", "application/json")

	w := NewRecorder()
	s.ServeHTTP(w, r)

	return DecodeClientResponse(w.Body, res)
}

func TestService(t *testing.T) {
	s := httprpcjson.NewServer()
	s.RegisterCodec(NewCodec(), "application/json")
	s.RegisterService(new(Service1), "")

	var res Service1Response
	if err := execute(t, s, "Service1.Mul", &Service1Request{4, 2}, &res); err != nil {
		t.Errorf("Expected to get nil,but got", err)
	}

	if res.Result != 8 {
		t.Errorf("Wrong response:%v.", res.Result)
	}

	if err := execute(t, s, "Service1.ResponseError", &Service1Request{4, 2}, &res); err == nil {
		t.Errorf("Expected to get %q,but got nil", ErrResponseError)
	} else if err.Error() != ErrResponseError.Error() {
		t.Errorf("Expected to get %q,but got %q", ErrResponseError, err)
	}
}
