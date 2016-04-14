package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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

func TestInsertMeasures(t *testing.T) {
	measures := TrafficMeasures{
		TrafficSnippet{Pool: "myPool", Context: "net", Key: "http", Value: 40},
		TrafficSnippet{Pool: "myPool", Context: "net", Key: "mail", Value: 60},
	}
	measures.save()
}

func TestInsertSnippet(t *testing.T) {

	snippet := TrafficSnippet{
		Pool: "myPool",
	}

	snippet.save()

	if snippet.Pool != "myPool" {
		t.Fatalf("my first test has failed")
	}

}

func TestLoadFromCsv(t *testing.T) {
	filename := filepath.Join(dataDir, "web.csv")
	fmt.Println("loading from: ", filename)
	CsvToDb(filename)
}
