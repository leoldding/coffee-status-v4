package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Output struct {
	Item map[string]map[string]string
}

type Status struct {
	Value string `json:"status"`
}

// server used during dev
func main() {
	var status string
	status = "no"

	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")

	// return the current status
	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		output := Output{Item: make(map[string]map[string]string)}
		output.Item["value"] = make(map[string]string)
		output.Item["value"]["S"] = status
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(output)
	}).Methods("GET")

	// update status to new status
	router.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		var newStatus Status
		if err := json.NewDecoder(r.Body).Decode(&newStatus); err != nil {
			log.Printf("Error decoding JSON body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		status = newStatus.Value
	}).Methods("POST")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
