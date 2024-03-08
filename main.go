package main

import (
	"UTS/controllers"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// ENDPOINT 1 : List Rooms berdasarkan Game
	router.HandleFunc("/rooms", controllers.GetAllRooms).Methods("GET")

	// ENDPOINT 2 : Get Room Details
	router.HandleFunc("/roomDetails", controllers.GetDetailRoom).Methods("GET")

	// ENDPOINT 3 : Insert Room
	router.HandleFunc("/participantInsert", controllers.InsertParticipant).Methods("POST")

	// ENDPOINT 4 : Leave Room
	router.HandleFunc("/participantLeave", controllers.LeaveParticipant).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
