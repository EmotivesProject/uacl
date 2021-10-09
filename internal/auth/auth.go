package auth

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"uacl/messages"
	"uacl/model"

	jwt "github.com/golang-jwt/jwt/v4"
)

const (
	expirationAccessTime  = 15    // 15 minutes
	expirationRefreshTime = 43800 // 1 month
)

func CreateToken(user model.User, refreshToken bool) (string, error) {
	var expiresAt int64

	currentTime := time.Now()

	if refreshToken {
		expiresAt = currentTime.Add(time.Minute * expirationRefreshTime).Unix()
	} else {
		expiresAt = currentTime.Add(time.Minute * expirationAccessTime).Unix()
	}

	short := model.ShortenedUser{
		Name:      user.Name,
		Username:  user.Username,
		UserGroup: user.UserGroup,
	}

	claims := make(jwt.MapClaims)
	claims["dat"] = short
	claims["exp"] = expiresAt
	claims["iss"] = currentTime.Unix()
	claims["nbf"] = currentTime.Unix()

	private, err := ioutil.ReadFile(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		return "", err
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
}

func Validate(token string) (model.ShortenedUser, error) {
	var shorten model.ShortenedUser

	public, err := ioutil.ReadFile(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		return shorten, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return shorten, err
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, messages.ErrUnexpectedMethod
		}

		return key, nil
	})
	if err != nil {
		return shorten, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return shorten, messages.ErrInvalid
	}

	jsonString, err := json.Marshal(claims["dat"])
	if err != nil {
		return shorten, err
	}

	err = json.Unmarshal(jsonString, &shorten)

	return shorten, err
}
