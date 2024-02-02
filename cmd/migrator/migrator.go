package migrator

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"rapidEx/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func main() {
	var migrationsPath string
	flag.StringVar(&migrationsPath, "", "migrations", "path to migrations")

	if migrationsPath == "" {
		panic("migrations path is required")
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	c, err := config.LoadConfig(pwd)
	if err != nil {
		panic(err)
	}
	m, err := migrate.New("file://"+migrationsPath, fmt.Sprintf("mysql://%s:%s@tcp(host:port)/%s", c.MysqlUser, c.MysqlPassword, c.MysqlDBName))
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("Migrations applied successfully")
}
