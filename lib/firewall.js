import mongoose from 'mongoose';

export var Service = mongoose.model('service', new mongoose.Schema({
  name: String
}));
