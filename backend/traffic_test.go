package backend

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

var dataDir = "fixtures"

func setUp() {
	viper.SetDefault(ROPCollection, "traffic")
	viper.SetDefault(DBName, "netdata")
}

func tearDown() {

}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func TestCleanDb(t *testing.T) {
	dropTraffic()
}

var r = rand.New(rand.NewSource(99))

func randomNormalizedArray(size int) []float64 {
	values := make([]float64, size)

	for i := range values {
		values[i] = r.Float64()
	}
	sum := 0.0
	for _, val := range values {
		sum += val
	}
	for i := range values {
		values[i] /= sum
	}

	return values
}

func randomArray(size int, maxValue float64) []float64 {
	values := make([]float64, size)

	for i := range values {
		values[i] = r.Float64() * maxValue
	}

	return values
}

func createContextTraffic(rop time.Time, pool string, context string, items []string) {
	size := len(items)

	//imap := make(map[string]TrafficSnippet, len)

	rateRx := randomNormalizedArray(size)
	rateTx := randomNormalizedArray(size)
	volumeRx := randomArray(size, 100)
	volumeTx := randomArray(size, 100)
	speedRx := randomArray(size, 100)
	speedTx := randomArray(size, 100)

	token := strings.Split(context, ".")

	shortContext := token[len(token)-1]
	parentContext := ""
	if len(token) > 1 {
		parentContext = token[len(token)-2]
	}

	for i, key := range items {
		ts := TrafficSnippet{
			ParentContext: parentContext,
			Rop:           rop,
			Pool:          pool,
			Context:       shortContext,
			Key:           key,
			TrafficKpi: TrafficKpi{
				RateRx:   float32(rateRx[i]),
				RateTx:   float32(rateTx[i]),
				VolumeRx: float32(volumeRx[i]),
				VolumeTx: float32(volumeTx[i]),
				SpeedRx:  float32(speedRx[i]),
				SpeedTx:  float32(speedTx[i]),
			},
		}
		ts.save()

	}

}

func TestInsertMeasures(t *testing.T) {

	timeZero := time.Date(2016, time.April, 14, 0, 0, 0, 0, time.UTC)

	numRop := 12
	lines := []string{"linea1"}
	services := []string{"http", "mail", "bittorrent"}
	servers := []string{"google", "facebook", "youtube", "noiportal.it", "others"}

	for i := 0; i < numRop; i++ {
		t := timeZero.Add(time.Duration(i*ROPPeriod) * time.Minute)

		for _, line := range lines {
			createContextTraffic(t, line, "net", services)
			createContextTraffic(t, line, "net.http", servers)
		}
	}

	for _, server := range servers {
		percentuale := randomArray(5, 200)
		fmt.Println(server, percentuale)
	}

}

// func TestLoadFromCsv(t *testing.T) {
// 	filename := filepath.Join(dataDir, "web.csv")
// 	fmt.Println("loading from: ", filename)
// 	CsvToDb(filename)
// }
