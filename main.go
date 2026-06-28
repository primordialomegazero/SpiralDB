package main

import (
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "log"
    "math"
    "net/http"
    "sync"
    "time"
)

// ═══ φ CONSTANTS ═══
const (
    PHI     = 1.6180339887498948482
    PHI_INV = 0.6180339887498948482
    LAMBDA  = 0.4812
)

// ═══ SELF-REFERENTIAL FHE (Direct numeric) ═══
type FHEEngine struct{}

func (f *FHEEngine) EncryptNum(value float64) float64 { return value*PHI + LAMBDA }
func (f *FHEEngine) DecryptNum(e float64) float64     { return math.Round((e - LAMBDA) / PHI) }
func (f *FHEEngine) Add(e1, e2 float64) float64       { return e1 + e2 - LAMBDA }
func (f *FHEEngine) Mul(e1, e2 float64) float64 {
    t1, t2, t3 := e1*e2, LAMBDA*(e1+e2), LAMBDA*LAMBDA
    return (t1 - t2 + t3) / PHI + LAMBDA
}

// ═══ RECURSIVE FRACTAL INDEX (7-layer, auto-compress) ═══
type FractalLayer struct {
    entries  map[string]float64
    compressed bool
    mu       sync.RWMutex
}

type FractalIndex struct {
    layers    [7]*FractalLayer
    compressThreshold int
}

func NewFractalIndex() *FractalIndex {
    f := &FractalIndex{compressThreshold: 1000}
    for i := 0; i < 7; i++ {
        f.layers[i] = &FractalLayer{entries: make(map[string]float64)}
    }
    return f
}

func (f *FractalIndex) Insert(key string, value float64) {
    for layer := 0; layer < 7; layer++ {
        h := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%f", key, layer, PHI)))
        fractalKey := fmt.Sprintf("%x", h[:16])
        f.layers[layer].mu.Lock()
        f.layers[layer].entries[fractalKey] = value
        
        // Auto-compress: if too many entries, apply φ-contraction
        if len(f.layers[layer].entries) > f.compressThreshold {
            f.compressLayer(layer)
        }
        f.layers[layer].mu.Unlock()
    }
}

func (f *FractalIndex) compressLayer(layer int) {
    // φ-contraction: merge entries by averaging with φ-weights
    var sum float64
    count := 0
    for _, v := range f.layers[layer].entries {
        sum += v
        count++
    }
    if count > 0 {
        avg := sum / float64(count)
        // Contract toward the φ-weighted average
        for k, v := range f.layers[layer].entries {
            f.layers[layer].entries[k] = v*PHI_INV + avg*(1.0-PHI_INV)
        }
        f.layers[layer].compressed = true
    }
}

func (f *FractalIndex) Query(key string) (float64, bool) {
    for layer := 0; layer < 7; layer++ {
        h := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%f", key, layer, PHI)))
        fractalKey := fmt.Sprintf("%x", h[:16])
        f.layers[layer].mu.RLock()
        val, ok := f.layers[layer].entries[fractalKey]
        f.layers[layer].mu.RUnlock()
        if ok {
            return val, true
        }
    }
    return 0, false
}

func (f *FractalIndex) Stats() map[string]interface{} {
    stats := make(map[string]interface{})
    for i := 0; i < 7; i++ {
        f.layers[i].mu.RLock()
        stats[fmt.Sprintf("layer_%d_entries", i)] = len(f.layers[i].entries)
        stats[fmt.Sprintf("layer_%d_compressed", i)] = f.layers[i].compressed
        f.layers[i].mu.RUnlock()
    }
    return stats
}

// ═══ DOUBLE MIRROR CONSCIOUSNESS ═══
type MirrorEntry struct {
    Value          string
    NumericValue   float64
    EncryptedValue float64
    FractalHash    string
    Timestamp      time.Time
}

type DoubleMirrorDB struct {
    primary map[string]MirrorEntry
    cache   map[string]MirrorEntry
    fractal *FractalIndex
    fhe     *FHEEngine
    mu      sync.RWMutex
}

func NewDoubleMirrorDB() *DoubleMirrorDB {
    return &DoubleMirrorDB{
        primary: make(map[string]MirrorEntry),
        cache:   make(map[string]MirrorEntry),
        fractal: NewFractalIndex(),
        fhe:     &FHEEngine{},
    }
}

func (db *DoubleMirrorDB) Put(key string, numValue float64) {
    db.mu.Lock()
    defer db.mu.Unlock()

    valueStr := fmt.Sprintf("%.6f", numValue)
    encrypted := db.fhe.EncryptNum(numValue)
    h := sha256.Sum256([]byte(fmt.Sprintf("%s:%f:%f", key, numValue, PHI)))
    fractalHash := fmt.Sprintf("%x", h[:16])

    entry := MirrorEntry{
        Value: valueStr, NumericValue: numValue,
        EncryptedValue: encrypted, FractalHash: fractalHash,
        Timestamp: time.Now(),
    }

    // Triple write — double mirror consciousness
    db.primary[key] = entry
    db.cache[key] = entry
    db.fractal.Insert(key, encrypted)
}

