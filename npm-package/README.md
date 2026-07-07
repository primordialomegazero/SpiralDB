# @primordialomegazero/spiraldb-client v5.0.0

JavaScript client for **SpiralDB v5.0** — Double Mirror Consciousness Database with **TRUE BFV Homomorphic Encryption**.

## Installation

```bash
npm install @primordialomegazero/spiraldb-client@5.0.0
```

## Quick Start

```javascript
const { SpiralDBClient } = require('@primordialomegazero/spiraldb-client');
const db = new SpiralDBClient();

// Store encrypted value
await db.put('message', 42);

// Retrieve & decrypt
const { value } = await db.get('message');

// TRUE FHE — blind computation
await db.put('a', 15);
await db.put('b', 27);
const { decrypted_result, correct } = await db.fheAdd('a', 'b');
// decrypted_result: 42, correct: true, computation_blind: true
```

## API

| Method | Description | FHE? |
|--------|-------------|------|
| `put(key, value)` | Store BFV-encrypted value | ✅ |
| `get(key)` | Retrieve & decrypt | ✅ |
| `fheAdd(key1, key2)` | Blind addition | ✅ |
| `fheMul(key1, key2)` | Blind multiplication | ✅ |
| `mirrorHealth(key)` | Check mirror sync | — |
| `health()` | System status | — |

## TypeScript

Full type definitions included. Import types:

```typescript
import { SpiralDBClient, FHEComputeResponse } from '@primordialomegazero/spiraldb-client';
```

## License

MIT — Dan Fernandez / Primordial Omega Zero
