// https://eladnava.com/write-synchronous-node-js-code-with-es6-generators/

// Dependencies
require('babel-polyfill');

var koa = require('koa');
var serve = require('koa-static');
var route = require('koa-route');

var mongoose  = require('mongoose');

import {TrafficSnippet} from './dpi';

// Set up MongoDB connection
var connection = mongoose.connect('localhost/netdata');


// Create koa app
var app = koa();

app.use(serve('app'));

app.use(route.get('/api/t', tdata));


function *tdata() {

    var item = new TrafficSnippet();

    item.pool = "xxxxxxx";
    yield item.save();

    var data = yield TrafficSnippet.find({});
    //var data = yield Bear.find({});
    this.body = data;
}


// Define configurable port
var port = process.env.PORT || 3000;

// Listen for connections
app.listen(port);

// Log port
console.log('Server listening on port ' + port);
