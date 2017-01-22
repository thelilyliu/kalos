package main

import (
	// "log"
	"math/rand"
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
	Rating float64 `json:"rating"`
}

/*
  ========================================
  Load Polls
  ========================================
*/

func loadPollsDB(polls *[]Poll) error {
	// create new MongoDB session
	collection, session := initMongoDB("poll")
	defer session.Close()

	selector := bson.M{"user": user}

	return collection.Find(selector).Sort("-time").All(polls)
}

/*
  ========================================
  Load Poll
  ========================================
*/

func loadPollDB(poll *Poll) error {
	// create new MongoDB session
	collection, session := initMongoDB("poll")
	defer session.Close()

	selector := bson.M{"id": poll.ID}

	return collection.Find(selector).One(poll)
}

/*
  ========================================
  Insert Poll
  ========================================
*/

func insertPollDB(poll *Poll) error {
	// create new MongoDB session
	collection, session := initMongoDB("poll")
	defer session.Close()

	timeEST := time.Now().Add(-4 * time.Hour)
	poll.Time = timeEST.Format("20060102150405")
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	poll.User = user
	poll.ID = bson.NewObjectId().String()
	poll.ID = poll.ID[13 : len(poll.ID)-2]
	poll.Code = 1000 + r1.Intn(9000)
	poll.Options = make([]string, 1)

	return collection.Insert(poll)
}

/*
  ========================================
  Update Poll
  ========================================
*/

func updatePollDB(poll *Poll) error {
	// create new MongoDB session
	collection, session := initMongoDB("poll")
	defer session.Close()

	selector := bson.M{"id": poll.ID}
	change := bson.M{"question": poll.Question, "options": poll.Options}
	update := bson.M{"$set": &change}

	return collection.Update(selector, update)
}

/*
  ========================================
  Delete Poll
  ========================================
*/

func deletePollDB(pollID string) error {
	// create new MongoDB session
	collection, session := initMongoDB("poll")
	defer session.Close()

	selector := bson.M{"id": pollID}

	return collection.Remove(selector)
}

/*
  ========================================
  Submit Code
  ========================================
*/

func submitCodeDB(poll *Poll) error {
	// create new MongoDB session
	collection, session := initMongoDB("poll")
	defer session.Close()

	selector := bson.M{"code": poll.Code}

	return collection.Find(selector).One(poll)
}

/*
  ========================================
  Submit Response
  ========================================
*/

func submitResponseDB(pollID string, response *Response) error {
	// create new MongoDB session
	collection, session := initMongoDB("poll")
	defer session.Close()

	selector := bson.M{"id": pollID}
	update := bson.M{"$push": bson.M{"responses": response}}

	return collection.Update(selector, update)
}
