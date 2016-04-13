package backend

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"gopkg.in/mgo.v2/bson"
)

const (
	// DBName is the db name config key
	DBName = "DB"

	// ROPCollection is the collection name config key
	ROPCollection = "ropCollection"
)

// TrafficSnippet is a aggregate traffic record
type TrafficSnippet struct {
	Pool    string
	Rop     time.Time
	Src     string
	Context string
	Key     string
	Value   int
}

// TrafficMeasures is a collection of traffic records
type TrafficMeasures []TrafficSnippet

func dropTraffic() {
	cName := viper.Get(ROPCollection).(string)
	dbName := viper.Get(DBName).(string)
	mongo := GetMongoDB()
	coll := mongo.Session().DB(dbName).C(cName)
	coll.DropCollection()
}

// AllTraffic return all ROPs
func AllTraffic(pool string, context string) TrafficMeasures {
	var measures = TrafficMeasures{}
	mongo := GetMongoDB()
	cName := viper.Get(ROPCollection).(string)
	dbName := viper.Get(DBName).(string)
	coll := mongo.Session().DB(dbName).C(cName)
	coll.Find(bson.M{"pool": pool, "context": context}).All(&measures)
	return measures
}

func (measures TrafficMeasures) save() {
	mongo := GetMongoDB()
	cName := viper.Get(ROPCollection).(string)
	dbName := viper.Get(DBName).(string)
	coll := mongo.Session().DB(dbName).C(cName)
	err := coll.Insert(I(measures)...)
	if err != nil {
		log.Fatal(err)
	}

}

func (ts *TrafficSnippet) save() {
	mongo := GetMongoDB()
	cName := viper.Get(ROPCollection).(string)
	dbName := viper.Get(DBName).(string)
	coll := mongo.Session().DB(dbName).C(cName)
	err := coll.Insert(ts)
	if err != nil {
		log.Fatal(err)
	}
}

// CsvToDb read a csv file and store into mongo
func CsvToDb(filename string) {

	timeNow := time.Now()

	f, _ := os.Open(filename)
	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		period, _ := strconv.Atoi(record[0])
		val, _ := strconv.Atoi(record[4])
		snippet := TrafficSnippet{
			Rop:     timeNow.Add(time.Duration(period*15) * time.Minute),
			Pool:    record[1],
			Context: record[2],
			Key:     record[3],
			Value:   val,
		}
		snippet.save()

	}
}
