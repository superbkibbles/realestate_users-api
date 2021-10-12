package users

import (
	"encoding/json"
)

type PublicUser struct {
	Id          int64  `json:"id"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Age         int64  `json:"age"`
	Email       string `json:"email"`
	UserName    string `json:"user_name"`
	PhoneNumber string `json:"phone_number"`
	Photo       string `json:"photo"`
	City        string `json:"city"`
	GPS         string `json:"gps"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Gender      string `json:"gender"`
	AppLanguage string `json:"app_language"`
}

func (users Users) Marshal(isPublic bool) []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Marshal(isPublic)
	}
	return result
}

func (user *User) Marshal(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:          user.Id,
			DateCreated: user.DateCreated,
			Status:      user.Status,
		}
	}
	userJson, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(userJson, &privateUser)
	return privateUser
}

func (user *UserForm) Marshal(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:          user.Id,
			DateCreated: user.DateCreated,
			Status:      user.Status,
		}
	}
	userJson, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(userJson, &privateUser)
	return privateUser
}
