// https://eladnava.com/write-synchronous-node-js-code-with-es6-generators/

// Dependencies
require('babel-polyfill');

var koa = require('koa');
var serve = require('koa-static');
var route = require('koa-route');

var mongoose = require('mongoose');

import {
  TrafficSnippet
} from './dpi';
import {
  Service
} from './firewall';

// Set up MongoDB connection
var connection = mongoose.connect('localhost/netdata');


// Create koa app
var app = koa();

app.use(serve('app'));

app.use(route.get('/api/net/:line', netdata));
app.use(route.get('/api/fwservices', fwservices));


function* netdata(line) {

  var item = new TrafficSnippet();
  var data = yield TrafficSnippet.find({
    "pool": line
  });
  this.body = data;
}

function* fwservices() {
  this.body = yield Service.find({}).select({
    "name": 1,
    "_id": 0
  });
}

// Define configurable port
var port = process.env.PORT || 3000;

// Listen for connections
app.listen(port);

// Log port
console.log('Server listening on port ' + port);
