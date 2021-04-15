package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"uacl/internal/auth"
	"uacl/internal/db"
	"uacl/internal/password"
	"uacl/internal/uacl_errors"
	"uacl/model"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/response"

	"github.com/go-chi/chi"
)

const (
	publicKeyLocation = "/app/jwt/public.key"
)

var (
	errFailedDecoding = errors.New("Failed during decoding request")
	errFailedCrypting = errors.New("Failed during encrypting password")
	errUnauthorised   = errors.New("Unauthorized")
)

func healthz(w http.ResponseWriter, r *http.Request) {
	response.MessageResponseJSON(w, http.StatusOK, response.Message{Message: "Health ok"})
}

func publicKey(w http.ResponseWriter, r *http.Request) {
	public, err := ioutil.ReadFile(publicKeyLocation)
	logger.Info("Fetching public key file")
	if err != nil {
		response.MessageResponseJSON(w, http.StatusInternalServerError, response.Message{Message: "Failed to find key"})
		return
	}

	response.ResultResponseJSON(w, http.StatusOK, model.Key{
		Key: string(public),
	})
}

// Should also refresh if required
func authorizeHeader(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	header = strings.Split(header, "Bearer ")[1]

	user, err := auth.Validate(header)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{
			Message: errUnauthorised.Error(),
		})
		return
	}
	logger.Infof("Validating %s", user.Username)
	response.ResultResponseJSON(w, http.StatusOK, user)
}

func getUserByEncodedID(w http.ResponseWriter, r *http.Request) {
	encodedID := chi.URLParam(r, "username")
	user, err := db.FindByUsername(encodedID, db.GetDB())
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: "Failed to find user"})
		return
	}

	logger.Infof("Fetched user %s", user.Username)
	response.ResultResponseJSON(w, http.StatusOK, model.ShortenedUser{
		Name:     user.Name,
		Username: user.Username,
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: errFailedDecoding.Error()})
		return
	}

	target, err := user.ValidateLogin()
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{
			Message: err.Error(),
			Target:  target,
		})
		return
	}

	databaseUser, err := db.FindByUsername(user.Username, db.GetDB())
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})
		return
	}

	correctPassword := password.ValidatePassword(user.Password, databaseUser.Password)
	if !correctPassword {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: uacl_errors.ErrInvalidCredentials.Error()})
		return
	}
	logger.Infof("Logging in user %s", user.Username)

	passTokenToUser(w, &databaseUser)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: errFailedDecoding.Error()})
		return
	}

	target, err := user.ValidateCreate()
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{
			Message: err.Error(),
			Target:  target,
		})
		return
	}

	encryptedPassword := password.CreatePassword(user.Password)
	if encryptedPassword == "" {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: errFailedCrypting.Error()})
		return
	}
	user.Password = encryptedPassword

	createdUser := db.GetDB().Create(user)
	if createdUser.Error != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: createdUser.Error.Error()})
		return
	}
	logger.Infof("Created user %s", user.Username)

	passTokenToUser(w, user)
}

func passTokenToUser(w http.ResponseWriter, user *model.User) {
	tokenString, err := auth.CreateToken(*user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})
		return
	}

	token := model.Token{
		Token: tokenString,
	}

	response.ResultResponseJSON(w, http.StatusCreated, token)
}
