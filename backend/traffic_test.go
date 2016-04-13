package backend

import (
	"fmt"
	"path/filepath"
	"testing"
)

var dataDir = "fixtures"

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
