package dbms

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlDriver struct {
	cfg config
}

func newMySQLDriver(cfg config) driver {
	return &mysqlDriver{cfg}
}

func (p *mysqlDriver) open(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}
func (p *mysqlDriver) dsn() string {
	// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		p.cfg.Username,
		p.cfg.Password,
		p.cfg.Host,
		p.cfg.Port,
		p.cfg.DBName,
	)
}
