package users

import (
	"fmt"

	"github.com/superbkibbles/bookstore_utils-go/logger"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_users-api/src/datasources/mysqlclient"
	"github.com/superbkibbles/realestate_users-api/src/utils/mysql_utils"
)

const (
	queryInsertUsers     = "INSERT INTO users(first_name, last_name, age, email, user_name, password, phone_number, photo, city, gps, date_created, status, gender, app_language) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	queryGetUserByID     = "SELECT id, first_name, last_name, age, email, user_name, phone_number, photo, city, gps, date_created, status, gender, app_language FROM users WHERE id=?;"
	queryGetAllUsers     = "SELECT id, first_name, last_name, age, email, user_name, phone_number, photo, city, gps, date_created, status, gender, app_language FROM users;"
	queryUpdateUser      = "UPDATE users SET first_name=?, last_name=?, age=?, email=?, user_name=?, phone_number=?, city=?, gps=?, status=?, gender=?, app_language=? WHERE id=?;"
	queryUpdateUserPhoto = "UPDATE users SET photo=? WHERE id=?;"
	queryDeleteUser      = "UPDATE users set status=? where id=?"
)

// NOT IMPLEMENTED
func (u *User) Get() (Users, rest_errors.RestErr) {
	stmt, err := mysqlclient.Session.Prepare(queryGetAllUsers)
	if err != nil {
		logger.Error("error while trying to prepare get user statment", err)
		return nil, rest_errors.NewInternalServerErr("Database error", nil)
	}
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

func (u *UserForm) Save(fileName string) (*User, rest_errors.RestErr) {
	stmt, err := mysqlclient.Session.Prepare(queryInsertUsers)
	if err != nil {
		logger.Error("error while trying to prepare get user statment", err)
		return nil, rest_errors.NewInternalServerErr("Database error", nil)
	}
	defer stmt.Close()

	photoLink := "http://localhost:8080/assets/" + fileName

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
	// Create instance of User and assign the id to it
	var user User
	user.Id = userID
	if err := user.GetByID(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetByID() rest_errors.RestErr {
	stmt, err := mysqlclient.Session.Prepare(queryGetUserByID)
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
	return nil
}

func (u *User) Update() rest_errors.RestErr {
	stmt, err := mysqlclient.Session.Prepare(queryUpdateUser)
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
	stmt, err := mysqlclient.Session.Prepare(queryUpdateUserPhoto)
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
	stmt, err := mysqlclient.Session.Prepare(queryDeleteUser)
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
