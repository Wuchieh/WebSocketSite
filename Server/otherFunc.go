package Server

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
	//a := string(b)
	//return &a
}

func genNewToken(session string) *UserToken {
	if token, ok := userTokens[session]; ok {
		token.UpdateTime = time.Now()
		token.ExpiredTime = time.Now().Add(time.Hour)
		token.Token = RandomString(32)
	} else {
		userToken := UserToken{
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
			ExpiredTime: time.Now().Add(time.Hour),
			Token:       RandomString(32),
		}
		userTokens[session] = &userToken
	}
	return userTokens[session]
}

func updateToken(session string) bool {
	if token, ok := userTokens[session]; ok {
		if token.UpdateTime.Add(1*time.Minute).Unix() > time.Now().Unix() {
			return false
		}
		token.UpdateTime = time.Now()
		token.ExpiredTime = time.Now().Add(time.Hour)
		token.Token = RandomString(32)
	} else {
		if genNewToken(session) == nil {
			return false
		}
	}
	return true
}
