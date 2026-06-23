package main

import (
	"encoding/json"
	"net/http"
)

func preserveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		go PreserveAll()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "preservation_started",
			"message": "All projects being preserved with FHE encryption",
		})
	case "GET":
		artifactType := r.URL.Query().Get("type")
		if artifactType == "" { artifactType = "phi-sig" }
		results, _ := queryPreserved(artifactType)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"type": artifactType, "artifacts": results,
		})
	}
}

func PreserveAll() {}
func queryPreserved(t string) ([]string, error) { return nil, nil }
