// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/foo", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "ok")
	})
	http.ListenAndServe(":9999", nil)
}
