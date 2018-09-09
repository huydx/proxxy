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

func (rql *RequestDup) String() string {
	return fmt.Sprintf("Method: %s; Header: %v; Body: %s; TimeNano: %d",
		rql.Method, rql.Header, string(rql.Body), rql.TimeNano)
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
