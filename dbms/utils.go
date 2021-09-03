package dbms

import (
	"database/sql"
	"fmt"
	"strings"
)

type Table struct {
	Catalog                      string `gorm:"column:table_catalog"`
	Schema                       string `gorm:"column:table_schema"`
	Name                         string `gorm:"column:table_name"`
	Type                         string `gorm:"column:table_type"`
	Self_referencing_column_name string `gorm:"column:self_referencing_column_name"`
}

func (*Table) TableName() string {

	return "information_schema.tables"
}
func wrapString(s *string) {
	*s = fmt.Sprintf(`"%s"`, *s)
}
func wrapStrings(s *[]string) {
	for i := range *s {
		wrapString(&(*s)[i])
	}
}
func (conn *GormConnection) InsertRow(table string, columns []string, values []interface{}) error {
	wrapStrings(&columns)
	wrapString(&table)
	cols := strings.Join(columns, ",")
	_vals := make([]string, len(values))
	for i := range values {
		_vals[i] = "?"
	}
	vals := strings.Join(_vals, ",")
	query := fmt.Sprintf(`insert into %s (%s) values (%s)`, table, cols, vals)
	conn.Exec(query, values...)
	return conn.Error
}
func (conn *GormConnection) GetTableColumns(name string) ([]string, error) {
	query := fmt.Sprintf(`select * from %s Limit 1`, name)
	rows, err := conn.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	return cols, nil
}
func (conn *GormConnection) GetTables() []Table {
	// query := `select * from information_schema.tables where table_schema = 'public'`
	// tables := make([]Table, 0)
	// rows, err := conn.Raw(query).Rows()
	// if err != nil {
	// 	return
	// }
	// rows.Scan()
	tables := []Table{}
	conn.Where("table_schema=?", "public").Find(&tables)
	fmt.Println("Found", len(tables), "tables")
	return tables
}

type TableInfo struct {
	Types []*sql.ColumnType
}

func (conn *GormConnection) GetTableRows(query string, limit int) ([][]interface{}, []string, error) {
	// query := fmt.Sprintf(`select * from %s limit %v`, name, limit)
	rows, err := conn.Raw(query).Rows()
	if err != nil {
		return nil, nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	var result [][]interface{}
	temp := make([]interface{}, len(cols))
	for rows.Next() {
		row := make([]interface{}, len(cols))
		for i := range temp {
			temp[i] = &row[i]
		}
		rows.Scan(temp...)
		result = append(result, row)
	}
	return result, cols, nil
}
func (conn *GormConnection) GetTableInfo(name string) (*TableInfo, error) {
	query := fmt.Sprintf(`select * from %s`, name)
	rows, err := conn.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	return &TableInfo{
		Types: types,
	}, nil
}
