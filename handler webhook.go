package main

import (
	"encoding/json"
	"errors"
	"net/http"
    "os"
	"example.com/v2/database"
)


func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {

    expectedAPIKey := os.Getenv("API_KEY")

    apiKeyHeader := r.Header.Get("Authorization")
    if apiKeyHeader == "" || apiKeyHeader != "ApiKey "+expectedAPIKey {
        respondWithError(w, http.StatusUnauthorized, "Invalid API key")
        return
    }

    type WebhookEvent struct {
        Event string `json:"event"`
        Data  struct {
            UserID int `json:"user_id"`
        } `json:"data"`
    }

    var webhookEvent WebhookEvent
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&webhookEvent)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if webhookEvent.Event != "user.upgraded" {
        respondWithJSON(w, http.StatusOK, map[string]string{"message": "Event not processed"})
        return
    }

    // Event is "user.upgraded", update user in the database
    user, err := cfg.DB.GetUser(webhookEvent.Data.UserID)
    if err != nil {
        if errors.Is(err, database.ErrNotExist) {
            respondWithError(w, http.StatusNotFound, "User not found")
        } else {
            respondWithError(w, http.StatusInternalServerError, "Error updating user")
        }
        return
    }

    // Mark the user as Chirpy Red
    user.IsChirpyRed = true
    _, err = cfg.DB.UpdateUser(user.ID, user.Email, user.HashedPassword, user.IsChirpyRed)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Error updating user")
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"message": "User upgraded successfully"})
}
