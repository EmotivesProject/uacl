package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"uacl/internal/auth"
	"uacl/internal/db"
	"uacl/internal/password"
	"uacl/messages"
	"uacl/model"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/response"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	response.MessageResponseJSON(w, http.StatusOK, response.Message{Message: messages.HealthResponse})
}

func publicKey(w http.ResponseWriter, r *http.Request) {
	public, err := ioutil.ReadFile(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, http.StatusOK, model.Key{
		Key: string(public),
	})
}

func authorizeHeader(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	header = strings.Split(header, "Bearer ")[1]

	user, err := auth.Validate(header)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{
			Message: messages.ErrUnauthorised.Error(),
		})

		return
	}

	logger.Infof("Validating %s", user.Username)
	response.ResultResponseJSON(w, http.StatusOK, user)
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	token := model.Token{}

	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	user, err := auth.Validate(token.RefreshToken)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	if user.Username != token.Username {
		logger.Error(messages.ErrMismatchUsername)
		response.MessageResponseJSON(w, http.StatusInternalServerError, response.Message{
			Message: messages.ErrMismatchUsername.Error(),
		})

		return
	}

	if !db.RefreshTokenIsValidForUsername(token.RefreshToken, token.Username) {
		logger.Error(messages.ErrWrongRefreshToken)
		response.MessageResponseJSON(w, http.StatusInternalServerError, response.Message{
			Message: messages.ErrWrongRefreshToken.Error(),
		})

		return
	}

	passTokenToUser(w, &model.User{
		Name:     user.Name,
		Username: user.Username,
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

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

	databaseUser, err := db.FindByUsername(user.Username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})

		return
	}

	correctPassword := password.ValidatePassword(user.Password, databaseUser.Password)
	if !correctPassword {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{
			Message: messages.ErrInvalidCredentials.Error(),
		})

		return
	}

	logger.Infof("Logging in user %s", user.Username)

	passTokenToUser(w, &databaseUser)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

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
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})

		return
	}

	user.Password = encryptedPassword

	err = db.CreateNewUser(user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Created user %s", user.Username)

	passTokenToUser(w, user)
}

func passTokenToUser(w http.ResponseWriter, user *model.User) {
	tokenString, err := auth.CreateToken(*user, false)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})

		return
	}

	refreshTokenString, err := auth.CreateToken(*user, true)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})

		return
	}

	token := model.Token{
		Username:     user.Username,
		Token:        tokenString,
		RefreshToken: refreshTokenString,
		UpdatedAt:    time.Now(),
	}

	err = db.UpsertToken(&token)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, http.StatusCreated, token)
}
