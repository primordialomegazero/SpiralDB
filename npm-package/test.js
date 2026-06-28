const { SpiralDBClient } = require('./index.js');

async function test() {
  const db = new SpiralDBClient();
  let passed = 0, total = 8;

  console.log('╔════════════════════════════════════════╗');
  console.log('║  SpiralDB v4.0 Client Test             ║');
  console.log('╚════════════════════════════════════════╝\n');

  // 1. Health
  const h = await db.health();
  console.log(`1. Health: ${h.system === 'SpiralDB' ? '✅' : '❌'}`);
  if (h.system === 'SpiralDB') passed++;

  // 2. Put
  const put = await db.put('test', 42);
  console.log(`2. Put: ${put.status === 'ok' ? '✅' : '❌'}`);
  if (put.status === 'ok') passed++;

  // 3. Get
  const get = await db.get('test');
  console.log(`3. Get: ${get.value === 42 ? '✅' : '❌'}`);
  if (get.value === 42) passed++;

  // 4. FHE Add
  await db.put('a', 15);
  await db.put('b', 27);
  const add = await db.fheAdd('a', 'b');
  console.log(`4. FHE Add 15+27=${add.decrypted_result} ${add.correct ? '✅' : '❌'}`);
  if (add.correct) passed++;

  // 5. FHE Mul
  await db.put('x', 6);
  await db.put('y', 7);
  const mul = await db.fheMul('x', 'y');
  console.log(`5. FHE Mul 6×7=${mul.decrypted_result} ${mul.correct ? '✅' : '❌'}`);
  if (mul.correct) passed++;

  // 6. Mirror Health
  const mh = await db.mirrorHealth('test');
  console.log(`6. Mirror Synced: ${mh.mirrors_synced ? '✅' : '❌'}`);
  if (mh.mirrors_synced) passed++;

  // 7. Encrypt/Decrypt
  const enc = db.encrypt(42);
  const dec = db.decrypt(enc);
  console.log(`7. Encrypt/Decrypt: ${dec === 42 ? '✅' : '❌'}`);
  if (dec === 42) passed++;

  // 8. Double Mirror
  console.log(`8. Double Mirror: ${put.double_mirror ? '✅' : '❌'}`);
  if (put.double_mirror) passed++;

  console.log(`\n╔════════════════════════════════════════╗`);
  console.log(`║  RESULT: ${passed}/${total} passed ${passed === total ? '✅' : '❌'}          ║`);
  console.log(`╚════════════════════════════════════════╝`);
}

test().catch(e => console.error('Error:', e.message));
