import http from 'http';
import assert from 'assert';

import '../lib/index.js';
import {TrafficSnippet} from '../lib/dpi';


function *randomTraffic(done) {
  var num = 5;
  console.log(`generating ${num} records`);

  for (var i = 0; i < num; i += 1) {
    console.log(`${i}`)
    var item = new TrafficSnippet({
      pool: 'linea1',
      rop: new Date(2016, 4, 10, 6 ,0, 0),
      trafficKpi: {
        raterx: Math.random(),
        ratetx: Math.random(),
        volumerx: Math.random()*100,
        volumetx: Math.random()*100,
        speedrx : Math.random()*50,
        speedtx : Math.random()*50
      }
    });

    yield item.save();
  }
  done();
}

describe('Example Node Server', () => {
  it('should return 200', done => {
    http.get('http://127.0.0.1:3000', res => {
      assert.equal(200, res.statusCode);
      done();
    });
  });
});

describe('Insert a record', () => {
  it('should save a traffic record', function *(done) {
    var item = new TrafficSnippet();
    item.pool = 'test';
    yield randomTraffic(done);
    //item.save();
    //done();
  });
});
