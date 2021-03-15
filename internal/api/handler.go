package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"uacl/internal/db"
	"uacl/model"
	"uacl/pkg/auth"
	"uacl/pkg/encode"

	"golang.org/x/crypto/bcrypt"
)

const (
	encodePrefix      = "USER"
	publicKeyLocation = "/app/jwt/public.key"
	encodedIDLength   = 9
)

var (
	database          = db.ConnectDB()
	errFailedDecoding = errors.New("Failed during decoding request")
	errFailedCrypting = errors.New("Failed during encrypting password")
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, "Health ok")
}

func publicKey(w http.ResponseWriter, r *http.Request) {
	public, err := ioutil.ReadFile(publicKeyLocation)
	if err != nil {
		messageResponseJSON(w, http.StatusInternalServerError, "Failed to find key")
		return
	}

	resultResponseJSON(w, http.StatusOK, model.Key{
		Key: string(public),
	})
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

	encodedID, err := encode.GenerateBase64ID(encodedIDLength, encodePrefix)
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	encodedID = encodedID[:len(encodedID)-1]
	user.EncodedID = encodedID
	database.Save(user)

	token, err := auth.CreateToken(*user)
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	resultResponseJSON(w, http.StatusCreated, token)
}
