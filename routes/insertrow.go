package routes

import (
	"fmt"
	"net/http"

	"github.com/empty-interface/goadmin/session"
)

const InsertRowPath = "/insertrow"

func HandleInsert(w http.ResponseWriter, r *http.Request, currentSession *session.Session) (int, error) {
	tableName, err := parseTableName(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	cols, err := currentSession.Conn.GetTableColumns(tableName)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	values, includebools, err := parseUrlParams(r, cols)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	fmt.Println("Columns to insert", values)
	if includebools {
		includeBools(currentSession, tableName, &values)
		fmt.Println("Columns to insert", values)
	}
	err = currentSession.Conn.InsertRow(tableName, values)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	http.Redirect(w, r, fmt.Sprintf("table?name=%s&section=select", tableName), http.StatusPermanentRedirect)
	return -1, nil
}
func parseTableName(r *http.Request) (string, error) {
	tableNames, exists := r.URL.Query()["table"]
	if !exists {
		return "", fmt.Errorf("Invalid table name")
	}
	tableName := tableNames[0]
	return tableName, nil
}
func parseUrlParams(r *http.Request, cols []string) (map[string]string, bool, error) {

	params := r.URL.Query()

	//Note : what if a column is named "table" ?
	values := map[string]string{}
	includebools := false
	for key, value := range params {
		if key == "table" {
			continue
		}
		if key == "includebools" {
			includebools = true
			continue
		}
		fmt.Println("key")
		if len(value) > 0 {
			if !validColName(cols, key) {
				return nil, false, fmt.Errorf("Invalid column name : %s", key)
			}
			values[key] = value[0]
		}
	}
	return values, includebools, nil
}

func validColName(cols []string, name string) bool {
	for _, col := range cols {
		if col == name {
			return true
		}
	}
	return false
}
func includeBools(currentSession *session.Session, name string, valuesToInsert *map[string]string) error {
	tableInfo, err := showStructure(currentSession, name)
	if err != nil {
		return err
	}
	inputs := parseColumnTypes(tableInfo)
	for k, v := range inputs {
		if v == "checkbox" {
			if _, exists := (*valuesToInsert)[k]; !exists {
				(*valuesToInsert)[k] = "false"
			}

		}
	}
	return nil
}
