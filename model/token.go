package model

import jwt "github.com/dgrijalva/jwt-go"

//Token struct declaration
type Token struct {
	EncodedID string
	Name      string
	Email     string
	*jwt.StandardClaims
}
