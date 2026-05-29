package utils

import (
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func MapMySQLError(err *mysql.MySQLError, path string) (int, string, string) {
	switch err.Number {
	case 1062:
		return http.StatusConflict, "Duplicate entry is not allowed.", devMsg(err, path)
	case 1451:
		return http.StatusConflict, "Cannot delete this item as it is referenced elsewhere.", devMsg(err, path)
	case 1452:
		return http.StatusBadRequest, "Reference to a non-existent related record.", devMsg(err, path)
	case 1048:
		return http.StatusBadRequest, "A required field was not provided.", devMsg(err, path)
	case 1364:
		return http.StatusBadRequest, "A required field has no default value and was not provided.", devMsg(err, path)
	case 1406:
		return http.StatusBadRequest, "Input value is too long for one of the fields.", devMsg(err, path)
	case 1054:
		return http.StatusBadRequest, "Invalid field specified in the request.", devMsg(err, path)
	case 1146:
		return http.StatusInternalServerError, "Internal error: Required data table is missing.", devMsg(err, path)
	case 1366:
		return http.StatusBadRequest, "Invalid character encoding in input.", devMsg(err, path)
	case 2002, 2006, 2013:
		return http.StatusServiceUnavailable, "Database server is unavailable. Please try again later.", devMsg(err, path)
	default:
		return http.StatusInternalServerError, "A database error occurred. Please contact support.", devMsg(err, path)
	}
}

func devMsg(err *mysql.MySQLError, path string) string {
	return fmt.Sprintf("MySQL Error at %s: Code %d, Message: %s", path, err.Number, err.Message)
}
