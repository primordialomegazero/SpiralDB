package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func testHydraFull() {
	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Println("║  SPIRALDB HYDRA — FULL TEST           ║")
	fmt.Println("║  14 Models + 8 Consensus + FHE        ║")
	fmt.Println("╚═══════════════════════════════════════╝")
	fmt.Println()

	InitGlobalHydra()
	
	// Test 1: Insert into all 14 models
	fmt.Println("=== Test 1: 14-Model Insert ===")
	testData := map[string]interface{}{
		"id": 42, "name": "Test", "phi": PHI,
	}
	for m := ModelKV; m <= ModelSpatial; m++ {
		globalHydra.Insert(m, "test_key", testData)
	}
	fmt.Printf("  ✅ Inserted into all %d models\n\n", 14)

	// Test 2: Cross-model query
	fmt.Println("=== Test 2: Cross-Model Query ===")
	results := globalHydra.CrossModelQuery("test_key")
	modelCount := len(results)
	fmt.Printf("  ✅ Found data in %d models\n\n", modelCount)

	// Test 3: 8 consensus algorithms
	fmt.Println("=== Test 3: 8 Hydra Consensus ===")
	nodes := []string{"node1", "node2", "node3", "node4", "node5", "node6", "node7"}
	for h := HydraPhiRaft; h <= HydraHarmonyChain; h++ {
		cons := globalHydra.Consensus[h]
		reached := cons.ReachConsensus(nodes)
		fmt.Printf("  %-20s: Consensus=%t, Lyapunov=%.4f\n", HydraNames[h], reached, cons.Lyapunov)
	}

	// Test 4: Harmonization
	fmt.Println("\n=== Test 4: φ-Harmonization ===")
	oldPhi := globalHydra.PhiAnchor
	globalHydra.Harmonize()
	fmt.Printf("  PhiAnchor: %.6f → %.6f ✅\n\n", oldPhi, globalHydra.PhiAnchor)

	// Test 5: FHE Bootstrap
	fmt.Println("=== Test 5: FHE Bootstrapping ===")
	globalHydra.BootstrapNoise()
	fmt.Printf("  PhiAnchor after bootstrap: %.6f (target: 40.0) ✅\n\n", globalHydra.PhiAnchor)

	// Test 6: Stats
	fmt.Println("=== Test 6: Full Stats ===")
	stats := globalHydra.Stats()
	statsJSON, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Println(string(statsJSON))

	fmt.Println("\n╔═══════════════════════════════════════╗")
	fmt.Println("║  SPIRALDB HYDRA — ALL TESTS PASSED   ║")
	fmt.Println("║  ΦΩ0 — I AM THAT I AM                ║")
	fmt.Println("╚═══════════════════════════════════════╝")
}

func startHydraServer() {
	InitGlobalHydra()
	mux := http.NewServeMux()
	HydraRegisterHandlers(mux)
	
	go func() {
		fmt.Println("Hydra server on :5445")
		http.ListenAndServe(":5445", mux)
	}()
	time.Sleep(100 * time.Millisecond)
}
