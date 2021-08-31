package dbms

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDriver struct {
	config Config
}

func newPostgresDriver(config Config) driver {
	return postgresDriver{config}
}

func (p postgresDriver) open(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
func (p postgresDriver) dsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.config.Host, p.config.Port, p.config.Username, p.config.Password, p.config.DBName)
}
