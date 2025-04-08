package auth

import (
	"fmt"

	"github.com/BingyanStudio/configIT/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func encrypt(pwd string) (string, error) {
	hashStr, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashStr), err
}

func validate(hash, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))

	return err == nil
}

func Login(username, password string) (*model.User, error) {
	user, err := model.GetUserBySub(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	if user.Password == "" {
		return nil, fmt.Errorf("user authenticated by OIDC, not password")
	}
	if validate(user.Password, password) {
		return user, nil
	}
	return nil, fmt.Errorf("password not match")
}
