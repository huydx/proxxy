// Copyright (c) 2016 LINE Corporation. All rights reserved.
// LINE Corporation PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.

package analysys

import (
	"github.com/huydx/proxxy/requestLog"
	"github.com/jinzhu/gorm"
	"github.com/huydx/proxxy/log"
)

var db *gorm.DB

// Measure for each URL
type Measure struct {
	Url         string
	Count       int
	Total       float64
	Mean        float64
	Stddev      float64
	Min         float64
	Percentiles []float64
	Max         float64
	S2xx        int
	S3xx        int
	S4xx        int
	S5xx        int
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "/tmp/requestlog.db")
	log.Fatal(err)
}

func MaxLatency() *requestLog.RequestDupJson {
	rql := requestLog.RequestDup{}
	if err := db.Order("time_nano DESC").First(&rql).Error; err != nil {
		log.Error(err)
		return nil
	} else {
		return rql.ToJSON()
	}
}

func MinLatency() *requestLog.RequestDupJson {
	rql := requestLog.RequestDup{}
	if err := db.Order("time_nano ASC").First(&rql).Error; err != nil{
		log.Error(err)
		return nil
	} else {
		return rql.ToJSON()
	}
}

func SortByLatency(limit int, page int) []*requestLog.RequestDupJson {
	rql := make([]*requestLog.RequestDup, 0)
	if err := db.Order("time_nano DESC").Offset(limit * page).Find(&rql); err != nil {
		log.Error(err)
		return nil
	} else {
		rqlj := make([]*requestLog.RequestDupJson, len(rql))
		for _, rq := range rql {
			rqlj = append(rqlj, rq.ToJSON())
		}
		return rqlj
	}
}

func Measures(limit int, page int) []*Measure {
	return nil
}
