package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
	"uacl/internal/auth"
	"uacl/internal/db"
	"uacl/internal/password"
	"uacl/internal/send"
	"uacl/messages"
	"uacl/model"

	"github.com/EmotivesProject/common/logger"
	"github.com/EmotivesProject/common/response"
	"github.com/go-chi/chi"
)

const (
	autologinLength       = 64
	successfulHeaderCount = 2
)

func authorizeHeader(w http.ResponseWriter, r *http.Request) {
	user, err := doAuthentication(r)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{
			Message: messages.ErrUnauthorised.Error(),
		})

		return
	}

	logger.Infof("Authorization request for %s", user.Username)
	response.ResultResponseJSON(w, false, http.StatusOK, user)
}

func doAuthentication(r *http.Request) (model.ShortenedUser, error) {
	header := r.Header.Get("Authorization")
	headerSplit := strings.Split(header, "Bearer ")

	if len(headerSplit) < successfulHeaderCount {
		return model.ShortenedUser{}, messages.ErrUnauthorised
	}

	return auth.Validate(headerSplit[1])
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	token := model.Token{}

	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	user, err := auth.Validate(token.RefreshToken)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: err.Error()})

		return
	}

	if !db.RefreshTokenIsValidForUsername(r.Context(), token.RefreshToken, user.Username) {
		logger.Error(messages.ErrWrongRefreshToken)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{
			Message: messages.ErrWrongRefreshToken.Error(),
		})

		return
	}

	logger.Infof("Created refresh token for %s", user.Username)

	passTokenToUser(r.Context(), w, &model.User{
		Name:      user.Name,
		Username:  user.Username,
		UserGroup: user.UserGroup,
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	user.Username = strings.ToLower(user.Username)
	user.UserGroup = strings.ToLower(user.UserGroup)

	target, err := user.ValidateLogin()
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{
			Message: err.Error(),
			Target:  target,
		})

		return
	}

	databaseUser, err := db.FindByUsername(r.Context(), user.Username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: err.Error()})

		return
	}

	correctPassword := password.ValidatePassword(user.Password, databaseUser.Password)
	if !correctPassword {
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{
			Message: messages.ErrInvalidCredentials.Error(),
		})

		return
	}

	logger.Infof("Logging in user %s", user.Username)

	passTokenToUser(r.Context(), w, &databaseUser)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	user.Username = strings.ToLower(user.Username)
	user.UserGroup = strings.ToLower(user.UserGroup)

	target, err := user.ValidateCreate()
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{
			Message: err.Error(),
			Target:  target,
		})

		return
	}

	encryptedPassword := password.CreatePassword(user.Password)
	if encryptedPassword == "" {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	user.Password = encryptedPassword

	err = db.CreateNewUser(r.Context(), user)
	if err != nil {
		logger.Error(err)

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			response.MessageResponseJSON(
				w, false, http.StatusUnprocessableEntity, response.Message{Message: "Username is already used"},
			)

			return
		}

		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Created user %s", user.Username)

	err = send.ChatterUser(user)
	if err != nil {
		logger.Error(err)
	}

	passTokenToUser(r.Context(), w, user)
}

func getAutologinTokens(w http.ResponseWriter, r *http.Request) {
	authUser, err := doAuthentication(r)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: err.Error()})

		return
	}

	authorizedUsers := strings.Split(os.Getenv("AUTOLOGIN_CREATE_USERS"), ",")

	in := stringInSlice(authUser.Username, authorizedUsers)
	if !in {
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: "no authorized"})

		return
	}

	autologins, err := db.FindAutologins(r.Context())
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Info("Fetched autologins")

	response.ResultResponseJSON(w, false, http.StatusCreated, autologins)
}

func getAutologinToken(w http.ResponseWriter, r *http.Request) {
	authUser, err := doAuthentication(r)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: err.Error()})

		return
	}

	autologinID, err := extractID(r, "token_id")
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	autoLogin, err := db.FindAutologinByID(r.Context(), autologinID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	authorizedUsers := strings.Split(os.Getenv("AUTOLOGIN_CREATE_USERS"), ",")

	inAuthorizedUsers := stringInSlice(authUser.Username, authorizedUsers)
	isUsernameEqual := autoLogin.Username == authUser.Username

	if !inAuthorizedUsers && !isUsernameEqual {
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: "no authorized"})

		return
	}

	logger.Infof("Fetched autologin for %s", authUser.Username)

	response.ResultResponseJSON(w, false, http.StatusCreated, autoLogin)
}

func createLoginToken(w http.ResponseWriter, r *http.Request) {
	authUser, err := doAuthentication(r)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: err.Error()})

		return
	}

	user := &model.AutologinRequest{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	user.Username = strings.ToLower(user.Username)

	authorizedUsers := strings.Split(os.Getenv("AUTOLOGIN_CREATE_USERS"), ",")

	in := stringInSlice(authUser.Username, authorizedUsers)
	isUsernameEqual := user.Username == authUser.Username

	if !in && !isUsernameEqual {
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: "no authorized"})

		return
	}

	dbUser, err := db.FindByUsername(r.Context(), user.Username)
	if err != nil {
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	id, err := generateRandomString(autologinLength)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	err = db.CreateNewAutologinToken(r.Context(), dbUser.Username, id)
	if err != nil {
		logger.Error(err)
		// assuming error with db is missing value
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	auto := model.AutologinToken{
		Username:       dbUser.Username,
		AutologinToken: id,
		Site:           os.Getenv("AUTOLOGIN_URL"),
	}

	logger.Infof("Created autologin for user %s %s", dbUser.Username, id)

	response.ResultResponseJSON(w, false, http.StatusCreated, auto)
}

func deleteAutologinToken(w http.ResponseWriter, r *http.Request) {
	authUser, err := doAuthentication(r)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: err.Error()})

		return
	}

	authorizedUsers := strings.Split(os.Getenv("AUTOLOGIN_CREATE_USERS"), ",")

	in := stringInSlice(authUser.Username, authorizedUsers)
	if !in {
		response.MessageResponseJSON(w, false, http.StatusUnauthorized, response.Message{Message: "no authorized"})

		return
	}

	autologinToken := chi.URLParam(r, "token")

	err = db.DeleteAutologinToken(r.Context(), autologinToken)
	if err != nil {
		logger.Error(err)
		// assuming error with db is missing value
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Deleted autologin %s", autologinToken)

	response.ResultResponseJSON(w, false, http.StatusOK, response.Message{Message: "Autologin Deleted"})
}

func authoriseLoginToken(w http.ResponseWriter, r *http.Request) {
	autologinToken := chi.URLParam(r, "token")

	autoLoginRequest, err := db.FindAutologinForUser(r.Context(), autologinToken)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	dbUser, err := db.FindByUsername(r.Context(), autoLoginRequest.Username)
	if err != nil {
		logger.Error(err)
		// assuming error with db is missing value
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Logging in user via autologin %s", dbUser.Username)

	passTokenToUser(r.Context(), w, &dbUser)
}

func passTokenToUser(ctx context.Context, w http.ResponseWriter, user *model.User) {
	tokenString, err := auth.CreateToken(*user, false)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	refreshTokenString, err := auth.CreateToken(*user, true)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	token := model.Token{
		Username:     user.Username,
		Token:        tokenString,
		RefreshToken: refreshTokenString,
		UpdatedAt:    time.Now(),
	}

	err = db.UpsertToken(ctx, &token)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, false, http.StatusCreated, token)
}
