package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToken(t *testing.T) {
	u := &User{ID: 123}
	token := NewToken(u)
	tokenStr, err := token.String()
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	token2, err := NewTokenFromString(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, token2)
	assert.Equal(t, u, token2.User())

	u3 := &User{ID: 1234}
	assert.NotEqual(t, u3, token2.User())
}

func TestTokenNegative(t *testing.T) {
	u := &User{ID: 123}
	token := NewToken(u)
	tokenStr, err := token.String()
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	token2, err := NewTokenFromString(tokenStr + "!")
	assert.Error(t, err)
	assert.Nil(t, token2)

	token2, err = NewTokenFromString("")
	assert.Error(t, err)
	assert.Nil(t, token2)
}
