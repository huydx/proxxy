// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package requestLog

import (
	"time"
	"net/http"
	"encoding/gob"
	"bytes"
	"sync"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/huydx/proxxy/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
)

type RequestLogRecord struct {
	gorm.Model
	Content []byte
	UUid    string // uuid
	Ts      time.Time
}

type RequestDup struct {
	Method string
	URL    string
	Proto  string
	Header map[string][]string
	Body   []byte
}

var lock *sync.Mutex
var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "/tmp/requestlog.db")
	db.AutoMigrate(&RequestLogRecord{})
	gob.Register(http.NoBody)
	log.Fatal(err)
	lock = &sync.Mutex{}
}

func Write(r *http.Request) {
	dataBuff := bytes.NewBuffer(make([]byte, 0))
	encoder := gob.NewEncoder(dataBuff)
	err := encoder.Encode(copyRequest(r))
	log.Fatal(err)
	flush(dataBuff.Bytes())
}

func WriteAsync(r *http.Request) {
	go Write(r)
}

func copyRequest(r *http.Request) *RequestDup {
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return &RequestDup{
		Method: r.Method,
		URL:    r.URL.String(),
		Proto:  r.Proto,
		Header: r.Header,
		Body:   bodyBytes,
	}
}

func flush(bytes []byte) {
	id := uuid.New().URN()
	record := &RequestLogRecord{
		Ts:      time.Now(),
		Content: bytes,
		UUid:      id,
	}
	db2 := db.Create(record)
	if db2.Error != nil {
		log.Error(db2.Error)
	}
}
