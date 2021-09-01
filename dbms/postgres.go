package dbms

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDriver struct {
	cfg config
}

func newPostgresDriver(cfg config) driver {
	return postgresDriver{cfg}
}

func (p postgresDriver) open(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
func (p postgresDriver) dsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.cfg.Host, p.cfg.Port, p.cfg.Username, p.cfg.Password, p.cfg.DBName)
}
