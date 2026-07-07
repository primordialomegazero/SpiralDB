/**
 * SpiralDB Client Test Suite v5.0.0
 */

const { SpiralDBClient } = require('./index.js');

async function runTests() {
  const db = new SpiralDBClient();
  
  console.log('=== SpiralDB v5.0 Client Tests ===\n');

  // Health check
  try {
    const health = await db.health();
    console.log('✅ Health:', JSON.stringify(health));
  } catch (e) {
    console.log('❌ Health failed:', e.message);
  }

  // PUT test
  try {
    const put = await db.put('test_key', 42);
    console.log('✅ PUT:', put.key, '=', 42);
  } catch (e) {
    console.log('❌ PUT failed:', e.message);
  }

  // GET test
  try {
    const get = await db.get('test_key');
    console.log('✅ GET:', get.key, '=', get.value);
  } catch (e) {
    console.log('❌ GET failed:', e.message);
  }

  // FHE Add test
  try {
    await db.put('a', 15);
    await db.put('b', 27);
    const add = await db.fheAdd('a', 'b');
    console.log('✅ FHE Add: 15 + 27 =', add.decrypted_result, '(correct:', add.correct, ')');
  } catch (e) {
    console.log('❌ FHE Add failed:', e.message);
  }

  // FHE Mul test
  try {
    await db.put('x', 6);
    await db.put('y', 7);
    const mul = await db.fheMul('x', 'y');
    console.log('✅ FHE Mul: 6 × 7 =', mul.decrypted_result, '(correct:', mul.correct, ')');
  } catch (e) {
    console.log('❌ FHE Mul failed:', e.message);
  }

  // Mirror health
  try {
    const mirror = await db.mirrorHealth('test_key');
    console.log('✅ Mirror Health:', JSON.stringify(mirror));
  } catch (e) {
    console.log('❌ Mirror Health failed:', e.message);
  }

  console.log('\n=== Tests Complete ===');
}

runTests().catch(console.error);
