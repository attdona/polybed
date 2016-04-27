import http from 'http';
import assert from 'assert';

import '../lib/index.js';
import {
  TrafficSnippet
} from '../lib/dpi';

import {Service} from '../lib/firewall.js'

function* randomTraffic() {
  var num = 5;
  let types = ['http', 'mail', 'p2p']
  let webs = ['google', 'facebook', 'noiportal.it', 'youporn']

  for (let type of types) {
    for (let web of webs) {
      for (var i = 0; i < num; i += 1) {
        console.log(`${i}`)
        let datetime = new Date(2016, 4, 10, 2+i, 0, 0);
        var item = new TrafficSnippet({
          pool: 'linea1',
          rop: datetime,
          trafficKpi: {
            raterx: Math.random(),
            ratetx: Math.random(),
            volumerx: Math.random() * 100,
            volumetx: Math.random() * 100,
            speedrx: Math.random() * 50,
            speedtx: Math.random() * 50
          }
        });
        yield item.save();
      }
      console.log(`generating ${type} records`);
    }
  }
  //done();
}

describe('Example Node Server', () => {
  it('should return 200', done => {
    http.get('http://127.0.0.1:3000', res => {
      assert.equal(200, res.statusCode);
      done();
    });
  });
});

describe('Insert FwServices', () => {
  it('should insert firewall services', function*() {
    let services = ['http (port 80)', 'you porn', 'bit torrent', 'facebook', 'trespolo']

    yield Service.collection.remove();

    for (let name of services) {
      var srv = new Service();
      srv.name = name;
      yield srv.save();
    }
  });
});

describe('Insert a batch of records', () => {
  it('should save a traffic record', function*() {
    var item = new TrafficSnippet();
    item.pool = 'test';
    yield randomTraffic();
    console.log("FIN");
  });
});
