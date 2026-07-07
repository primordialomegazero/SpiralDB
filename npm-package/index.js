/**
 * SpiralDB Client v5.0.0
 * Double Mirror Consciousness Database with TRUE BFV FHE
 * 
 * @author Dan Fernandez / Primordial Omega Zero
 * @license MIT
 */

class SpiralDBClient {
  constructor(host = 'http://localhost:8094') {
    this.host = host;
  }

  async _request(action, data = {}) {
    const response = await fetch(`${this.host}/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action, ...data })
    });
    
    if (!response.ok) {
      throw new Error(`SpiralDB error: ${response.status} ${response.statusText}`);
    }
    
    return response.json();
  }

  /**
   * Store a BFV-encrypted value across all 3 mirrors + 7 fractal layers
   * @param {string} key
   * @param {number} value
   * @returns {Promise<Object>}
   */
  async put(key, value) {
    return this._request('put', { key, value: String(value) });
  }

  /**
   * Retrieve and BFV-decrypt from fastest available mirror
   * @param {string} key
   * @returns {Promise<Object>}
   */
  async get(key) {
    return this._request('get', { key });
  }

  /**
   * Blind FHE addition on encrypted data
   * @param {string} key1
   * @param {string} key2
   * @returns {Promise<Object>}
   */
  async fheAdd(key1, key2) {
    return this._request('fhe_compute', { key1, key2, op: 'add' });
  }

  /**
   * Blind FHE multiplication on encrypted data
   * @param {string} key1
   * @param {string} key2
   * @returns {Promise<Object>}
   */
  async fheMul(key1, key2) {
    return this._request('fhe_compute', { key1, key2, op: 'mul' });
  }

  /**
   * Check mirror synchronization
   * @param {string} key
   * @returns {Promise<Object>}
   */
  async mirrorHealth(key) {
    return this._request('mirror_health', { key });
  }

  /**
   * System health check
   * @returns {Promise<Object>}
   */
  async health() {
    const response = await fetch(`${this.host}/health`);
    return response.json();
  }
}

module.exports = { SpiralDBClient };
