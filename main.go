// Copyright (c) 2016  Jay R. Wren
// Licensed under the AGPLv3.

package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dustinkirkland/golang-petname"
	"github.com/juju/httprequest"
	"github.com/julienschmidt/httprouter"
)

func main() {
	html, err := ioutil.ReadFile("index.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New("f").Parse(string(html))
	if err != nil {
		log.Fatal(err)
	}
	f := func(p httprequest.Params) (*petnameHandler, error) {
		return &petnameHandler{t}, nil
	}
	router := httprouter.New()
	for _, h := range petnameErrorMapper.Handlers(f) {
		router.Handle(h.Method, h.Path, h.Handle)
	}

	log.Fatal(http.ListenAndServe(":2016", router))
}

type petnameHandler struct {
	t *template.Template
}

type petnamex struct{ Name string }

func (h *petnameHandler) Root(p httprequest.Params,
	arg *struct {
		httprequest.Route `httprequest:"GET /"`
	}) (petnamex, error) {
	name := petname.Generate(3, "-")
	a := newAccept(p.Request.Header["Accept"])
	r := petnamex{name}
	if a.Contains("text/html") {
		p.Response.Header()["Content-Type"] = []string{"FUCKYOU"}
		h.t.Execute(p.Response, r)
		return r, nil
	}
	return r, nil
}

type petnameErrorResponse struct{ Message string }

var petnameErrorMapper httprequest.ErrorMapper = func(err error) (int, interface{}) {
	return http.StatusInternalServerError, &petnameErrorResponse{
		Message: err.Error(),
	}
}

type accept struct{ a []string }

func newAccept(a []string) accept {
	joinedheaders := strings.Join(a, "")
	semis := strings.Split(joinedheaders, ";")
	aa := make([]string, 0, len(semis))
	for i := range semis {
		c := strings.Split(semis[i], ",")
		for j := range c {
			aa = append(aa, c[j])
		}
	}
	return accept{aa}
}
func (a accept) Contains(s string) bool {
	for i := range a.a {
		if a.a[i] == s {
			return true
		}
	}
	return false
}
