package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	c, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/assets/ready", func(w http.ResponseWriter, r *http.Request) {
		var job AssetJob
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := deliver(r.Context(), c, job); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, `{"status":"notification-published"}`)
	})
	log.Println("creator notification service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
