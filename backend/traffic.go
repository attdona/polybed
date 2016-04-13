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
	mongo := GetMongoDB()
	coll := mongo.Session().DB("netdata").C("traffic")
	coll.DropCollection()
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
