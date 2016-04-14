package backend

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// TrafficSnippet is a aggregate traffic record
type TrafficSnippet struct {
	ParentContext string
	Pool     string
	Rop      time.Time
	Src      string
	Context  string
	Key      string
	TrafficKpi TrafficKpi
}

type TrafficKpi struct {
	// RateRx is the bandwidth rate in download
	RateRx int
	// RateTx is the bandwidth rate in upload
	RateTx int
	// VolumeRx is the total traffic volume in download (MB)
	VolumeRx float32
	// VolumeTx is the total traffic volume in upload (MB)
	VolumeTx float32
	// SpeedRx is the mean speed in download (Mbps)
	SpeedRx float32
	// SpeedTx is the mean speed in upload (Mbps)
	SpeedTx float32
}



// TrafficMeasures is a collection of traffic records
type TrafficMeasures []TrafficSnippet

func dropTraffic() {
	mongo := GetMongoDB()
	coll := mongo.Session().DB("netdata").C("traffic")
	coll.DropCollection()
}

// AllTraffic return all ROPs
func AllTraffic(pool string, context string) TrafficMeasures {
	var measures = TrafficMeasures{}
	mongo := GetMongoDB()
	coll := mongo.Session().DB("netdata").C("traffic")
	fmt.Println("pool, context: ", pool, context)
	coll.Find(bson.M{"pool": pool, "context": context}).All(&measures)
	return measures
}

func (measures TrafficMeasures) save() {
	mongo := GetMongoDB()
	coll := mongo.Session().DB("netdata").C("traffic")
	err := coll.Insert(I(measures)...)
	if err != nil {
		log.Fatal(err)
	}

}

func (ts *TrafficSnippet) save() {
	mongo := GetMongoDB()
	coll := mongo.Session().DB("netdata").C("traffic")
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
		// Display record.
		// ... Display record length.
		// ... Display all individual elements of the slice.
		fmt.Println(record)
		fmt.Println(len(record))

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

		for value := range record {
			fmt.Printf("  %v\n", record[value])
		}
	}
}
