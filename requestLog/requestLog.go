// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package requestLog

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/huydx/proxxy/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "/tmp/requestlog.db")
	db.AutoMigrate(&RequestLogRecord{})
	gob.Register(http.NoBody)
	log.Fatal(err)
}

func Log(r *http.Request) {
	dataBuff := bytes.NewBuffer(make([]byte, 0))
	encoder := gob.NewEncoder(dataBuff)
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
	dup := &RequestDup{
		Method: r.Method,
		URL:    r.URL.String(),
		Proto:  r.Proto,
		Header: r.Header,
		Body:   bodyBytesDup,
	}
	fmt.Println(string(dup.Body))
	err := encoder.Encode(dup)
	log.Fatal(err)
	go flush(dataBuff.Bytes())
}

func loadRequestLog(from time.Time, to time.Time) []*RequestLogRecord {
	rql := make([]*RequestLogRecord, 0)
	db.Where("ts <= ? and ts >= ?", to, from).Find(&rql)
	return rql
}

func LoadRequest(from time.Time, to time.Time) []*http.Request {
	rqs := make([]*http.Request, 0)
	rql := loadRequestLog(from, to)
	for _, rq := range rql {
		bs := rq.Content
		if bs == nil {
			// log
		} else {
			rqs = append(rqs, decode(bs))
		}
	}

	return rqs
}

func decode(bs []byte) *http.Request {
	dcd := gob.NewDecoder(bytes.NewBuffer(bs))
	dup := & RequestDup{}
	dcd.Decode(dup)
	req, err := http.NewRequest(
		dup.Method,
		dup.URL,
		bytes.NewBuffer(dup.Body),
	)
	if err != nil {

	}
	return req
}

func flush(bytes []byte) {
	id := uuid.New().URN()
	record := &RequestLogRecord{
		Ts:      time.Now(),
		Content: bytes,
		UUid:    id,
	}
	db2 := db.Create(record)
	if db2.Error != nil {
		log.Error(db2.Error)
	}
}
