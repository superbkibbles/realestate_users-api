package users

import (
	"fmt"

	"github.com/superbkibbles/bookstore_utils-go/logger"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_users-api/src/constants"
	"github.com/superbkibbles/realestate_users-api/src/datasources/mysqlclient"
	"github.com/superbkibbles/realestate_users-api/src/utils/mysql_utils"
)

const (
	queryLikeProperty = "INSERT INTO(property_id, user_id) values(?, ?)"
)

func (u *User) Get() (Users, rest_errors.RestErr) {
	stmt, err := mysqlclient.Session.Prepare(constants.QUERY_GET_ALL_USERS)
	if err != nil {
		logger.Error("error while trying to prepare get user statment", err)
		return nil, rest_errors.NewInternalServerErr("Database error", nil)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Error("error while trying to Get all users", err)
		return nil, rest_errors.NewInternalServerErr("Database error", nil)
	}
	defer rows.Close()
	results := make(Users, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Age, &user.Email, &user.UserName, &user.PhoneNumber, &user.Photo, &user.City, &user.GPS, &user.DateCreated, &user.Status, &user.Gender, &user.AppLanguage); err != nil {
			logger.Error("error while trying to scan user into user struct", err)
			return nil, rest_errors.NewInternalServerErr("Database error", nil)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundErr(fmt.Sprintf("No user Found"))
	}

	return results, nil
}

func (u *LikePrpertyReq) LikeProperty() rest_errors.RestErr {
	stmt, err := mysqlclient.Session.Prepare(constants.QUERY_INSERT_LIKE_PROPERTY)
	if err != nil {
		logger.Error("error while trying to prepare like property statment", err)
		return rest_errors.NewInternalServerErr("Database error", nil)
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.PropertyID, u.UserID)
	if err != nil {
		logger.Error("error while trying to save user", err)
		return rest_errors.NewInternalServerErr("Database error", nil)
	}

	return nil
}

func (u *UserForm) Save(fileName string) (*User, rest_errors.RestErr) {
	var photoLink string
	stmt, err := mysqlclient.Session.Prepare(constants.QUERY_INSERT_USER)
	if err != nil {
		logger.Error("error while trying to prepare get user statment", err)
		return nil, rest_errors.NewInternalServerErr("Database error", nil)
	}
	defer stmt.Close()

	if fileName != "" {
		photoLink = "http://localhost:8080/assets/" + fileName
	}

	insertRes, err := stmt.Exec(u.FirstName, u.LastName, u.Age, u.Email, u.UserName, u.Password, u.PhoneNumber, photoLink, u.City, u.GPS, u.DateCreated, u.Status, u.Gender, u.AppLanguage)
	if err != nil {
		logger.Error("error while trying to save user", err)
		return nil, rest_errors.NewInternalServerErr("Database error", nil)
	}

	userID, err := insertRes.LastInsertId()
	if err != nil {
		logger.Error("error while trying to get created user id", err)
		return nil, rest_errors.NewInternalServerErr("Database error", nil)
	}
	var user User
	user.Id = userID
	if err := user.GetByID(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetLikedProperties() rest_errors.RestErr {
	stmt, err := mysqlclient.Session.Prepare(constants.QUERY_GET_LIKED_PROPERTY)
	if err != nil {
		logger.Error("error while trying to prepare get liked propety statment", err)
		return rest_errors.NewInternalServerErr("Database Error", nil)
	}
	defer stmt.Close()

	rows, err := stmt.Query(u.Id)
	if err != nil {
		logger.Error("error while trying to get property by id", err)
		return mysql_utils.ParseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			logger.Error("error while trying to scan user into user struct", err)
			return rest_errors.NewInternalServerErr("Database error", nil)
		}
		u.LikedProperty = append(u.LikedProperty, s)
	}
	return nil
}

func (u *User) GetByID() rest_errors.RestErr {
	// Get user
	stmt, err := mysqlclient.Session.Prepare(constants.QUERY_GET_USER_ID)

	if err != nil {
		logger.Error("error while trying to prepare get by ID user statment", err)
		return rest_errors.NewInternalServerErr("Database Error", nil)
	}
	defer stmt.Close()
	res := stmt.QueryRow(u.Id)
	if getErr := res.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Email, &u.UserName, &u.PhoneNumber, &u.Photo, &u.City, &u.GPS, &u.DateCreated, &u.Status, &u.Gender, &u.AppLanguage); getErr != nil {
		logger.Error("error while trying to get user by ID", getErr)
		return mysql_utils.ParseError(getErr)
	}
	return u.GetLikedProperties()
}

func (u *User) Update() rest_errors.RestErr {
	stmt, err := mysqlclient.Session.Prepare(constants.QUERY_UPDATE_USER)
	if err != nil {
		logger.Error("error while trying to prepare Update user statment", err)
		return rest_errors.NewInternalServerErr("Database Error", nil)
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.Age, u.Email, u.UserName, u.PhoneNumber, u.City, u.GPS, u.Status, u.Gender, u.AppLanguage, u.Id)
	if err != nil {
		logger.Error("error while trying to prepare Update user statment", err)
		return rest_errors.NewInternalServerErr("Database Error", nil)
	}
	return nil
}

func (u *User) UpdatePhoto(fileName string) rest_errors.RestErr {
	stmt, err := mysqlclient.Session.Prepare(constants.QUERY_UPDATE_USER_PHOTO)
	if err != nil {
		logger.Error("error while trying to prepare Update user photo statment", err)
		return rest_errors.NewInternalServerErr("Database Error", nil)
	}
	defer stmt.Close()
	_, err = stmt.Exec("http://localhost:8080/assets/"+fileName, u.Id)
	if err != nil {
		logger.Error("error while trying to Exec Update user photo statment", err)
		return rest_errors.NewInternalServerErr("Database Error", nil)
	}

	u.Photo = "http://localhost:8080/assets/" + fileName

	return nil
}

func (u *User) Delete() rest_errors.RestErr {
	stmt, err := mysqlclient.Session.Prepare(constants.QUERY_DELETE_USER)
	if err != nil {
		logger.Error("error while trying to prepare Deactivate user", err)
		return rest_errors.NewInternalServerErr("Database Error", nil)
	}
	defer stmt.Close()

	_, err = stmt.Exec(StatusDeactivated, u.Id)
	if err != nil {
		logger.Error("error while trying to Deactivate user", err)
		return rest_errors.NewInternalServerErr("Database Error", nil)
	}

	return nil
}
