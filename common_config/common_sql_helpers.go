package common_config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

// Generates all "VALUES('xxx', 'yyy')..." for insert statements
func GenerateSQLInsertValues(testdata [][]interface{}) (sqlInsertValuesString string) {

	sqlInsertValuesString = ""

	// Loop over both rows and values
	for rowCounter, rowValues := range testdata {
		if rowCounter == 0 {
			// Only add 'VALUES' for first row
			sqlInsertValuesString = sqlInsertValuesString + "VALUES("
		} else {
			sqlInsertValuesString = sqlInsertValuesString + ",("
		}

		for valueCounter, value := range rowValues {
			switch valueType := value.(type) {

			case bool:
				sqlInsertValuesString = sqlInsertValuesString + fmt.Sprint(value)

			case int:
				sqlInsertValuesString = sqlInsertValuesString + fmt.Sprint(value)

			case int64:
				sqlInsertValuesString = sqlInsertValuesString + fmt.Sprint(value)

			case uint32:
				sqlInsertValuesString = sqlInsertValuesString + fmt.Sprint(value)

			case string:
				valuePrepared := strings.ReplaceAll(fmt.Sprint(value), "'", "''")
				sqlInsertValuesString = sqlInsertValuesString + "'" + valuePrepared + "'"

			default:
				Logger.WithFields(logrus.Fields{
					"id":    "924f5986-68b7-48cc-b093-17bdf7608476",
					"value": value,
				}).Fatal("Unhandled type, %valueType:", valueType)
			}

			// After the last value then add ')'
			if valueCounter == len(rowValues)-1 {
				sqlInsertValuesString = sqlInsertValuesString + ") "
			} else {
				// Not last value, so Add ','
				sqlInsertValuesString = sqlInsertValuesString + ", "
			}

		}

	}

	return sqlInsertValuesString
}

// Generates incoming values in the following form:  "('monkey', 'tiger'. 'fish')"
func GenerateSQLINArray(testdata []string) (sqlInsertValuesString string) {

	// Create a list with '' as only element if there are no elements in array
	if len(testdata) == 0 {
		sqlInsertValuesString = "('')"

		return sqlInsertValuesString
	}

	sqlInsertValuesString = "("

	// Loop over both rows and values
	for counter, value := range testdata {

		if counter == 0 {
			// Only used for first row
			sqlInsertValuesString = sqlInsertValuesString + "'" + value + "'"

		} else {

			sqlInsertValuesString = sqlInsertValuesString + ", '" + value + "'"
		}
	}

	sqlInsertValuesString = sqlInsertValuesString + ") "

	return sqlInsertValuesString
}
