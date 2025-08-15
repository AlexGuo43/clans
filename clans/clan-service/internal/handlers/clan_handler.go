package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AlexGuo43/clans/clan-service/internal/models"
	"github.com/AlexGuo43/clans/clan-service/internal/services"
	"github.com/gorilla/mux"
)

type ClanHandler struct {
	clanService *services.ClanService
}

func NewClanHandler(clanService *services.ClanService) *ClanHandler {
	return &ClanHandler{
		clanService: clanService,
	}
}

func (h *ClanHandler) CreateClan(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ClanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	clan, err := h.clanService.CreateClan(r.Context(), &req, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(clan)
}

func (h *ClanHandler) GetClan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid clan ID", http.StatusBadRequest)
		return
	}

	clan, err := h.clanService.GetClan(r.Context(), id)
	if err != nil {
		http.Error(w, "Clan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clan)
}

func (h *ClanHandler) GetClanByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	clan, err := h.clanService.GetClanByName(r.Context(), name)
	if err != nil {
		http.Error(w, "Clan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clan)
}

func (h *ClanHandler) GetClans(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	clans, err := h.clanService.GetClans(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clans)
}

func (h *ClanHandler) UpdateClan(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid clan ID", http.StatusBadRequest)
		return
	}

	var req models.ClanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	clan, err := h.clanService.UpdateClan(r.Context(), id, &req, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clan)
}

func (h *ClanHandler) DeleteClan(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid clan ID", http.StatusBadRequest)
		return
	}

	err = h.clanService.DeleteClan(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ClanHandler) JoinClan(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid clan ID", http.StatusBadRequest)
		return
	}

	err = h.clanService.JoinClan(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Successfully joined clan"}`))
}

func (h *ClanHandler) LeaveClan(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid clan ID", http.StatusBadRequest)
		return
	}

	err = h.clanService.LeaveClan(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Successfully left clan"}`))
}

func (h *ClanHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid clan ID", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	members, err := h.clanService.GetMembers(r.Context(), id, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func (h *ClanHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	clanID, err := strconv.Atoi(vars["clanId"])
	if err != nil {
		http.Error(w, "Invalid clan ID", http.StatusBadRequest)
		return
	}

	targetUserID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.ClanMembershipUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.clanService.UpdateMemberRole(r.Context(), clanID, targetUserID, userID, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Role updated successfully"}`))
}

func (h *ClanHandler) GetUserClans(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	clans, err := h.clanService.GetUserClans(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clans)
}

func (h *ClanHandler) GetMembership(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid clan ID", http.StatusBadRequest)
		return
	}

	membership, err := h.clanService.GetMembership(r.Context(), id, userID)
	if err != nil {
		http.Error(w, "Not a member", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(membership)
}

func getUserIDFromHeader(r *http.Request) int {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0
	}
	return userID
}