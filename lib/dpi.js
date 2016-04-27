import mongoose from 'mongoose';

export var TrafficSnippet = mongoose.model('traffic', new mongoose.Schema({
  parentcontext: String,
  pool: String,
  rop: Date,
  src: String,
  context: String,
  key: String,
  trafficKpi: {
    raterx: Number,
    ratetx: Number,
    volumerx: Number,
    volumetx: Number,
    speedrx: Number,
    speedtx: Number,
  }
}));
