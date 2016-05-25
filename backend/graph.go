package backend

import (
	"reflect"
	"time"
)

/*
* Graph data grid creator
 */
type Row []interface{}

type Grid []Row

type Grids struct {
	WeightedAvgRateRx Grid
	WeightedAvgRateTx Grid
	SumVolumeRx       Grid
	SumVolumeTx       Grid
	RateRx            Grid
	RateTx            Grid
	VolumeRx          Grid
	VolumeTx          Grid
	SpeedRx           Grid
	SpeedTx           Grid
}

func GetGraphData(tmf TrafficMeasureFilter) Grids {

	// Retrieve data from DB
	measures := GetTrafficMeasures(tmf)
	keys := getKeys(measures)

	// Create graph grids ....
	grids := Grids{}

	grids.SumVolumeRx = createSumGrid(keys, measures, "VolumeRx")
	grids.SumVolumeTx = createSumGrid(keys, measures, "VolumeTx")

	grids.RateRx = createRopGrid(keys, measures, "RateRx")
	grids.RateTx = createRopGrid(keys, measures, "RateTx")
	grids.VolumeRx = createRopGrid(keys, measures, "VolumeRx")
	grids.VolumeTx = createRopGrid(keys, measures, "VolumeTx")
	grids.SpeedRx = createRopGrid(keys, measures, "SpeedRx")
	grids.SpeedTx = createRopGrid(keys, measures, "SpeedTx")

	return grids
}

func getKeys(measures TrafficMeasures) map[string]int {
	result := make(map[string]int)
	var i int = 0
	for _, m := range measures {
		if _, ok := result[m.Key]; !ok {
			result[m.Key] = i
			i++
		}
	}
	return result
}

func createSumGrid(keys map[string]int, measures TrafficMeasures, measureName string) Grid {

	rowIndex := 0
	grid := make(Grid, len(keys)+1)
	grid[rowIndex] = Row{"context", measureName}

	for kName, kIndex := range keys {
		grid[kIndex+1] = make(Row, 2)
		grid[kIndex+1][0] = kName
	}

	for _, m := range measures {
		if i, ok := keys[m.Key]; ok {
			r := reflect.ValueOf(&m.TrafficKpi)
			f := reflect.Indirect(r).FieldByName(measureName)
			if grid[i+1][1] != nil {
				grid[i+1][1] = grid[i+1][1].(float64) + f.Float()
			} else {
				grid[i+1][1] = f.Float()
			}
		}
	}

	return grid
}

func createRopGrid(keys map[string]int, measures TrafficMeasures, measureName string) Grid {
	rowIndex := 0
	grid := make(Grid, 1)
	grid[rowIndex] = make(Row, len(keys)+1)
	grid[rowIndex][0] = "rop"
	for kName, kIndex := range keys {
		grid[rowIndex][kIndex+1] = kName
	}

	// Fill data grid
	lastRop := time.Time{}
	for _, m := range measures {
		if !lastRop.Equal(m.Rop) {
			lastRop = m.Rop
			rowIndex += 1
			grid = append(grid, make(Row, len(keys)+1))
			grid[rowIndex][0] = m.Rop.Format("15:04")
		}
		if i, ok := keys[m.Key]; ok {
			r := reflect.ValueOf(&m.TrafficKpi)
			f := reflect.Indirect(r).FieldByName(measureName)
			grid[rowIndex][i+1] = f.Float()
		}
	}

	return grid
}
