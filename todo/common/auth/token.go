package auth

import (
	"encoding/json"
	"fmt"
)

// Token is an authorization token like JWT.
type Token struct {
	user *User
}

// NewToken create token for user
func NewToken(u *User) *Token {
	return &Token{user: u}
}

// NewTokenFromString constructor of Token from string
func NewTokenFromString(tokenString string) (*Token, error) {
	u := new(User)
	err := json.Unmarshal([]byte(tokenString), u)
	if err != nil {
		return nil, err
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("invalid token. UserID == 0 in token: %s", tokenString)
	}
	t := NewToken(u)
	return t, nil
}

// User return authorized User object
func (t *Token) User() *User {
	return t.user
}

// String should return encoded token as string
func (t *Token) String() (string, error) {
	b, err := json.Marshal(t.user)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
