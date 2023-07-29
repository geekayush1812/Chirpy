package database

import (
	"errors"
	"reflect"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	RevokedRefreshTokens map[string]string `json:"revokedRefreshTokens"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func (db *DB) IsUserExists(userId int) (error, bool) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()

	if err != nil {
		return err, false
	}

	_ , ok := dbData.Users[userId];

	return nil, ok
}

func (db *DB) CreateUser(email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	for _, user := range dbData.Users {
		if user.Email == email {
			return User{}, errors.New("user with this email already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)

	if err!= nil {
		return User{}, err
	}

	newUser := User{
		Id:    len(dbData.Users) + 1,
		Email: email,
		Password: string(hashedPassword),
		RevokedRefreshTokens: map[string]string{},
		IsChirpyRed: false,
	}

	dbData.Users[newUser.Id] = newUser
	err = db.writeDB(dbData)

	// setting password to zero value
	newUser.Password = ""
	newUser.RevokedRefreshTokens = make(map[string]string)

	return newUser, err
}

func (db *DB) GetUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	for _, user := range dbData.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user does not exists")
}

func (db *DB) UpdateUser(userId int, email, password string, isChirpyRed bool) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	newUser, ok := dbData.Users[userId]

	if !ok {
		return User{}, errors.New("something went wrong")
	}

	if email != "" {
		newUser.Email = email
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)

		if err != nil {
			return User{}, errors.New("something went wrong")
		}

		newUser.Password = string(hashedPassword)
	}

	if isChirpyRed {
		newUser.IsChirpyRed = isChirpyRed
	}

	dbData.Users[userId] = newUser

	err = db.writeDB(dbData)

	if err != nil {
		return User{}, err
	}

	// setting password to zero value
	newUser.Password = ""
	newUser.RevokedRefreshTokens = make(map[string]string)

	return newUser, nil
}

func (db *DB) IsUserRefreshTokenRevoked(userId int, userToken string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()

	if err != nil {
		return err
	}

	newUser := User{}

	for _, user := range dbData.Users {
		if user.Id == userId {
			newUser = user
			break;
		}
	}

	if reflect.DeepEqual(newUser, User{}) {
		return errors.New("something went wrong")
	}

	for token := range newUser.RevokedRefreshTokens {
		if token == userToken {
			return errors.New("token is revoked")
		}
	}

	return nil
}

func (db *DB) RevokeUserRefreshToken(userId int, userToken string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()

	if err != nil {
		return err
	}

	for _, user := range dbData.Users {
		if user.Id == userId {
			user.RevokedRefreshTokens[userToken] = time.Now().String()
			dbData.Users[userId] = user;
		
			err = db.writeDB(dbData)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("something went wrong")
}