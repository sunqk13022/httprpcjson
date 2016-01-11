package httprpcjson

import (
	"net/http"
	"testing"
)

type MapTestRequest struct {
	A, B int
}

type MapTestResponse struct {
	Result int
}

type MapTestService1 struct {
}

func (t *MapTestService1) Mul(r *http.Request, req *MapTestRequest, res *MapTestResponse) error {
	return nil
}

func TestRegiter(t *testing.T) {
	serviceMap1 := &serviceMap{}
	serviceMap1.register(new(MapTestService1), "")
	if len(serviceMap1.services) != 1 {
		t.Error("Expected len=1, but got:", len(serviceMap1.services))
	}

	if _, ok := serviceMap1.services["MapTestService1"]; !ok {
		t.Error("Expected true, but got:", ok)
	}
}
