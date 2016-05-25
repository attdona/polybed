
function createArray(name, objectValues, weight) {
  var measures = [
    ["context", name]
  ]

  for (var key of Object.keys(objectValues)) {
    if (weight) {
      measures.push([key, 100 * objectValues[key] / weight[key]])
    } else {
      measures.push([key, objectValues[key]])
    }
  }

  return measures
}

function createLinearRopArray(name, srcobj) {
  var measures = [
    ['rop']
  ]

  for (var first in srcobj[name]) break;
  let keys = Object.keys(srcobj[name][first])
  measures[0].push(...keys)

  for (var rop of Object.keys(srcobj[name]).sort()) {
    let measure = [rop]
    for (var i = 0; i < keys.length; i++) {
      measure.push(srcobj[name][rop][keys[i]])
    }
    measures.push(measure)
  }
  return measures;
}

function linearRop(rops, ts, measureName) {
  if (!(ts.rop in rops[measureName])) {
    rops[measureName][ts.rop] = {}
  }
  rops[measureName][ts.rop][ts.key] = ts.trafficKpi[measureName]
}

function weighted(rops, values, view, volume_sum) {
  let volumerx = rops[view]

  for (var rop in rops[view]) {
    //console.log(rop)
    let item = rops[view][rop]
    for (var key in item) {
      //console.log("key: ", key);
      //console.log("vol: ", volumerx[rop][key])
      if (!(key in values)) {
        values[key] = item[key] * volumerx[rop][key]
        volume_sum[key] = volumerx[rop][key]
        //console.log("item: ", item[key]);
      } else {
        //console.log("item else: ", item[key]);
        values[key] += item[key] * volumerx[rop][key]
        volume_sum[key] += volumerx[rop][key]
      }
    }
  }

  let s = 0
  // console.log(values)
  for (let item in values) {
    //values[item] = values[item] / volume_sum[item]
    s += values[item]
  }

  for (let item in values) {
    //console.log(`${item} : ${values[item]/s}`);
    values[item] = values[item]/s
  }
  return values
}

export function buildGraphData(rawData) {
  //console.log(rawData);
  var weightedAvgRateRx = {}
  var totVolumeRx = {}

  var weightedAvgRateTx = {}
  var totVolumeTx = {}

  var rateTx = {}

  var res = {}

  var rops = {
    raterx: {},
    ratetx: {},
    volumerx: {},
    volumetx: {},
    speedrx: {},
    speedtx: {}
  }

  for (let ts of rawData) {
    linearRop(rops, ts, 'raterx')
    linearRop(rops, ts, 'ratetx')
    linearRop(rops, ts, 'volumerx')
    linearRop(rops, ts, 'volumetx')
    linearRop(rops, ts, 'speedrx')
    linearRop(rops, ts, 'speedtx')
  }

  weighted(rops, weightedAvgRateRx, 'raterx', totVolumeRx)
  weighted(rops, weightedAvgRateTx, 'ratetx', totVolumeTx)

  //console.log(rops)
  var measures = {
    "SumVolumeRx": createArray("VolumeRx", totVolumeRx),
    "SumVolumeTx": createArray("VolumeTx", totVolumeTx),
    "RateRx": createLinearRopArray('raterx', rops),
    "RateTx": createLinearRopArray('ratetx', rops),
    "VolumeRx": createLinearRopArray('volumerx', rops),
    "VolumeTx": createLinearRopArray('volumetx', rops),
    "SpeedRx": createLinearRopArray('speedrx', rops),
    "SpeedTx": createLinearRopArray('speedtx', rops)
  }

  //console.log(measures)
  return measures
}
