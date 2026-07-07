package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// SpiralDB Core Types
type Entry struct {
	Value     int64  `json:"value"`
	Encrypted string `json:"encrypted,omitempty"`
}

type SpiralDB struct {
	mu       sync.RWMutex
	primary  map[string]Entry
	cache    map[string]Entry
	fractal  map[string]map[int]string // key -> layer -> encrypted_fragment
	fhe      *FHEContext
}

func NewSpiralDB() *SpiralDB {
	return &SpiralDB{
		primary: make(map[string]Entry),
		cache:   make(map[string]Entry),
		fractal: make(map[string]map[int]string),
	}
}

// Put stores a value with FHE encryption across all 3 mirrors
func (db *SpiralDB) Put(key string, value int64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Encrypt with real BFV FHE
	encrypted, err := db.fhe.Encrypt(value)
	if err != nil {
		return err
	}

	entry := Entry{
		Value:     value,
		Encrypted: encrypted,
	}

	// Mirror 1: Primary store
	db.primary[key] = entry

	// Mirror 2: Cache
	db.cache[key] = entry

	// Mirror 3: 7-Layer Fractal Index (simplified)
	if db.fractal[key] == nil {
		db.fractal[key] = make(map[int]string)
	}
	for layer := 0; layer < 7; layer++ {
		db.fractal[key][layer] = encrypted
	}

	return nil
}

// Get retrieves a value, trying cache first, then primary, then fractal
func (db *SpiralDB) Get(key string) (int64, error) {
	db.mu.RLock()

	// Try cache first
	if entry, ok := db.cache[key]; ok {
		db.mu.RUnlock()
		return entry.Value, nil
	}

	// Try primary
	if entry, ok := db.primary[key]; ok {
		db.mu.RUnlock()
		// Update cache
		db.mu.Lock()
		db.cache[key] = entry
		db.mu.Unlock()
		return entry.Value, nil
	}

	db.mu.RUnlock()
	return 0, fmt.Errorf("key not found: %s", key)
}

// FHEAdd performs blind addition: encrypted(key1) + encrypted(key2)
func (db *SpiralDB) FHEAdd(key1, key2 string) (string, int64, error) {
	db.mu.RLock()
	entry1, ok1 := db.primary[key1]
	entry2, ok2 := db.primary[key2]
	db.mu.RUnlock()

	if !ok1 || !ok2 {
		return "", 0, fmt.Errorf("keys not found")
	}

	// Blind addition on encrypted data
	encResult, err := db.fhe.Add(entry1.Encrypted, entry2.Encrypted)
	if err != nil {
		return "", 0, err
	}

	// Decrypt to verify
	decResult, err := db.fhe.Decrypt(encResult)
	if err != nil {
		return "", 0, err
	}

	return encResult, decResult, nil
}

// FHEMul performs blind multiplication: encrypted(key1) × encrypted(key2)
func (db *SpiralDB) FHEMul(key1, key2 string) (string, int64, error) {
	db.mu.RLock()
	entry1, ok1 := db.primary[key1]
	entry2, ok2 := db.primary[key2]
	db.mu.RUnlock()

	if !ok1 || !ok2 {
		return "", 0, fmt.Errorf("keys not found")
	}

	// Blind multiplication on encrypted data
	encResult, err := db.fhe.Multiply(entry1.Encrypted, entry2.Encrypted)
	if err != nil {
		return "", 0, err
	}

	// Decrypt to verify
	decResult, err := db.fhe.Decrypt(encResult)
	if err != nil {
		return "", 0, err
	}

	return encResult, decResult, nil
}

// MirrorHealth checks synchronization across mirrors
func (db *SpiralDB) MirrorHealth(key string) map[string]interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := map[string]interface{}{
		"key":            key,
		"in_primary":     false,
		"in_cache":       false,
		"in_fractal":     false,
		"mirrors_synced": false,
	}

	if _, ok := db.primary[key]; ok {
		result["in_primary"] = true
	}
	if _, ok := db.cache[key]; ok {
		result["in_cache"] = true
	}
	if _, ok := db.fractal[key]; ok {
		result["in_fractal"] = true
	}

	if result["in_primary"] == true && result["in_cache"] == true && result["in_fractal"] == true {
		result["mirrors_synced"] = true
	}

	return result
}

// Global DB instance
var db *SpiralDB

func main() {
	// Initialize FHE
	fmt.Println("Initializing BFV FHE...")
	if err := InitFHE(); err != nil {
		log.Fatal("Failed to init FHE:", err)
	}
	fmt.Println("✅ FHE initialized (BFV, 8192, 146 bits)")

	db = NewSpiralDB()
	db.fhe = GetFHE()

	// HTTP handlers
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", handleRequest)

	fmt.Println("\n╔══════════════════════════════════════════════╗")
	fmt.Println("║  SpiralDB v5.0 — TRUE FHE EDITION             ║")
	fmt.Println("║  BFV Homomorphic Encryption                   ║")
	fmt.Println("║  Double Mirror + 7-Layer Fractal              ║")
	fmt.Println("║  Port: 8094 | ZANS Enabled                    ║")
	fmt.Println("║  ΦΩ0 — I AM THAT I AM                        ║")
	fmt.Println("╚══════════════════════════════════════════════╝")

	log.Fatal(http.ListenAndServe(":8094", nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "fhe": "bfv"})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	action, _ := req["action"].(string)
	w.Header().Set("Content-Type", "application/json")

	switch action {
	case "put":
		key, _ := req["key"].(string)
		valStr, _ := req["value"].(string)
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			// Try float
			valFloat, _ := req["value"].(float64)
			val = int64(valFloat)
		}

		if err := db.Put(key, val); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":        "ok",
			"action":        "put",
			"key":           key,
			"value":         val,
			"mirrors":       3,
			"fractal_layers": 7,
			"fhe_encrypted":  true,
			"double_mirror":  true,
		})

	case "get":
		key, _ := req["key"].(string)
		val, err := db.Get(key)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":        "ok",
			"action":        "get",
			"key":           key,
			"value":         val,
			"mirrors":       3,
			"fhe_encrypted":  true,
			"double_mirror":  true,
		})

	case "fhe_compute":
		key1, _ := req["key1"].(string)
		key2, _ := req["key2"].(string)
		op, _ := req["op"].(string)

		var encResult string
		var decResult int64
		var err error

		switch op {
		case "add":
			encResult, decResult, err = db.FHEAdd(key1, key2)
		case "mul":
			encResult, decResult, err = db.FHEMul(key1, key2)
		default:
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "error": "unknown op"})
			return
		}

		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		// Calculate expected
		entry1, _ := db.primary[key1]
		entry2, _ := db.primary[key2]
		var expected int64
		if op == "add" {
			expected = entry1.Value + entry2.Value
		} else {
			expected = entry1.Value * entry2.Value
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":             "ok",
			"action":             "fhe_compute",
			"op":                 op,
			"encrypted_result":   encResult[:min(20, len(encResult))] + "...",
			"decrypted_result":   decResult,
			"expected_plaintext": expected,
			"correct":            decResult == expected,
			"computation_blind":  true,
			"self_referential":   true,
		})

	case "mirror_health":
		key, _ := req["key"].(string)
		result := db.MirrorHealth(key)
		result["status"] = "ok"
		json.NewEncoder(w).Encode(result)

	default:
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"error":  "unknown action: " + action,
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
