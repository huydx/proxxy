// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/huydx/proxxy/log"
	"github.com/huydx/proxxy/proxy"
	"github.com/huydx/proxxy/requestLog"
)

type proxyHandler struct {
	rvp *proxy.ReverseProxy
}

func (m *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bodyBytesDup []byte
	if r.Body != nil {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		bodyBytesDup = make([]byte, len(bodyBytes))
		copy(bodyBytesDup, bodyBytes)
		if err != nil {
			log.Error(err)
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	t := m.rvp.ServeHTTP(w, r)
	headerBytes := bytes.NewBuffer(make([]byte, 0))
	enc := gob.NewEncoder(headerBytes)
	err := enc.Encode(r.Header)
	if err != nil {
		log.Error(err)
	}
	dup := &requestLog.RequestDup{
		Method:   r.Method,
		URL:      r.URL.String(),
		Proto:    r.Proto,
		Header:   headerBytes.Bytes(),
		Body:     bodyBytesDup,
		TimeNano: t,
	}
	requestLog.Log(dup)
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
		rqs := requestLog.LoadRequestDup(from, to)
		for _, r := range rqs {
			fmt.Println(r)
		}
	}).Methods("GET")

	go http.ListenAndServe(":8080", &proxyHandler{rvp: rvp})
	go http.ListenAndServe(":8081", router)
	http.ListenAndServe(":8082", &staticHandler{fs: http.FileServer(http.Dir("view"))})
}
