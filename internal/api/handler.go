package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"uacl/internal/db"
	"uacl/model"
	"uacl/pkg/auth"
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

// Should also refresh if required
func authorizeHeader(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("authorization")
	header = strings.Split(header, "Bearer ")[1]

	user, err := auth.Validate(header)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{
			Message: errUnauthorised.Error(),
		})
		return
	}
	fmt.Println("AUTHORIZING" + user.Username)
	resultResponseJSON(w, http.StatusOK, user)
}

func getUserByEncodedID(w http.ResponseWriter, r *http.Request) {
	encodedID := chi.URLParam(r, "username")
	user, err := db.FindByUsername(encodedID, database)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "Failed to find user"})
		return
	}

	resultResponseJSON(w, http.StatusOK, model.ShortenedUser{
		Name:     user.Name,
		Username: user.Username,
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

	databaseUser, err := db.FindByUsername(user.Username, database)
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

	sendUserToPostit(user)

	passTokenToUser(w, user)
}

func sendUserToPostit(user *model.User) {
	baseHost := os.Getenv("BASE_HOST")
	postitURL := baseHost + "postit/user"

	// Clear the id since postit will create a new one
	user.ID = 0
	requestBody, err := json.Marshal(user)
	if err != nil {
		fmt.Println("FAILED TO SEND NEW USER")
	}

	_, err = http.Post(postitURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("FAILED TO SEND NEW USER")
	}
	fmt.Println("Sent user to postit")
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
