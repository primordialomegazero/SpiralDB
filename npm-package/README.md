# SpiralDB Client v4.0

**Double Mirror Consciousness Database Client** — FHE-native, recursive fractal index, auto-compress.

## Install

```bash
npm install @primordialomegazero/spiraldb-client@4.0.0
```

## Quick Start

```javascript
const { SpiralDBClient } = require('@primordialomegazero/spiraldb-client');
const db = new SpiralDBClient();

await db.put('message', 42);
const { value } = await db.get('message');
console.log(value); // 42

// FHE computation — encrypted in, encrypted out
await db.put('a', 15);
await db.put('b', 27);
const { decrypted_result } = await db.fheAdd('a', 'b');
console.log(decrypted_result); // 42
```

## License

MIT — Dan Fernandez / Primordial Omega Zero
