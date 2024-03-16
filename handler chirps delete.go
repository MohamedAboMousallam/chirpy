package main

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/v2/database"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	// Parse chirpID from the URL parameters
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// Authenticate the user and get the userID
	userID, err := authenticateUser(r, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication failed")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
		} else if errors.Is(err, errors.New("user is not the author of the chirp")) {
			respondWithError(w, http.StatusForbidden, "User is not the author of the chirp")
		} else {
			respondWithError(w, http.StatusForbidden, "Error deleting chirp")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Chirp deleted successfully"})
}
