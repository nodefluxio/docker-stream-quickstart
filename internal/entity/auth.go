package entity

import "github.com/dgrijalva/jwt-go"

// ClaimsStruct represent data injected to OAUTH2 token claims
type ClaimsStruct struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
	jwt.StandardClaims
}