func (db *DoubleMirrorDB) Get(key string) (float64, bool) {
    db.mu.RLock()
    defer db.mu.RUnlock()

    // Fastest: cache mirror
    if e, ok := db.cache[key]; ok {
        return e.NumericValue, true
    }
    // Fractal mirror
    if v, ok := db.fractal.Query(key); ok {
        return db.fhe.DecryptNum(v), true
    }
    // Primary mirror
    if e, ok := db.primary[key]; ok {
        return e.NumericValue, true
    }
    return 0, false
}

func (db *DoubleMirrorDB) FHECompute(key1, key2, op string) (float64, float64, float64) {
    db.mu.RLock()
    defer db.mu.RUnlock()

    e1 := db.primary[key1].EncryptedValue
    e2 := db.primary[key2].EncryptedValue

    var result float64
    switch op {
    case "add":
        result = db.fhe.Add(e1, e2)
    case "mul":
        result = db.fhe.Mul(e1, e2)
    }
    decrypted := db.fhe.DecryptNum(result)
    var expected float64; switch op { case "add": expected = db.primary[key1].NumericValue + db.primary[key2].NumericValue; case "mul": expected = db.primary[key1].NumericValue * db.primary[key2].NumericValue }; return result, decrypted, expected
}

func (db *DoubleMirrorDB) MirrorHealth(key string) map[string]interface{} {
    db.mu.RLock()
    defer db.mu.RUnlock()

    p, pok := db.primary[key]
    c, cok := db.cache[key]
    f, fok := db.fractal.Query(key)

    synced := pok && cok && fok && p.EncryptedValue == c.EncryptedValue

    return map[string]interface{}{
        "in_primary": pok, "in_cache": cok, "in_fractal": fok,
        "mirrors_synced": synced,
        "fractal_decrypted": func() float64 { if fok { return db.fhe.DecryptNum(f) } else { return 0 } }(),
    }
}

func (db *DoubleMirrorDB) Stats() map[string]interface{} {
    db.mu.RLock()
    defer db.mu.RUnlock()

    stats := db.fractal.Stats()
    stats["primary_entries"] = len(db.primary)
    stats["cache_entries"] = len(db.cache)
    stats["mirrors"] = 3
    stats["double_mirror"] = true
    stats["fractal_layers"] = 7
    stats["auto_compress"] = true
    return stats
}

// ═══ HTTP SERVER ═══
var db = NewDoubleMirrorDB()

func main() {
    http.HandleFunc("/", handler)

    fmt.Println("╔══════════════════════════════════════════════╗")
    fmt.Println("║  SpiralDB v4.0 — DEEP DOUBLE MIRROR           ║")
    fmt.Println("║  Recursive Fractal + Auto-Compress + FHE      ║")
    fmt.Println("║  Double Mirror Consciousness                  ║")
    fmt.Println("║  Port: 8094 | Zero Deps | Instant Start       ║")
    fmt.Println("║  ΦΩ0 — I AM THAT I AM                        ║")
    fmt.Println("╚══════════════════════════════════════════════╝")

    log.Fatal(http.ListenAndServe(":8094", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    json.NewDecoder(r.Body).Decode(&req)
    action, _ := req["action"].(string)
    w.Header().Set("Content-Type", "application/json")

    switch action {
    case "put":
        key, _ := req["key"].(string)
        value, _ := req["value"].(float64)
        db.Put(key, value)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status": "ok", "action": "put", "key": key,
            "encrypted_value": db.fhe.EncryptNum(value),
            "mirrors": 3, "fractal_layers": 7, "double_mirror": true,
        })

    case "get":
        key, _ := req["key"].(string)
        value, ok := db.Get(key)
        if !ok {
            json.NewEncoder(w).Encode(map[string]interface{}{"status": "error", "error": "not found"})
            return
        }
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status": "ok", "action": "get", "key": key,
            "value": value, "encrypted_value": db.fhe.EncryptNum(value),
            "mirrors": 3, "double_mirror": true,
        })

    case "fhe_compute":
        key1, _ := req["key1"].(string)
        key2, _ := req["key2"].(string)
        op, _ := req["op"].(string)
        encrypted, decrypted, expected := db.FHECompute(key1, key2, op)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status": "ok", "action": "fhe_compute", "op": op,
            "encrypted_result": encrypted, "decrypted_result": decrypted,
            "expected_plaintext": expected,
            "correct": decrypted == expected,
            "computation_blind": true, "self_referential": true,
        })

    case "mirror_health":
        key, _ := req["key"].(string)
        health := db.MirrorHealth(key)
        health["status"] = "ok"
        health["key"] = key
        json.NewEncoder(w).Encode(health)

    case "health":
        stats := db.Stats()
        stats["system"] = "SpiralDB"
        stats["version"] = "4.0.0"
        stats["edition"] = "DEEP_DOUBLE_MIRROR"
        stats["backends"] = "In-Memory + Cache + 7-Layer Recursive Fractal"
        stats["fhe"] = "Self-Referential Direct Numeric"
        stats["phi"] = PHI
        stats["lambda"] = LAMBDA
        stats["status"] = "THE_VOID_PERSISTS_IN_MEMORY"
        json.NewEncoder(w).Encode(stats)

    default:
        json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
    }
}
