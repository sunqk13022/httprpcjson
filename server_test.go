package httprpcjson

import (
	"net/http"
	"testing"
)

type ServiceServerTestRequest struct {
	A, B int
}

type ServiceServerTestResponse struct {
	Result int
}

type ServiceServer1 struct {
}

func (t *ServiceServer1) Mul(r *http.Request, req *ServiceServerTestRequest, res *ServiceServerTestResponse) error {
	res.Result = req.A * req.B
	return nil
}

type ServiceServer2 struct {
}

func TestRegisterService(t *testing.T) {
	var err error
	s := NewServer()
	service1 := new(ServiceServer1)
	service2 := new(ServiceServer2)

	err = s.RegisterService(service1, "")
	if err != nil || !s.HasMethod("ServiceServer1.Mul") {
		t.Errorf("Expected to be registered:ServiceServer1.Mul")
	}
	err = s.RegisterService(service2, "")
	if err == nil {
		t.Errorf("Expected error on service2")
	}
}
