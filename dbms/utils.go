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
func (conn *GormConnection) InsertRow(table string, row map[string]string) error {
	_vals, _cols, values := []string{}, []string{}, []interface{}{}
	for k, v := range row {
		wrapString(&k)
		_cols = append(_cols, k)
		_vals = append(_vals, "?")
		values = append(values, v)
	}
	cols := strings.Join(_cols, ",")
	vals := strings.Join(_vals, ",")
	query := fmt.Sprintf(`insert into %s (%s) values (%s)`, table, cols, vals)
	err := conn.Exec(query, values...).Error
	return err
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

type Schema struct {
	CatalogName                string `gorm:"column:catalog_name"`
	SchemaName                 string `gorm:"column:schema_name"`
	SchemaOwner                string `gorm:"column:schema_owner"`
	DefaultCharacterSetCatalog string `gorm:"column:default_character_set_catalog"`
	DefaultCharacterSetSchema  string `gorm:"column:default_character_set_schema"`
	DefaultCharacterSetName    string `gorm:"column:default_character_set_name"`
	SqlPath                    string `gorm:"column:sql_path"`
}

func (*Schema) TableName() string {
	return "information_schema.schemata"
}
func (conn *GormConnection) GetDBSchemas() []Schema {
	var schemas []Schema
	conn.Find(&schemas)
	return schemas
}
func (conn *GormConnection) GetSchemaTables(schema string) []Table {
	// query := `select * from information_schema.tables where table_schema = 'public'`
	// tables := make([]Table, 0)
	// rows, err := conn.Raw(query).Rows()
	// if err != nil {
	// 	return
	// }
	// rows.Scan()
	tables := []Table{}
	conn.Where("table_schema=?", schema).Find(&tables)
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
	result := scanRows(rows, len(cols))
	return result, cols, nil
}
func scanRows(rows *sql.Rows, len_ int) [][]interface{} {
	var result [][]interface{}
	temp := make([]interface{}, len_)
	for rows.Next() {
		row := make([]interface{}, len_)
		for i := range temp {
			temp[i] = &row[i]
		}
		rows.Scan(temp...)
		result = append(result, row)
	}
	return result
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
