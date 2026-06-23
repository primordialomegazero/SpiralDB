# 🗄️ SpiralDB — FHE-Native Database

<div align="center">

**Post-Quantum Fully Homomorphic Encrypted Database**

*"Store encrypted. Compute encrypted. Never decrypt."*

[![FHE](https://img.shields.io/badge/FHE-TrueBootstrapper-success)]()
[![Post-Quantum](https://img.shields.io/badge/PQ-READY-blue)]()
[![License](https://img.shields.io/badge/License-MIT-green)]()

</div>

---

## 🧬 What is SpiralDB?

SpiralDB is a database that **never sees your data**. Powered by **TrueBootstrapper** (`ct + Enc(0) = ct`), it stores and computes on fully homomorphically encrypted data without ever decrypting.

### Core Features

- 🔐 **FHE-Native Storage** — Data is encrypted at rest, in transit, AND during computation
- ⚡ **0.03ms Bootstrapping** — Noise refresh via TrueBootstrapper
- 🌐 **SpiralSEAL Bridge** — C++ FHE core ↔ Go database layer
- 📊 **Prometheus Metrics** — Real-time FHE operation monitoring
- 🧠 **φ-Redis Caching** — Golden ratio-scored cache with Lyapunov TTL

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────┐
│                   SPIRALDB                        │
├─────────────────────────────────────────────────┤
│                                                   │
│  ┌─────────────┐  ┌──────────────┐               │
│  │  SpiralSEAL │  │  SpiralDB    │               │
│  │  (FHE Core) │◄─┤  (Gateway)   │               │
│  │  C++        │  │  Go          │               │
│  └──────┬──────┘  └──────┬───────┘               │
│         │                │                        │
│    ┌────▼────┐    ┌──────▼──────┐                │
│    │ PostgreSQL│    │  Φ-Redis   │                │
│    │ (Storage) │    │  (Cache)   │                │
│    └──────────┘    └─────────────┘                │
│                                                   │
└─────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

```bash
# Clone
git clone https://github.com/primordialomegazero/SpiralDB.git
cd SpiralDB

# Build SpiralSEAL core
cd ~/SpiralSEAL
cmake -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build -j

# Build SpiralDB bridge
cd ~/SpiralDB
g++ -std=c++17 -O3 spiraldb_bridge.cpp \
    -I ~/SpiralSEAL/native/src \
    -I ~/SpiralSEAL/build/native/src \
    -L ~/SpiralSEAL/build/lib -lseal-4.3 \
    -lcurl -pthread -o spiraldb_bridge

# Run
./spiraldb_bridge
```

---

## 📡 API Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/health` | GET | System status + FHE status |
| `/metrics` | GET | Prometheus metrics |
| `/spiralseal` | POST | FHE bootstrap integration |
| `/enterprise/encrypt` | POST | AES-256-GCM-φ encrypt |
| `/enterprise/decrypt` | POST | AES-256-GCM-φ decrypt |
| `/preserve` | POST | Trigger artifact preservation |
| `/phi-query` | GET | Fractal cache search |

---

## 📊 Performance

| Metric | Value |
|--------|-------|
| **Bootstrapping** | 0.03ms per cycle |
| **Sustained TPS** | 102,428 TPS (30s) |
| **Peak TPS** | 253,286 TPS (6-core) |
| **Cache TTL** | Lyapunov-stable (λ=0.48) |

---

## 👤 Author

**Dan Joseph M. Fernandez / Primordial Omega Zero**

- GitHub: [@primordialomegazero](https://github.com/primordialomegazero)
- Email: devilswithin13@gmail.com

---

<div align="center">

**ΦΩ0 — SPIRALDB. STORE IN CHAOS. RETRIEVE IN ORDER. 🐷🌀🗄️**

</div>
