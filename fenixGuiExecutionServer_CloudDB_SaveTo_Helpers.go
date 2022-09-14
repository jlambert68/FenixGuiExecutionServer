package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// Generates all "VALUES('xxx', 'yyy')..." for insert statements
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) generateSQLInsertValues(testdata [][]interface{}) (sqlInsertValuesString string) {

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
			case string:

				sqlInsertValuesString = sqlInsertValuesString + "'" + fmt.Sprint(value) + "'"

			default:
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"id": "33e11bc9-bfc7-4c2f-8440-30f8d9a89ab0",
				}).Fatal("Unhandled type, %valueType", valueType)
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
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) generateSQLINArray(testdata []string) (sqlInsertValuesString string) {

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
