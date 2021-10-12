package mysql_utils

import (
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
)

const (
	ErrorNoRow = "no rows in result set"
)

func ParseError(err error) rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRow) {
			return rest_errors.NewNotFoundErr("No record matching the given id")
		}
		return rest_errors.NewInternalServerErr("error parsing database response", nil)
	}
	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestErr("invalid data")
	}
	return rest_errors.NewInternalServerErr("Database Error", nil)
}
