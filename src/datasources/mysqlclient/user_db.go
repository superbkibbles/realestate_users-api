package mysqlclient

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Client mysqlInterface = &mysqlClient{}
)

type mysqlInterface interface {
	setClient(*sql.DB)
}

type mysqlClient struct {
	client *sql.DB
}

var (
	Session *sql.DB
)

func Init() {
	dataSourceName := "root:sky-sharing12@/realestate_users"

	var err error
	Session, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err := Session.Ping(); err != nil {
		panic(err)
	}

	Client.setClient(Session)
}

func (c *mysqlClient) setClient(client *sql.DB) {
	c.client = client
}
