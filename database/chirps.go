package database

import (
	"errors"
	"sort"
	"strconv"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
	AuthorId int `json:"author_id"`
}


func sortChirps(chirps []Chirp, sortType string) {
	sort.Slice(chirps, func(i, j int) bool {
		if sortType == "asc" {
			return chirps[i].Id < chirps[j].Id
		} else {
			return chirps[i].Id > chirps[j].Id
		}
	})
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(authorId int, body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	newChirp := Chirp{
		Id:   len(dbData.Chirps) + 1,
		Body: body,
		AuthorId: authorId,
	}

	dbData.Chirps[newChirp.Id] = newChirp

	err = db.writeDB(dbData)

	return newChirp, err
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps(authorId string, sortType string) ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbData, err := db.loadDB()

	if err != nil {
		return []Chirp{}, err
	}

	chirps := []Chirp{}

	for _, chirp := range dbData.Chirps {
		if authorId != "" {
			authorId, err := strconv.ParseInt(authorId, 10, 32)

			if err != nil {
				return []Chirp{}, err
			}

			if chirp.AuthorId == int(authorId) {
				chirps = append(chirps, chirp)
			}
		} else {
			chirps = append(chirps, chirp)
		}
	}

	if (sortType == "") {
		sortType = "asc"
	}

	sortChirps(chirps, sortType)

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbData, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbData.Chirps[id]

	if !ok {
		return Chirp{}, errors.New("Chirp not found")
	}

	return chirp, nil
}


func (db *DB) DeleteChirp(id int) error {
	db.mux.Lock()

	defer db.mux.Unlock()

	dbData, err := db.loadDB()

	if err != nil {
		return err
	}

	if _, ok := dbData.Chirps[id]; !ok {
		return errors.New("internal server error")
	}

	return nil
}