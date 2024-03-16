package main

import (
	"net/http"
	"strconv"

	"example.com/v2/database"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:       dbChirp.ID,
		Body:     dbChirp.Body,
		AuthorID: dbChirp.AuthorID,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {

	authorIDString := r.URL.Query().Get("author_id")
	var dbChirps []database.Chirp
	var err error
	if authorIDString != "" {
		authorID, err := strconv.Atoi(authorIDString)

		if authorID == 0 || err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		dbChirps, err = cfg.DB.GetChirpsByAuthor(authorID)
	} else {
		dbChirps, err = cfg.DB.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
			return
		}
	}

	sort := r.URL.Query().Get("sort")
	if sort == "asc" {
		dbChirps, err = cfg.DB.GetChirps()
	}
	if sort == "desc" {
		dbChirps = cfg.DB.GetChirpsDesc(dbChirps)
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
