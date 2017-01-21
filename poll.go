package main

import (
	// "log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Poll struct {
	User      string     `json:"user"`
	ID        string     `json:"id"`
	Time      string     `json:"time"`
	Code      int        `json:"code"`
	Question  string     `json:"question"`
	Options   []string   `json:"options"`
	Responses []Response `json:"responses"`
	Results   []Result   `json:"results"`
}

type Response struct {
	Name    string    `json:"name"`
	Ratings []float64 `json:"ratings"`
}

type Result struct {
	Option string  `json:"option"`
	Score  float64 `json:"score"`
}
