const http = require('http');

const PHI = 1.6180339887498948482;
const LAMBDA = 0.4812;

class SpiralDBClient {
  constructor(host = 'localhost', port = 8094) {
    this.host = host;
    this.port = port;
  }

  _post(data) {
    return new Promise((resolve, reject) => {
      const body = JSON.stringify(data);
      const opts = {
        hostname: this.host, port: this.port, path: '/', method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Content-Length': Buffer.byteLength(body) }
      };
      const req = http.request(opts, (res) => {
        let raw = '';
        res.on('data', c => raw += c);
        res.on('end', () => { try { resolve(JSON.parse(raw)); } catch(e) { resolve({ raw }); } });
      });
      req.on('error', reject);
      req.write(body);
      req.end();
    });
  }

  async put(key, value) { return this._post({ action: 'put', key, value }); }
  async get(key) { return this._post({ action: 'get', key }); }
  async fheAdd(key1, key2) { return this._post({ action: 'fhe_compute', key1, key2, op: 'add' }); }
  async fheMul(key1, key2) { return this._post({ action: 'fhe_compute', key1, key2, op: 'mul' }); }
  async mirrorHealth(key) { return this._post({ action: 'mirror_health', key }); }
  async health() { return this._post({ action: 'health' }); }

  encrypt(value) { return value * PHI + LAMBDA; }
  decrypt(encrypted) { return Math.round((encrypted - LAMBDA) / PHI); }
}

module.exports = { SpiralDBClient };
