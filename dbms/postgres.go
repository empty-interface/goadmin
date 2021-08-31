package dbms

import (
	"database/sql"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(Config) Driver {
	return &Postgres{}
}

func (p *Postgres) Connect(config Config) error {
	return nil
}
func (p *Postgres) Query(text string, args ...interface{}) (*sql.Rows, error) {
	return p.DB.Query(text, args...)
}
func (p *Postgres) Close() {

}
