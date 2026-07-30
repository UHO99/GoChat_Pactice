package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type roomResponse struct {
	Name string `json:"name"`
}

type userResponse struct {
	Nickname string `json:"nickname"`
}

type createRoomRequest struct {
	Name string `json:"name"`
}

type createUserRequest struct {
	Nickname string `json:"name"`
}

const pgUniqueViolation = "23505"

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Nickname == "" {
		http.Error(w, "Nickname parameter Error", http.StatusBadRequest)
		return
	}

	user, err := s.store.CreateUser(r.Context(), req.Nickname)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			http.Error(w, "room already exists", http.StatusConflict)
			return
		}
		http.Error(w, "cannot create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse{Nickname: user.Nickname})
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "room name required", http.StatusBadRequest)
		return
	}

	room, err := s.store.CreateRoom(r.Context(), req.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			http.Error(w, "room already exists", http.StatusConflict)
			return
		}
		http.Error(w, "cannot create room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(roomResponse{Name: room.Name})
}

func (s *Server) handleListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := s.store.ListRooms(r.Context())
	if err != nil {
		http.Error(w, "cannot list rooms", http.StatusInternalServerError)
		return
	}

	resp := make([]roomResponse, 0, len(rooms))
	for _, room := range rooms {
		resp = append(resp, roomResponse{Name: room.Name})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
