// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package main

import (
	"net/url"
	"net/http"

	"github.com/huydx/proxxy/proxy"
	"github.com/huydx/proxxy/log"
	"github.com/huydx/proxxy/requestLog"
)

func main() {
	u, err := url.Parse("http://127.0.0.1:9090")
	log.Fatal(err)
	rvp := proxy.NewSingleHostReverseProxy(u)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got request %v", r)
		requestLog.WriteAsync(r)
		rvp.ServeHTTP(w, r)
	})
	http.ListenAndServe(":8080", nil)
}
