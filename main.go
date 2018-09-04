// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/huydx/proxxy/log"
	"github.com/huydx/proxxy/proxy"
	"github.com/huydx/proxxy/requestLog"
)

type proxyHandler struct {
	rvp *proxy.ReverseProxy
}

func (m *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestLog.Log(r)
	m.rvp.ServeHTTP(w, r)
}

type staticHandler struct {
	fs http.Handler
}

func (m *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.fs.ServeHTTP(w, r)
}

func main() {
	u, err := url.Parse("http://127.0.0.1:9999")
	log.Fatal(err)
	rvp := proxy.NewSingleHostReverseProxy(u)
	router := mux.NewRouter()
	router.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()
		from := n.AddDate(-1, 0, 0)
		to := n
		rqs := requestLog.LoadRequest(from, to)
		for _, r := range rqs {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(body))
			}
		}
	}).Methods("GET")

	go http.ListenAndServe(":8080", &proxyHandler{rvp: rvp})
	go http.ListenAndServe(":8081", router)
	http.ListenAndServe(":8082", &staticHandler{fs: http.FileServer(http.Dir("view"))})
}
