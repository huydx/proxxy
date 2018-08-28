// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func init() {
	log.SetOutput(os.Stdout)
}

func Fatal(v ...interface{}) {
	if len(v) > 0 && v[0] != nil {
		log.Printf("[FATAL] %s %v", fileLine(), v)
		os.Exit(1)
	}
}

func Error(v ...interface{}) {
	log.Printf("[ERROR] %s %v", fileLine(), v)
}

func Warn(v ...interface{}) {
	log.Printf("[WARN] %s %v", fileLine(), v)
}

func Printf(s string, v ...interface{}) {
	log.Printf(fileLine()+" "+s, v)
}

func fileLine() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
