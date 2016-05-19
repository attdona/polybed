package backend

import (
	"fmt"
	"reflect"
)

/*
* Graph data grid creator
 */
type PlotlyLabelsValuesData struct {
	Labels []string
	Values []interface{}
}

type PlotlyXYData struct {
	Name string
	X    []interface{}
	Y    []interface{}
}

type PlotlyGrids struct {
	WeightedAvgRateRx PlotlyLabelsValuesData
	WeightedAvgRateTx PlotlyLabelsValuesData
	SumVolumeRx       PlotlyLabelsValuesData
	SumVolumeTx       PlotlyLabelsValuesData
	RateRx            []PlotlyXYData
	RateTx            []PlotlyXYData
	VolumeRx          []PlotlyXYData
	VolumeTx          []PlotlyXYData
	SpeedRx           []PlotlyXYData
	SpeedTx           []PlotlyXYData
}

func GetPlotlyGraphData(tmf TrafficMeasureFilter) PlotlyGrids {

	// Retrieve data from DB
	measures := GetTrafficMeasures(tmf)
	keys := getKeys(measures)

	// Create graph grids ....
	plotlyGrids := PlotlyGrids{}

	plotlyGrids.WeightedAvgRateRx = createPlotlyWeightedAvgGrid(keys, measures, "RateRx", "VolumeRx")
	plotlyGrids.WeightedAvgRateTx = createPlotlyWeightedAvgGrid(keys, measures, "RateTx", "VolumeTx")

	plotlyGrids.SumVolumeRx = createPlotlySumGrid(keys, measures, "VolumeRx")
	plotlyGrids.SumVolumeTx = createPlotlySumGrid(keys, measures, "VolumeTx")

	plotlyGrids.RateRx = createPlotlyRopGrid(keys, measures, "RateRx")
	plotlyGrids.RateTx = createPlotlyRopGrid(keys, measures, "RateTx")
	plotlyGrids.VolumeRx = createPlotlyRopGrid(keys, measures, "VolumeRx")
	plotlyGrids.VolumeTx = createPlotlyRopGrid(keys, measures, "VolumeTx")
	plotlyGrids.SpeedRx = createPlotlyRopGrid(keys, measures, "SpeedRx")
	plotlyGrids.SpeedTx = createPlotlyRopGrid(keys, measures, "SpeedTx")

	return plotlyGrids
}

func createPlotlyWeightedAvgGrid(keys map[string]int, measures TrafficMeasures, measureName string, measureWeightName string) PlotlyLabelsValuesData {

	labels := make([]string, len(keys))
	values := make([]interface{}, len(keys))
	sum := make([]interface{}, len(keys))

	for kName, kIndex := range keys {
		labels[kIndex] = kName
	}

	for _, m := range measures {
		if i, ok := keys[m.Key]; ok {
			r := reflect.ValueOf(&m.TrafficKpi)
			f := reflect.Indirect(r).FieldByName(measureName)
			fw := reflect.Indirect(r).FieldByName(measureWeightName)
			if values[i] != nil {
				values[i] = values[i].(float64) + (f.Float() * fw.Float())
				sum[i] = sum[i].(float64) + fw.Float()
			} else {
				values[i] = f.Float() * fw.Float()
				sum[i] = fw.Float()
			}
		}
	}

	for i, _ := range values {
		if values[i] != nil {
			values[i] = values[i].(float64) / sum[i].(float64)
		}
	}

	result := PlotlyLabelsValuesData{labels, values}
	fmt.Println(result)

	return result
}

func createPlotlySumGrid(keys map[string]int, measures TrafficMeasures, measureName string) PlotlyLabelsValuesData {

	labels := make([]string, len(keys))
	values := make([]interface{}, len(keys))

	for kName, kIndex := range keys {
		labels[kIndex] = kName
	}

	for _, m := range measures {
		if i, ok := keys[m.Key]; ok {
			r := reflect.ValueOf(&m.TrafficKpi)
			f := reflect.Indirect(r).FieldByName(measureName)
			if values[i] != nil {
				values[i] = values[i].(float64) + f.Float()
			} else {
				values[i] = f.Float()
			}
		}
	}

	result := PlotlyLabelsValuesData{labels, values}
	fmt.Println(result)

	return result

}

func createPlotlyRopGrid(keys map[string]int, measures TrafficMeasures, measureName string) []PlotlyXYData {

	result := make([]PlotlyXYData, len(keys))
	for kName, kIndex := range keys {
		result[kIndex] = PlotlyXYData{kName, make([]interface{}, 0), make([]interface{}, 0)}
	}

	// Fill data grid
	for _, m := range measures {
		if i, ok := keys[m.Key]; ok {
			r := reflect.ValueOf(&m.TrafficKpi)
			f := reflect.Indirect(r).FieldByName(measureName)
			result[i].X = append(result[i].X, m.Rop.Format("15:04"))
			result[i].Y = append(result[i].Y, f.Float())
		}
	}

	return result
}
