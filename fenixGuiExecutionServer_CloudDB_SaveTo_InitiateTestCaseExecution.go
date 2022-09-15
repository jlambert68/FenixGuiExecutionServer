package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	uuidGenerator "github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Prepare for Saving the Initiation of a new TestCaseExecution in the CloudDB
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) prepareInitiateTestCaseExecutionSaveToCloudDB(initiateSingleTestCaseExecutionRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionRequestMessage) (initiateSingleTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":    "306edce0-7a5a-4a0f-992b-5c9b69b0bcc6",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareInitiateTestCaseExecutionSaveToCloudDB'")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCaseExecutionUuid: "",
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem when saving to database",
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
			},
		}

		return initiateSingleTestCaseExecutionResponseMessage
	}
	defer txn.Commit(context.Background())

	// Generate a new TestCaseExecution-UUID
	testCaseExecutionUuid := uuidGenerator.New().String()

	// Generate TimeStamp
	placedOnTestExecutionQueueTimeStamp := time.Now()

	// Extract TestCase-information to be added to TestCaseExecution-data
	//testCaseToExecute := fenixTestCaseBuilderServerGrpcApi.BasicTestCaseInformationMessage{}
	testCaseToExecute, err := fenixGuiTestCaseBuilderServerObject.loadTestCaseBasicInformation(initiateSingleTestCaseExecutionRequestMessage.TestCaseUuid)
	if err != nil {

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCaseExecutionUuid: "",
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem when Loading TestCase Basic Information from database",
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
			},
		}

		return initiateSingleTestCaseExecutionResponseMessage
	}
	//TODO Load TestCase-data  from Database

	// Prepare TestDataExecution-data to be saved in database
	testCaseExecutionToBeSaved := fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage{
		DomainUuid:                          testCaseToExecute.NonEditableInformation.DomainUuid,
		DomainName:                          testCaseToExecute.NonEditableInformation.DomainName,
		TestSuiteUuid:                       "",
		TestSuiteName:                       "",
		TestSuiteVersion:                    0,
		TestSuiteExecutionUuid:              "",
		TestSuiteExecutionVersion:           0,
		TestCaseUuid:                        testCaseToExecute.NonEditableInformation.TestCaseUuid,
		TestCaseName:                        testCaseToExecute.EditableInformation.TestCaseName,
		TestCaseVersion:                     testCaseToExecute.NonEditableInformation.TestCaseVersion,
		TestCaseExecutionUuid:               testCaseExecutionUuid,
		TestCaseExecutionVersion:            1,
		PlacedOnTestExecutionQueueTimeStamp: timestamppb.New(placedOnTestExecutionQueueTimeStamp),
		TestDataSetUuid:                     initiateSingleTestCaseExecutionRequestMessage.TestDataSetUuid,
		ExecutionPriority:                   fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_HIGH_SINGLE_TESTCASE,
	}

	// Save the Initiation of a new TestCaseExecution in the CloudDB
	err = fenixGuiTestCaseBuilderServerObject.saveInitiateTestCaseExecutionSaveToCloudDB(txn, &testCaseExecutionToBeSaved)
	if err != nil {

		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":    "bc6f1da5-3c8c-493e-9882-0b20e0da9e2e",
			"error": err,
		}).Error("Couldn't Save TestCaseExecution in CloudDB")

		// Rollback any SQL transactions
		txn.Rollback(context.Background())

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCaseExecutionUuid: "",
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem when saving to database",
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
			},
		}

		return initiateSingleTestCaseExecutionResponseMessage

	}

	initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
		TestCaseExecutionUuid: testCaseExecutionUuid,
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
		},
	}

	return initiateSingleTestCaseExecutionResponseMessage
}

// Save the newly created TestExecution in database
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) saveInitiateTestCaseExecutionSaveToCloudDB(dbTransaction pgx.Tx, testCaseExecutionToBeSaved *fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "653795a1-d686-4823-9b5e-909dc37acc7d",
	}).Debug("Entering: saveInitiateTestCaseExecutionSaveToCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "b24a7105-5057-4ce3-a0bd-bdbbd5e138e0",
		}).Debug("Exiting: saveInitiateTestCaseExecutionSaveToCloudDB()")
	}()

	// Get a common dateTimeStamp to use
	//currentDataTimeStamp := fenixSyncShared.GenerateDatetimeTimeStampForDB()

	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""

	// Create Insert Statement for TestCaseExecution that will be put on ExecutionQueue
	// Data to be inserted in the DB-table
	dataRowsToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.DomainName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.DomainName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestSuiteUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestSuiteName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestSuiteVersion)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestSuiteExecutionUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestSuiteExecutionVersion)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestCaseUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestCaseName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestCaseVersion)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestCaseExecutionUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestCaseExecutionVersion)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.PlacedOnTestExecutionQueueTimeStamp)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.TestDataSetUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.ExecutionPriority)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute = sqlToExecute + "INSERT INTO \"" + usedDBSchema + "\".\"TestCaseExecutionQueue\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestSuiteUuid\", \"TestSuiteName\", \"TestSuiteVersion\", " +
		"\"TestSuiteExecutionUuid\", \"TestSuiteExecutionVersion\", \"TestCaseUuid\", \"TestCaseName\", \"TestCaseVersion\"," +
		" \"TestCaseExecutionUuid\", \"TestCaseExecutionVersion\", \"QueueTimeStamp\", \"TestDataSetUuid\", \"ExecutionPriority\") "
	sqlToExecute = sqlToExecute + fenixGuiTestCaseBuilderServerObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		return err
	}

	// Log response from CloudDB
	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id":                       "dcb110c2-822a-4dde-8bc6-9ebbe9fcbdb0",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
	}).Debug("Return data for SQL executed in database")

	// No errors occurred
	return nil

}

// Load BasicInformation for TestCase to be able to populate the TestCaseExecution
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestCaseBasicInformation(testCaseUuid string) (testCaseBasicInformation *fenixTestCaseBuilderServerGrpcApi.BasicTestCaseInformationMessage, err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TC.\"TestCaseBasicInformationAsJsonb\" "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCases\" TC "
	sqlToExecute = sqlToExecute + "WHERE TC.\"TestCaseUuid\" = '" + testCaseUuid + "' AND "
	sqlToExecute = sqlToExecute + "TC.\"TestCaseVersion\" = (SELECT MAX(TC2.\"TestCaseVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixGuiBuilder\".\"TestCases\" TC2 "
	sqlToExecute = sqlToExecute + "WHERE TC2.\"TestCaseUuid\" = '" + testCaseUuid + "');"

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "79c43b90-7539-4bab-bff9-41acfeb2b5bc",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return &fenixTestCaseBuilderServerGrpcApi.BasicTestCaseInformationMessage{}, err
	}

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		testCaseBasicInformation := fenixTestCaseBuilderServerGrpcApi.BasicTestCaseInformationMessage{}

		err := rows.Scan(
			&testCaseBasicInformation,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "ee291325-3517-4c3f-b571-13a3c37fc4b8",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

	}

	return testCaseBasicInformation, err

}

// See https://www.alexedwards.net/blog/using-postgresql-jsonb
// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a myAttrStruct) Value() (driver.Value, error) {

	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *myAttrStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type myAttrStruct struct {
	fenixTestCaseBuilderServerGrpcApi.BasicTestCaseInformationMessage
}
