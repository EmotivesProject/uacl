package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"uacl/internal/db"
	"uacl/model"

	"golang.org/x/crypto/bcrypt"
)

var (
	database          = db.ConnectDB()
	errFailedDecoding = errors.New("Failed during decoding request")
	errFailedCrypting = errors.New("Failed during encrypting password")
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, "Health ok")
}

func FetchItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(userID)
	fmt.Printf("Fetching information for %+v", ctx)
	messageResponseJSON(w, http.StatusOK, "CONNECTED")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, errFailedDecoding.Error())
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, errFailedCrypting.Error())
		return
	}
	user.Password = string(pass)

	createdUser := database.Create(user)
	if createdUser.Error != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, createdUser.Error.Error())
		return
	}

	resultResponseJSON(w, http.StatusCreated, createdUser)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, errFailedDecoding.Error())
		return
	}
	resp, err := db.FindOne(user.Email, user.Password, database)
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	resultResponseJSON(w, http.StatusCreated, resp)
}
