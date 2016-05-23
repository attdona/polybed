import http from 'http';
import assert from 'assert';

import '../lib/index.js';
import {
  TrafficSnippet
} from '../lib/dpi';

import {
  Service
} from '../lib/firewall.js'

function* randomWebTraffic() {
  let num = 1;
  let types = ['http', 'mail', 'p2p']
  let webs = ['noiportal.it', 'youporn']

  for (let type of types) {
    for (let web of webs) {
      for (var i = 0; i < num; i += 1) {
        //console.log(`${i}`)
        let datetime = new Date(2016, 4, 10, 2 + i, 0, 0);
        var item = new TrafficSnippet({
          pool: 'linea1',
          rop: datetime,
          context: type,
          key: web,
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

function* randomTraffic() {
  let num = 2;
  let types = ['net']
  let keys = ['http', 'mail']

  for (let type of types) {
    for (let key of keys) {
      for (var i = 0; i < num; i += 1) {
        //console.log(`${i}`)
        let datetime = new Date(2016, 4, 10, 2 + i, 0, 0);
        var item = new TrafficSnippet({
          pool: 'linea1',
          rop: datetime,
          context: type,
          key: key,
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
}

function* mathTraffic(type, keys) {
  let num = 10;
  //let types = ['net']
  //let keys = ['http', 'mail', 'others']
  let frequency = 1

  var raterx_sum = 0
  var ratetx_sum = 0

  var raterx_grid = []
  var ratetx_grid = []

  for (var i = 0; i < num; i += 1) {
    var raterx_vals = []
    var ratetx_vals = []
    let rxsum = 0
    let txsum = 0
    for (var key_idx = 0; key_idx < keys.length; key_idx++) {
      let rxval = 0.5 * (1 + Math.sin(i * frequency + key_idx))
      let txval = 0.5 * (1 + Math.cos(i * frequency + key_idx))
      raterx_vals.push(rxval)
      ratetx_vals.push(txval)
      rxsum += rxval
      txsum += txval
    }
    raterx_vals = raterx_vals.map(v => {
      return v / rxsum
    })
    ratetx_vals = ratetx_vals.map(v => {
      return v / txsum
    })
    raterx_grid.push(raterx_vals)
    ratetx_grid.push(ratetx_vals)

  }

  for (var key_idx = 0; key_idx < keys.length; key_idx++) {
    for (var i = 0; i < num; i += 1) {
      //console.log(`${i}`)
      let datetime = new Date(2016, 4, 10, 2 + i, 0, 0);
      var item = new TrafficSnippet({
        pool: 'linea1',
        rop: datetime,
        context: type,
        key: keys[key_idx],
        trafficKpi: {
          raterx: raterx_grid[i][key_idx],
          ratetx: ratetx_grid[i][key_idx],
          volumerx: 100 * (Math.sin(i * frequency * 2 + key_idx) + 1),
          volumetx: 100 * (Math.cos(i * frequency * 2 + key_idx) + 1),
          speedrx: 50,
          speedtx: 80
        }
      });

      yield item.save();
    }
  }

}


describe('Example Node Server', () => {
  it('should return 200', done => {
    http.get('http://127.0.0.1:3000', res => {
      assert.equal(200, res.statusCode);
      done();
    });
  });
});

describe('cleanup database', () => {
  it('shouls remove old data', function*() {
    yield TrafficSnippet.collection.remove()
  })
})

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
    yield mathTraffic('net', ['http', 'mail', 'others']);
    yield mathTraffic('http', ['noi.portal', 'facebook', 'google', 'youporn']);
  });
});

describe('get the measures', () => {
  it('get the measures', done => {
    http.get('http://127.0.0.1:3000/api/net/linea1', res => {
      assert.equal(200, res.statusCode);
      let body = ''
      res.on('data', d => {
        body += d;
      })
      res.on('end', () => {
        //console.log(body)
        done()
      })
    })
  })
})
