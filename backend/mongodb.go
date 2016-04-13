package backend

import (
	"log"
	"reflect"
	"sync"

	"gopkg.in/mgo.v2"
)

var (
	// URL is the default connection endpoint
	URL = "mongodb://127.0.0.1:27017/?connect=direct"
)

type (

	// MongoDB interface is simply a Session getter
	MongoDB interface {
		Session() *mgo.Session
	}

	mongoDB struct {
		session *mgo.Session
	}
)

// I converts []struct{} to []interface{}
func I(array interface{}) []interface{} {

	v := reflect.ValueOf(array)
	t := v.Type()

	if t.Kind() != reflect.Slice {
		log.Panicf("`array` should be %s but got %s", reflect.Slice, t.Kind())
	}

	result := make([]interface{}, v.Len(), v.Len())

	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result
}

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

// GetMongoDB singleton
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
