// Copyright (c) 2016  Jay R. Wren
// Licensed under the AGPLv3.

package main

import (
	"log"
	"net/http"

	"github.com/dustinkirkland/golang-petname"
	"github.com/julienschmidt/httprouter"

	"github.com/juju/httprequest"
)

func main() {
	f := func(p httprequest.Params) (petnameHandler, error) {
		return petnameHandler{}, nil
	}
	router := httprouter.New()
	for _, h := range petnameErrorMapper.Handlers(f) {
		router.Handle(h.Method, h.Path, h.Handle)
	}

	log.Fatal(http.ListenAndServe(":2016", router))
}

type petnameHandler struct{}

type petnamex struct{ Name string }

func (petnameHandler) Root(arg *struct {
	httprequest.Route `httprequest:"GET /"`
}) (petnamex, error) {
	name := petname.Generate(3, "-")
	return petnamex{name}, nil
}

type petnameErrorResponse struct{ Message string }

var petnameErrorMapper httprequest.ErrorMapper = func(err error) (int, interface{}) {
	return http.StatusInternalServerError, &petnameErrorResponse{
		Message: err.Error(),
	}
}
