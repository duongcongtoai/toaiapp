package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string
	Password string
}

// HasPermission Todo
func (u *User) HasPermission(_ string) bool {
	return true
}

func (u *User) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["ID"] = u.ID
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Duration(345600 * time.Second)).Unix()
	tokenString, err := token.SignedString(SignKey)
	if err != nil {
		return "", fmt.Errorf("Failed to generate JWT token, error was: %v", err)
	}
	return tokenString, nil
}
