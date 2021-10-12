package users

import (
	"mime/multipart"
	"strings"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
)

const (
	StatusActive      = "active"
	StatusDeactivated = "deactivated"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Age         int64  `json:"age"`
	Email       string `json:"email"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Photo       string `json:"photo"`
	City        string `json:"city"`
	GPS         string `json:"gps"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Gender      string `json:"gender"`
	AppLanguage string `json:"app_language"`
}

type UserForm struct {
	Id          int64                 `form:"id" binding:"-"`
	FirstName   string                `form:"first_name" binding:"-"`
	LastName    string                `form:"last_name" binding:"-"`
	Age         int64                 `form:"age" binding:"-"`
	Email       string                `form:"email" binding:"-"`
	UserName    string                `form:"user_name" binding:"-"`
	Password    string                `form:"password" binding:"-"`
	PhoneNumber string                `form:"phone_number" binding:"-"`
	Photo       *multipart.FileHeader `form:"photo" binding:"-"`
	City        string                `form:"city" binding:"-"`
	GPS         string                `form:"gps" binding:"-"`
	DateCreated string                `form:"date_created" binding:"-"`
	Status      string                `form:"status" binding:"-"`
	Gender      string                `form:"gender" binding:"-"`
	AppLanguage string                `form:"app_language" binding:"-"`
}

type Users []User

func (user *User) Validate() rest_errors.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestErr("invalid Email address")
	}
	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return rest_errors.NewBadRequestErr("invalid Password")
	}
	return nil
}
