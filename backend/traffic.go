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
	// ROPPeriod is the observed period in minutes
	ROPPeriod = 60

	// DBName is the db name config key
	DBName = "DB"

	// ROPCollection is the collection name config key
	ROPCollection = "ropCollection"
)

// TrafficSnippet is a aggregate traffic record
type TrafficSnippet struct {
	Parentcontext string
	Pool          string
	Rop           time.Time
	Src           string
	Context       string
	Key           string
	Traffickpi    Traffickpi `bson:"trafficKpi"`
}

// Traffickpi is the traffic characterization
type Traffickpi struct {
	// RateRx is the bandwidth rate in download
	RateRx float32
	// RateTx is the bandwidth rate in upload
	RateTx float32
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
	ccc
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
		rateRxVal, _ := strconv.ParseFloat(record[4], 32)
		rateTxVal, _ := strconv.ParseFloat(record[5], 32)
		volumeRxVal, _ := strconv.ParseFloat(record[6], 32)
		volumeTxVal, _ := strconv.ParseFloat(record[7], 32)
		speedRxVal, _ := strconv.ParseFloat(record[8], 32)
		speedTxVal, _ := strconv.ParseFloat(record[9], 32)
		snippet := TrafficSnippet{
			Rop:     timeNow.Add(time.Duration(period*15) * time.Minute),
			Pool:    record[1],
			Context: record[2],
			Key:     record[3],
			Traffickpi: Traffickpi{
				RateRx:   float32(rateRxVal),
				RateTx:   float32(rateTxVal),
				VolumeRx: float32(volumeRxVal),
				VolumeTx: float32(volumeTxVal),
				SpeedRx:  float32(speedRxVal),
				SpeedTx:  float32(speedTxVal),
			},
		}
		snippet.save()

	}
}
