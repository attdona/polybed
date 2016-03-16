package db

import (
	"log"
	"sync"

	"gopkg.in/mgo.v2"
)

var (
	URL = "mongodb://127.0.0.1:27017/?connect=direct"
)

type (
	MongoDB interface {
		Session() *mgo.Session
	}

	mongoDB struct {
		session *mgo.Session
	}
)

func newMongoDB() *mongoDB {
	log.Printf("Connecting to Mongo: %s ...", URL)
	s, e := mgo.Dial(URL)
	if e != nil {
		panic(e)
	}

	log.Printf("Successfully connected to Mongo")
	return &mongoDB{s}
}

func (db *mongoDB) Session() *mgo.Session {
	return db.session.Copy()
}

var (
	mdb  *mongoDB
	once sync.Once
)

func GetMongoDB() MongoDB {
	once.Do(func() {
		defer func() {
			if x := recover(); x != nil {
				log.Fatalf("run time panic: %v", x)
			}
		}()

		mdb = newMongoDB()
	})
	return mdb
}
