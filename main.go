package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type WebhookPayload struct {
	Secret      string  `json:"secret"`
	Env         string  `json:"env"`
	Strategy    string  `json:"strategy"`
	SignalID    string  `json:"signal_id"`
	TVSymbol    string  `json:"tv_symbol"`
	Side        string  `json:"side"`
	Qty         int     `json:"qty"`
	RefPrice    float64 `json:"ref_price"`
	SLPoints    float64 `json:"sl_points"`
	TPPoints    float64 `json:"tp_points"`
	TimestampMs int64   `json:"timestamp_ms"`
}

func main() {
	sharedSecret := os.Getenv("TV_SHARED_SECRET")
	if sharedSecret == "" {
		sharedSecret = "CHANGE_ME_SHARED_SECRET"
	}

	http.HandleFunc("/tv-webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var p WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if p.Secret != sharedSecret {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if p.SignalID == "" || p.Strategy == "" || p.TVSymbol == "" {
			http.Error(w, "missing required fields", http.StatusBadRequest)
			return
		}

		if p.Side != "buy" && p.Side != "sell" {
			http.Error(w, "invalid side", http.StatusBadRequest)
			return
		}

		if p.Qty <= 0 {
			http.Error(w, "invalid qty", http.StatusBadRequest)
			return
		}

		log.Printf(
			"[TV] env=%s strategy=%s signal=%s %s %d %s @ %.2f SL=%v TP=%v",
			p.Env,
			p.Strategy,
			p.SignalID,
			p.Side,
			p.Qty,
			p.TVSymbol,
			p.RefPrice,
			p.SLPoints,
			p.TPPoints,
		)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
