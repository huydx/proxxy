// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package requestLog

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/huydx/proxxy/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type RequestDup struct {
	gorm.Model
	UUid     string // uuid
	Method   string
	URL      string
	Proto    string
	Header   []byte
	Body     []byte
	TimeNano int64
	Ts       time.Time
}

type RequestDupJson struct {
	UUid     string            `json:"u_uid"`
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Proto    string            `json:"proto"`
	Header   map[string]string `json:"header"`
	Body     string            `json:"body"`
	TimeNano int64             `json:"time_nano"`
	Ts       time.Time         `json:"time_stampt"`
}

func (dup *RequestDup) ToJSON() *RequestDupJson {
	return &RequestDupJson{
		UUid:     dup.UUid,
		Method:   dup.Method,
		URL:      dup.URL,
		Proto:    dup.Proto,
		Header:   decoderHeader(dup.Header),
		Body:     string(dup.Body),
		TimeNano: dup.TimeNano,
		Ts:       dup.Ts,
	}
}

func decoderHeader(headerBytes []byte) map[string]string {
	res := make(map[string]string)
	dc := gob.NewDecoder(bytes.NewBuffer(headerBytes))
	dc.Decode(&res)
	return res
}

func (rql *RequestDup) String() string {
	return fmt.Sprintf("Method: %s; Header: %v; Body: %s; TimeMs: %f",
		rql.Method, decoderHeader(rql.Header), string(rql.Body), float64(rql.TimeNano)/1000000.0)
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "/tmp/requestlog.db")
	db.AutoMigrate(&RequestDup{})
	gob.Register(http.NoBody)
	log.Fatal(err)
}

func Log(dup *RequestDup) {
	go func() {
		id := uuid.New().URN()
		dup.UUid = id
		db2 := db.Create(dup)
		if db2.Error != nil {
			log.Error(db2.Error)
		}
	}()
}

func LoadRequestDup(from time.Time, to time.Time) []*RequestDup {
	rql := make([]*RequestDup, 0)
	db.Where("ts <= ? and ts >= ?", to, from).Find(&rql)
	fmt.Println(len(rql))
	return rql
}

func LoadRequest(from time.Time, to time.Time) []*http.Request {
	rqs := make([]*http.Request, 0)
	rql := LoadRequestDup(from, to)
	for _, dup := range rql {
		bs := dup.Body
		if bs == nil {
			// log
		} else {
			req, err := http.NewRequest(
				dup.Method,
				dup.URL,
				bytes.NewBuffer(dup.Body),
			)
			decoder := gob.NewDecoder(bytes.NewBuffer(dup.Header))
			var hd = make(map[string]string)
			decoder.Decode(hd)
			for k, v := range hd {
				req.Header.Add(k, v)
			}

			if err != nil {
				log.Error(err)
			} else {
				rqs = append(rqs, req)
			}
		}
	}
	return rqs
}
