package model

import jwt "github.com/dgrijalva/jwt-go"

//Token struct declaration
type Token struct {
	ID    int
	Name  string
	Email string
	*jwt.StandardClaims
}
