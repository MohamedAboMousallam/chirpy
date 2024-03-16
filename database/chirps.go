package database

import "errors"

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, userid int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: userid,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(chirpID, userID int) error {

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Check if chirp exists
	chirp, ok := dbStructure.Chirps[chirpID]
	if !ok {
		return ErrNotExist
	}

	// Check if the user is the author of the chirp
	if chirp.AuthorID != userID {
		return errors.New("user is not the author of the chirp")
	}

	// Delete the chirp
	delete(dbStructure.Chirps, chirpID)

	// Update the database file
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetChirpsByAuthor(authorID int) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0)
	for _, chirp := range dbStructure.Chirps {
		if chirp.AuthorID == authorID {
			chirps = append(chirps, chirp)
		}
	}

	return chirps, nil
}

func (db *DB) GetChirpsDesc(chirps []Chirp) []Chirp {
	for i := 0; i < len(chirps); i++ {
		for j := i + 1; j < len(chirps); j++ {
			if chirps[i].ID < chirps[j].ID {
				chirps[i], chirps[j] = chirps[j], chirps[i]
			}
		}
	}
	return chirps
}
