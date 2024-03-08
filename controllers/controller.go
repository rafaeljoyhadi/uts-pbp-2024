package controllers

import (
	m "UTS/models"
	// "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "strconv"
	// "strings"

	// * UNTUK ENDPOINT 2
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// NO 1
func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		http.Error(w, "Please input the game id", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`
        SELECT r.id, r.room_name 
        FROM rooms r
        INNER JOIN games g ON r.id_game = g.id
        WHERE g.id = %s`, gameID)

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch rooms", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rooms []m.Room

	for rows.Next() {
		var room m.Room
		if err := rows.Scan(&room.ID, &room.RoomName); err != nil {
			log.Println(err)
			http.Error(w, "Failed to scan room data", http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, room)
	}

	response := struct {
		Status int      `json:"status"`
		Data   []m.Room `json:"data"`
	}{
		Status: 200,
		Data:   rooms,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// NO 2
func GetDetailRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	roomID := vars["room_id"]

	query := `
        SELECT r.id, r.room_name, p.id AS participant_id, p.id_account
        FROM rooms r
        LEFT JOIN participants p ON r.id = p.id_room
        WHERE r.id = ?`

	rows, err := db.Query(query, roomID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch room details", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var room m.Room
	var participants []m.Participant

	for rows.Next() {
		var participant m.Participant
		if err := rows.Scan(&room.ID, &room.RoomName, &participant.ID, &participant.AccountID); err != nil {
			log.Println(err)
			http.Error(w, "Failed to scan room data", http.StatusInternalServerError)
			return
		}
		participants = append(participants, participant)
	}

	response := struct {
		Status int              `json:"status"`
		Data   m.RoomDetailData `json:"data"`
	}{
		Status: 200,
		Data: m.RoomDetailData{
			Room:         room,
			Participants: participants,
		},
	}

	// Set response headers and encode the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// NO 3
func InsertParticipant(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	roomID := r.URL.Query().Get("room_id")
	accountID := r.URL.Query().Get("account_id")
	var response m.InsertResponse

	row := db.QueryRow("SELECT G.max_player FROM Games G JOIN Rooms R ON G.id = R.id_game WHERE R.id = ? ", roomID)
	var maxPlayer int
	err := row.Scan(&maxPlayer)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	countRow := db.QueryRow("SELECT COUNT(*) FROM participants WHERE id_room = ?", roomID)
	var participantCount int
	err = countRow.Scan(&participantCount)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if participantCount >= maxPlayer {
		response.Status = 400
		response.Message = "Room is full"
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = db.Exec("INSERT INTO participants (id_room, id_account) VALUES (?, ?)", roomID, accountID)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response.Status = 200
	response.Message = "Participant inserted to room successfully"
	json.NewEncoder(w).Encode(response)
}

// NO 4
func LeaveParticipant(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	roomID := r.URL.Query().Get("room_id")
	accountID := r.URL.Query().Get("account_id")
	var response m.InsertResponse

	_, err := db.Exec("DELETE FROM participants WHERE id_room = ? AND id_account = ?", roomID, accountID)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response.Status = 200
	response.Message = "Participant left the room successfully"
	json.NewEncoder(w).Encode(response)
}
