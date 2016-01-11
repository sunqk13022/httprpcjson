package main

import (
	"fmt"
	"github.com/sunqk13022/httprpcjson"
	"github.com/sunqk13022/httprpcjson/json"
	"log"
	"net/http"
)

type Counter struct {
	Count int
}

type GetReq struct {
	A, B int
}

func (c *Counter) Get(r *http.Request, req *GetReq, res *Counter) error {
	log.Printf("<- Get %+v", *req)
	*res = *c
	log.Printf("-> %v", *res)
	return nil
}

func main() {
	s := httprpcjson.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(Counter), "")
	http.Handle("/jsonrpc/", s)
	fmt.Println("listen 1234 port")
	http.ListenAndServe(":1234", nil)
}
