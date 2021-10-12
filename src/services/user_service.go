package services

import (
	"mime/multipart"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_users-api/src/domain/users"
	"github.com/superbkibbles/realestate_users-api/src/utils/date_utils"
	"github.com/superbkibbles/realestate_users-api/src/utils/file_utils"
)

var (
	UserService userServiceInterface = &userService{}
)

type userServiceInterface interface {
	Get() (users.Users, rest_errors.RestErr)
	Create(users.UserForm, *multipart.FileHeader, multipart.File) (*users.User, rest_errors.RestErr)
	GetByID(int64) (*users.User, rest_errors.RestErr)
	UpdateUser(users.User) (*users.User, rest_errors.RestErr)
	UpdatePhoto(users.UserPhotoUpdate, *multipart.FileHeader, multipart.File) (*users.User, rest_errors.RestErr)
	DeleteUser(int64) rest_errors.RestErr
}

type userService struct{}

func (*userService) Get() (users.Users, rest_errors.RestErr) {
	dao := users.User{}
	return dao.Get()
}

func (*userService) GetByID(id int64) (*users.User, rest_errors.RestErr) {
	var user users.User
	user.Id = id
	if err := user.GetByID(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (*userService) Create(userForm users.UserForm, header *multipart.FileHeader, file multipart.File) (*users.User, rest_errors.RestErr) {
	userForm.DateCreated = date_utils.GetNowDBFromat()
	userForm.Status = users.StatusActive
	// Saving Pic
	fileName, err := file_utils.SaveFile(header, file)
	if err != nil {
		return nil, err
	}
	user, err := userForm.Save(fileName)
	if err != nil {
		file_utils.DeleteFile(fileName)
		return nil, err
	}
	return user, nil
}

func (us *userService) UpdateUser(user users.User) (*users.User, rest_errors.RestErr) {
	current, err := us.GetByID(user.Id)
	if err != nil {
		return nil, err
	}

	updateUser(current, user)
	// userJSON, _ := json.Marshal(user)
	// json.Unmarshal(userJSON, current)

	current.Update()

	return current, nil
}

func (*userService) UpdatePhoto(photoUpdateRequest users.UserPhotoUpdate, header *multipart.FileHeader, file multipart.File) (*users.User, rest_errors.RestErr) {
	var user users.User
	user.Id = photoUpdateRequest.UserID

	if err := user.GetByID(); err != nil {
		return nil, err
	}

	fileName, err := file_utils.UpdateFile(header, file, user.Photo)
	if err != nil {
		return nil, err
	}

	if err := user.UpdatePhoto(fileName); err != nil {
		file_utils.DeleteFile(fileName)
		return nil, err
	}

	return &user, nil
}

func (*userService) DeleteUser(userID int64) rest_errors.RestErr {
	user := &users.User{Id: userID}
	if err := user.Delete(); err != nil {
		return err
	}
	return nil
}

func updateUser(current *users.User, user users.User) {
	if user.FirstName != "" {
		current.FirstName = user.FirstName
	}
	if user.LastName != "" {
		current.LastName = user.LastName
	}
	if user.Age != 0 {
		current.Age = user.Age
	}
	if user.Email != "" {
		current.Email = user.Email
	}
	if user.UserName != "" {
		current.UserName = user.UserName
	}
	if user.PhoneNumber != "" {
		current.PhoneNumber = user.PhoneNumber
	}
	if user.City != "" {
		current.City = user.City
	}

	if user.GPS != "" {
		current.GPS = user.GPS
	}

	if user.Status != "" {
		current.Status = user.Status
	}

	if user.Gender != "" {
		current.Gender = user.Gender
	}
	if user.AppLanguage != "" {
		current.AppLanguage = user.AppLanguage
	}
}
