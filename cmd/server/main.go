package main

import (
	"customerfollowup/internal/api"
	"customerfollowup/internal/audit"
	"customerfollowup/internal/service"
	"customerfollowup/internal/store"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := store.Open("followup.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	svc := service.New(db, audit.New(db))
	mux := api.New(svc)
	fmt.Println("customer follow-up service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
