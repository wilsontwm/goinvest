package models

import (
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"time"
)

// Token : The JWT token for logging in the user
type Token struct {
	UserID uuid.UUID
	Name   string
	Email  string
	PicURL string
	Expiry time.Time
	jwt.StandardClaims
}

// User : Users that have registered an account on the platform
type User struct {
	Base
	Name          string
	Email         string
	Source        string
	PicURL        string
	FirebaseID    string
	LoginAt       time.Time
	FirebaseToken string `gorm:"-"`
	Token         string `gorm:"-"`
}

// UserLogin : Login the user, update the user profile accordingly (if necessary)
func UserLogin(user *User) error {
	db := GetDB()
	defer db.Close()

	// Verify the firebase token first
	opt := option.WithCredentialsFile("./firebase.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	decoded, err := client.VerifyIDToken(context.Background(), user.FirebaseToken)
	if err != nil {
		return fmt.Errorf("User is not valid, please try again later")
	}

	user.LoginAt = time.Now()
	user.FirebaseID = decoded.UID
	temp := &User{}
	db.Table("users").Where("email = ?", user.Email).First(temp)

	// 1. If user is not found, then create new user
	// 2. Login the user if email and firebase token is valid
	// 3. Return an error
	if temp.ID == uuid.Nil {
		db.Create(user)
	} else if temp.ID != uuid.Nil {
		db.Model(user).Where("email = ?", user.Email).Updates(map[string]interface{}{
			"name":     user.Name,
			"pic_url":  user.PicURL,
			"login_at": user.LoginAt,
		})
	}

	// Create a new JWT token for the newly login account
	expiry := time.Now().Add(time.Hour * 2) // Only valid for 2 hours
	tk := &Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		PicURL: user.PicURL,
		Expiry: expiry,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(jwtKey))
	user.Token = tokenString

	return nil
}
