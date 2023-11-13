package mysqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"rapidEx/config"
)

func SetMysqlConnection() (*sql.DB, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	c, err := config.LoadConfig(pwd)
	if err != nil {
		return nil, err
	}
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", c.MysqlUser,c.MysqlHost ,c.MysqlPassword, c.MysqlDBName))
}