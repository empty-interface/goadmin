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
	cols, values, err := parseUrlParams(r, cols)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = currentSession.Conn.InsertRow(tableName, cols, values)
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
func parseUrlParams(r *http.Request, cols []string) ([]string, []interface{}, error) {

	params := r.URL.Query()

	//Note : what if a column is named "table" ?
	values := []interface{}{}
	names := []string{}
	for key, value := range params {
		if key == "table" {
			continue
		}
		if len(value) > 0 {
			if !validColName(cols, key) {
				return nil, nil, fmt.Errorf("Invalid column name : %s", key)
			}
			values = append(values, value[0])
			names = append(names, key)
		}
	}
	return names, values, nil
}

func validColName(cols []string, name string) bool {
	for _, col := range cols {
		if col == name {
			return true
		}
	}
	return false
}
