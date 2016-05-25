
function createSumLabelsValues(keys, rawData, measurename) {

  var result = {
     Labels: [],
     Values: [],
  }

  for (let ts of rawData) {
     let i = keys[ts.key]
     if (!result.Labels[i]) {
        result.Labels[i] = ts.key
        result.Values[i] = 0
     }
     result.Values[i] += ts.trafficKpi[measurename]
  }

  return result
}

function extractKeys(rawData) {
   let keys = {}
   let i = 0
   for (let ts of rawData) {
      if (!(ts.key in keys)) {
         keys[ts.key] = i++
      }
   }
   return keys
}

function createLinearXY(keys, rawData, name) {
  var measures = []
  for (let key in keys) {
     measures.push({Name: key, X: [], Y:[]})
  }


  for (let ts of rawData) {
     measures[keys[ts.key]].X.push(getFormattedRop(ts.rop))
     measures[keys[ts.key]].Y.push(ts.trafficKpi[name])
  }
  return measures;
}

function getFormattedRop(rop) {
   return ("0" + rop.getHours()).slice(-2) + ":" + ("0" + rop.getMinutes()).slice(-2)
}

export function buildGraphData(rawData) {

  var keys = extractKeys(rawData)

  //console.log(rawData);
  var measures = {
    "SumVolumeRx": createSumLabelsValues(keys, rawData, "volumerx"),
    "SumVolumeTx": createSumLabelsValues(keys, rawData, "volumetx"),
    "RateRx": createLinearXY(keys, rawData, 'raterx'),
    "RateTx": createLinearXY(keys, rawData, 'ratetx'),
    "VolumeRx": createLinearXY(keys, rawData, 'volumerx'),
    "VolumeTx": createLinearXY(keys, rawData, 'volumetx'),
    "SpeedRx": createLinearXY(keys, rawData, 'speedrx'),
    "SpeedTx": createLinearXY(keys, rawData, 'speedtx')
  }

  return measures
}
