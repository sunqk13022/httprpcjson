package httprpcjson

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var (
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfRequest = reflect.TypeOf((*http.Request)(nil)).Elem()
)

type service struct {
	name     string
	rcvr     reflect.Value
	rcvrType reflect.Type
	methods  map[string]*serviceMethod
}

type serviceMethod struct {
	method    reflect.Method
	argsType  reflect.Type
	replyType reflect.Type
}

type serviceMap struct {
	mutex    sync.Mutex
	services map[string]*service
}

func (m *serviceMap) register(rcvr interface{}, name string) error {
	s := &service{
		name:     name,
		rcvr:     reflect.ValueOf(rcvr),
		rcvrType: reflect.TypeOf(rcvr),
		methods:  make(map[string]*serviceMethod),
	}

	if name == "" {
		s.name = reflect.Indirect(s.rcvr).Type().Name()
		if !isExported(s.name) {
			return fmt.Errorf("rpc: type %q is not exported", s.name)
		}
	}
	if s.name == "" {
		return fmt.Errorf("rpc: no service name for type %q", s.rcvrType.String())
	}

	for i := 0; i < s.rcvrType.NumMethod(); i++ {
		method := s.rcvrType.Method(i)
		mtype := method.Type

		if method.PkgPath != "" {
			continue
		}

		if mtype.NumIn() != 4 {
			continue
		}

		reqType := mtype.In(1)
		if reqType.Kind() != reflect.Ptr || reqType.Elem() != typeOfRequest {
			continue
		}

		args := mtype.In(2)
		if args.Kind() != reflect.Ptr || !isExportedOrBuildin(args) {
			continue
		}

		reply := mtype.In(3)
		if reply.Kind() != reflect.Ptr || !isExportedOrBuildin(reply) {
			continue
		}

		if mtype.NumOut() != 1 {
			continue
		}

		if returnType := mtype.Out(0); returnType != typeOfError {
			continue
		}
		s.methods[method.Name] = &serviceMethod{
			method:    method,
			argsType:  args.Elem(),
			replyType: reply.Elem(),
		}
	}

	if len(s.methods) == 0 {
		return fmt.Errorf("rpc: %q has no exported methods of suitdable type", s.name)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.services == nil {
		m.services = make(map[string]*service)
	} else if _, ok := m.services[s.name]; ok {
		return fmt.Errorf("rpc:service already defined: %q", s.name)
	}
	m.services[s.name] = s
	return nil
}

func (m *serviceMap) get(method string) (*service, *serviceMethod, error) {
	parts := strings.Split(method, ".")
	if len(parts) != 2 {
		err := fmt.Errorf("rpc: service/method request ill-formed: %q", method)
		return nil, nil, err
	}

	m.mutex.Lock()
	service := m.services[parts[0]]
	m.mutex.Unlock()

	if service == nil {
		err := fmt.Errorf("rpc:cannot find service: %q", method)
		return nil, nil, err
	}

	serviceMethod := service.methods[parts[1]]
	if serviceMethod == nil {
		err := fmt.Errorf("rpc:cannot find method: %q", method)
		return nil, nil, err
	}
	return service, serviceMethod, nil
}

func isExported(name string) bool {
	rune1, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune1)
}

func isExportedOrBuildin(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return isExported(t.Name()) || t.PkgPath() == ""
}
