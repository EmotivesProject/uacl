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
	"uacl/pkg/password"
	"uacl/pkg/uacl_errors"

	"github.com/go-chi/chi"
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
	messageResponseJSON(w, http.StatusOK, model.Message{Message: "Health ok"})
}

func publicKey(w http.ResponseWriter, r *http.Request) {
	public, err := ioutil.ReadFile(publicKeyLocation)
	if err != nil {
		messageResponseJSON(w, http.StatusInternalServerError, model.Message{Message: "Failed to find key"})
		return
	}

	resultResponseJSON(w, http.StatusOK, model.Key{
		Key: string(public),
	})
}

func getUserByEncodedID(w http.ResponseWriter, r *http.Request) {
	encodedID := chi.URLParam(r, "encoded_id")
	user, err := db.FindByEncodedID(encodedID, database)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "Failed to find user"})
		return
	}

	resultResponseJSON(w, http.StatusOK, model.ShortenedUser{
		Name:      user.Name,
		Email:     user.Email,
		EncodedID: user.EncodedID,
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: errFailedDecoding.Error()})
		return
	}

	target, err := user.ValidateLogin()
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, model.Message{
			Message: err.Error(),
			Target:  target,
		})
		return
	}

	databaseUser, err := db.FindByEmail(user.Email, database)
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, model.Message{Message: err.Error()})
		return
	}

	correctPassword := password.ValidatePassword(user.Password, databaseUser.Password)
	if !correctPassword {
		messageResponseJSON(w, http.StatusUnprocessableEntity, model.Message{Message: uacl_errors.ErrInvalidCredentials.Error()})
		return
	}

	passTokenToUser(w, &databaseUser)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: errFailedDecoding.Error()})
		return
	}

	target, err := user.ValidateCreate()
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, model.Message{
			Message: err.Error(),
			Target:  target,
		})
		return
	}

	encryptedPassword := password.CreatePassword(user.Password)
	if encryptedPassword == "" {
		messageResponseJSON(w, http.StatusUnprocessableEntity, model.Message{Message: errFailedCrypting.Error()})
		return
	}
	user.Password = encryptedPassword

	createdUser := database.Create(user)
	if createdUser.Error != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, model.Message{Message: createdUser.Error.Error()})
		return
	}

	encodedID, err := encode.GenerateBase64ID(encodedIDLength, encodePrefix)
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, model.Message{Message: err.Error()})
		return
	}
	encodedID = encodedID[:len(encodedID)-1]
	user.EncodedID = encodedID
	database.Save(user)

	passTokenToUser(w, user)
}

func passTokenToUser(w http.ResponseWriter, user *model.User) {
	tokenString, err := auth.CreateToken(*user)
	if err != nil {
		messageResponseJSON(w, http.StatusUnprocessableEntity, model.Message{Message: err.Error()})
		return
	}

	token := model.Token{
		Token: tokenString,
	}

	resultResponseJSON(w, http.StatusCreated, token)
}
