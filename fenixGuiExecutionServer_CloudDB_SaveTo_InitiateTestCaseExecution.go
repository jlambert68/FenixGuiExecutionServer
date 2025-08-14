package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	uuidGenerator "github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// After all stuff is done, then Commit or Rollback depending on result
var doCommitNotRoleBack bool

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) commitOrRoleBack(dbTransaction pgx.Tx) {
	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())
	} else {
		dbTransaction.Rollback(context.Background())
	}
}

// Holds the specific TestSuiteInformation when 'InitiateTestCaseExecution' is started from a TestSuite
type testSuiteInformationStruct struct {
	suiteUuid             string
	suiteName             string
	suiteVersion          uint32
	suiteExecutionUuid    string
	suiteExecutionVersion uint32
}

// Prepare for Saving the Initiation of a new TestCaseExecution in the CloudDB
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) prepareInitiateTestCaseExecutionSaveToCloudDB(
	txnToUse pgx.Tx,
	initiateSingleTestCaseExecutionRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionRequestMessage,
	executionPriority fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum,
	testSuiteInformation testSuiteInformationStruct) (
	initiateSingleTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "2a2009f9-af78-4216-b1d8-b1a0519e7041",
		"initiateSingleTestCaseExecutionRequestMessage": initiateSingleTestCaseExecutionRequestMessage,
	}).Debug("Incoming 'prepareInitiateTestCaseExecutionSaveToCloudDB'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "1efe7107-888d-487b-8445-9f81fe1a2c62",
	}).Debug("Outgoing 'prepareInitiateTestCaseExecutionSaveToCloudDB'")

	var err error
	var txn pgx.Tx

	// Begin SQL Transaction
	if txnToUse != nil {
		// WHen called from Create TestSuiteExecution then use incoming 'sql-txn'
		txn = txnToUse
	} else {
		txn, err = fenixSyncShared.DbPool.Begin(context.Background())
		if err != nil {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
				TestCasesInExecutionQueue: nil,
				AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
					AckNack:                      false,
					Comments:                     "Problem when saving to database",
					ErrorCodes:                   errorCodes,
					ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				},
			}

			return initiateSingleTestCaseExecutionResponseMessage
		}
	}

	// Standard is to do a Rollback
	doCommitNotRoleBack = false
	// Begin SQL Transaction
	if txnToUse == nil {
		// WHen called from Create TestSuiteExecution then use incoming 'sql-txn':s Commit/Rollback
		defer fenixGuiExecutionServerObject.commitOrRoleBack(txn) //txn.Commit(context.Background())
	}

	// Generate a new TestCaseExecution-UUID
	testCaseExecutionUuid := uuidGenerator.New().String()

	// Generate TimeStamp
	placedOnTestExecutionQueueTimeStamp := time.Now().UTC()

	// Extract TestCase-information to be added to TestCaseExecution-data
	//testCaseToExecuteBasicInformation := fenixGuiExecutionServerObject.BasicTestCaseInformationMessage{}
	testCaseToExecuteBasicInformation, err := fenixGuiExecutionServerObject.loadTestCaseBasicInformation(
		txn,
		initiateSingleTestCaseExecutionRequestMessage.TestCaseUuid)

	if err != nil {

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem when Loading TestCase Basic Information from database",
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
		}

		return initiateSingleTestCaseExecutionResponseMessage
	}
	//TODO Load TestCase-data  from Database

	var testCaseExecutionToBeSaved fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage

	// Prepare TestDataExecution-data to be saved in database based in priority which depends on what initiated it
	switch executionPriority {
	case fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_ExecutionPriorityEnum_DEFAULT_NOT_SET:

		errMsg := "ExecutionPriority is not set: " + executionPriority.String()

		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "642415a7-e53f-43ff-8c24-5566340e4d4c",
			"error": err,
		}).Error(errMsg)

		// Create Return message
		initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     errMsg,
				ErrorCodes:                   nil,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
		}

		return initiateSingleTestCaseExecutionResponseMessage

	case fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_HIGHEST_PROBES,
		fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_HIGH_SINGLE_TESTSUITE,
		fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_MEDIUM_MULTIPLE_TESTSUITES,
		fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_LOW_SCHEDULED_TESTSUITES:

		// Secure that TestSuiteInformation exists
		if len(testSuiteInformation.suiteUuid) != 36 ||
			len(testSuiteInformation.suiteName) == 0 ||
			testSuiteInformation.suiteVersion == 0 ||
			len(testSuiteInformation.suiteExecutionUuid) != 36 ||
			testSuiteInformation.suiteExecutionVersion == 0 {

			errMsg := "TestSuiteInformation is not correct set"

			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"id":                   "b101b5c7-6636-403e-b136-ccd426377e3e",
				"testSuiteInformation": testSuiteInformation,
			}).Error(errMsg)

			// Create Return message
			initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
				TestCasesInExecutionQueue: nil,
				AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
					AckNack:                      false,
					Comments:                     errMsg,
					ErrorCodes:                   nil,
					ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				},
			}

			return initiateSingleTestCaseExecutionResponseMessage

		}

		testCaseExecutionToBeSaved = fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage{
			DomainUuid:                          testCaseToExecuteBasicInformation.domainUuid,
			DomainName:                          testCaseToExecuteBasicInformation.domainName,
			TestSuiteUuid:                       testSuiteInformation.suiteUuid,
			TestSuiteName:                       testSuiteInformation.suiteName,
			TestSuiteVersion:                    testSuiteInformation.suiteVersion,
			TestSuiteExecutionUuid:              testSuiteInformation.suiteExecutionUuid,
			TestSuiteExecutionVersion:           testSuiteInformation.suiteExecutionVersion,
			TestCaseUuid:                        testCaseToExecuteBasicInformation.testCaseUuid,
			TestCaseName:                        testCaseToExecuteBasicInformation.testCaseName,
			TestCaseVersion:                     uint32(testCaseToExecuteBasicInformation.testCaseVersion),
			TestCaseExecutionUuid:               testCaseExecutionUuid,
			TestCaseExecutionVersion:            1,
			PlacedOnTestExecutionQueueTimeStamp: timestamppb.New(placedOnTestExecutionQueueTimeStamp),
			TestDataSetUuid:                     initiateSingleTestCaseExecutionRequestMessage.TestDataSetUuid,
			ExecutionPriority:                   executionPriority,
		}

	case fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_HIGH_SINGLE_TESTCASE,
		fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_MEDIUM_MULTIPLE_TESTCASES:

		testCaseExecutionToBeSaved = fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage{
			DomainUuid:                          testCaseToExecuteBasicInformation.domainUuid,
			DomainName:                          testCaseToExecuteBasicInformation.domainName,
			TestSuiteUuid:                       common_config.ZeroUuid,
			TestSuiteName:                       "",
			TestSuiteVersion:                    0,
			TestSuiteExecutionUuid:              common_config.ZeroUuid,
			TestSuiteExecutionVersion:           0,
			TestCaseUuid:                        testCaseToExecuteBasicInformation.testCaseUuid,
			TestCaseName:                        testCaseToExecuteBasicInformation.testCaseName,
			TestCaseVersion:                     uint32(testCaseToExecuteBasicInformation.testCaseVersion),
			TestCaseExecutionUuid:               testCaseExecutionUuid,
			TestCaseExecutionVersion:            1,
			PlacedOnTestExecutionQueueTimeStamp: timestamppb.New(placedOnTestExecutionQueueTimeStamp),
			TestDataSetUuid:                     initiateSingleTestCaseExecutionRequestMessage.TestDataSetUuid,
			ExecutionPriority:                   executionPriority,
		}

	default:

		errMsg := "Unknown ExecutionPriority: " + executionPriority.String()

		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "cd67ded8-3f7f-410f-8964-b4dad0197d31",
			"error": err,
		}).Error(errMsg)

		// Create Return message
		initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     errMsg,
				ErrorCodes:                   nil,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
		}

		return initiateSingleTestCaseExecutionResponseMessage

	}

	// Save the Initiation of a new TestCaseExecution in the CloudDB
	err = fenixGuiExecutionServerObject.saveInitiateTestCaseExecutionSaveToCloudDB(
		txn,
		&testCaseExecutionToBeSaved,
		initiateSingleTestCaseExecutionRequestMessage.GetExecutionStatusReportLevel())
	if err != nil {

		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
			TestCasesInExecutionQueue: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem when saving to database",
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
		}

		return initiateSingleTestCaseExecutionResponseMessage

	}

	// Save the TestData to be used for of a new TestCaseExecution in the CloudDB
	err = fenixGuiExecutionServerObject.saveTestDataForTestCaseExecutionToCloudDB(
		txn,
		&testCaseExecutionToBeSaved,
		initiateSingleTestCaseExecutionRequestMessage.TestDataForTestCaseExecution)
	if err != nil {

		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "98502463-3670-45bd-be5a-29b8157fdd78",
			"error": err,
		}).Error("Couldn't Save TestData in CloudDB")

		// Rollback any SQL transactions
		txn.Rollback(context.Background())

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem when saving testdata to database",
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
		}

		return initiateSingleTestCaseExecutionResponseMessage

	}

	initiateSingleTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
		TestCasesInExecutionQueue: &testCaseExecutionToBeSaved,
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
	}

	// Commit every database change
	doCommitNotRoleBack = true

	return initiateSingleTestCaseExecutionResponseMessage
}

// Save the newly created TestExecution in database
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) saveInitiateTestCaseExecutionSaveToCloudDB(
	dbTransaction pgx.Tx,
	testCaseExecutionToBeSaved *fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage,
	executionStatusReportLevel fenixExecutionServerGuiGrpcApi.ExecutionStatusReportLevelEnum) (
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"Id": "653795a1-d686-4823-9b5e-909dc37acc7d",
	}).Debug("Entering: saveInitiateTestCaseExecutionSaveToCloudDB()")

	defer func() {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.DomainUuid)
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
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, int(executionStatusReportLevel))

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute = sqlToExecute + "INSERT INTO \"" + usedDBSchema + "\".\"TestCaseExecutionQueue\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestSuiteUuid\", \"TestSuiteName\", \"TestSuiteVersion\", " +
		"\"TestSuiteExecutionUuid\", \"TestSuiteExecutionVersion\", \"TestCaseUuid\", \"TestCaseName\", \"TestCaseVersion\"," +
		" \"TestCaseExecutionUuid\", \"TestCaseExecutionVersion\", \"QueueTimeStamp\", \"TestDataSetUuid\", \"ExecutionPriority\", "
	sqlToExecute = sqlToExecute + "\"ExecutionStatusReportLevel\") "
	sqlToExecute = sqlToExecute + fenixGuiExecutionServerObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "374c7116-2c1b-4ec8-9318-105d25c08aab",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'saveInitiateTestCaseExecutionSaveToCloudDB'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "5bfd73be-d0f6-482e-9f75-243028f83b39",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Log response from CloudDB
	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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

// Save the TestData, in the database, to be used when doing the TestExecution
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) saveTestDataForTestCaseExecutionToCloudDB(
	dbTransaction pgx.Tx,
	testCaseExecutionToBeSaved *fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage,
	testDataForTestCaseExecutionMessage *fenixExecutionServerGuiGrpcApi.TestDataForTestCaseExecutionMessage) (
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"Id": "d17a8d3e-72a3-439b-9d74-1f63341a3414",
	}).Debug("Entering: saveTestDataForTestCaseExecutionToCloudDB()")

	defer func() {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id": "b51778ea-83ba-4cd3-a84e-75032f317d63",
		}).Debug("Exiting: saveTestDataForTestCaseExecutionToCloudDB()")
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

	// Convert TestDataForTestCaseExecution into json to be later stored as jsonb
	var tempTestDataForTestCaseExecutionAsJsonb string
	tempTestDataForTestCaseExecutionAsJsonb = protojson.Format(testDataForTestCaseExecutionMessage)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseExecutionToBeSaved.DomainUuid)
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
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestDataForTestCaseExecutionAsJsonb)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute = sqlToExecute + "INSERT INTO \"" + usedDBSchema + "\".\"TestDataForTestCaseExecution\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestSuiteUuid\", \"TestSuiteName\", \"TestSuiteVersion\", " +
		"\"TestSuiteExecutionUuid\", \"TestSuiteExecutionVersion\", \"TestCaseUuid\", \"TestCaseName\", \"TestCaseVersion\"," +
		" \"TestCaseExecutionUuid\", \"TestCaseExecutionVersion\", \"InsertedTimeStamp\", "
	sqlToExecute = sqlToExecute + "\"TestDataForTestCaseExecutionAsJsonb\") "
	sqlToExecute = sqlToExecute + fenixGuiExecutionServerObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "58af29a9-53a0-419b-9cb6-fc9c6c258cca",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'saveTestDataForTestCaseExecutionToCloudDB'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "5bfd73be-d0f6-482e-9f75-243028f83b39",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Log response from CloudDB
	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"Id":                       "9afc0f63-c373-4434-850c-7a426414be1a",
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

// Temporary variabel for storing temp result from database
type tempTestCaseBasicInformationStruct struct {
	domainUuid      string
	domainName      string
	testCaseUuid    string
	testCaseName    string
	testCaseVersion int
}

// Load BasicInformation for TestCase to be able to populate the TestCaseExecution
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadTestCaseBasicInformation(dbTransaction pgx.Tx, testCaseUuid string) (testCaseBasicInformation tempTestCaseBasicInformationStruct, err error) {

	usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TC.\"DomainUuid\", TC.\"DomainName\", TC.\"TestCaseUuid\", TC.\"TestCaseName\", TC.\"TestCaseVersion\""
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCases\" TC "
	sqlToExecute = sqlToExecute + "WHERE TC.\"TestCaseUuid\" = '" + testCaseUuid + "' AND "
	sqlToExecute = sqlToExecute + "TC.\"TestCaseVersion\" = (SELECT MAX(TC2.\"TestCaseVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" TC2 "
	sqlToExecute = sqlToExecute + "WHERE TC2.\"TestCaseUuid\" = '" + testCaseUuid + "');"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "e1955234-00ce-4cfd-a1de-cfae9bf46792",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadTestCaseBasicInformation'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "79c43b90-7539-4bab-bff9-41acfeb2b5bc",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return tempTestCaseBasicInformationStruct{}, err
	}

	// USed to secure that exactly one row was found
	numberOfRowFromDB := 0

	// Extract data from DB result set
	for rows.Next() {

		numberOfRowFromDB = numberOfRowFromDB + 1

		err := rows.Scan(
			&testCaseBasicInformation.domainUuid,
			&testCaseBasicInformation.domainName,
			&testCaseBasicInformation.testCaseUuid,
			&testCaseBasicInformation.testCaseName,
			&testCaseBasicInformation.testCaseVersion,
		)

		if err != nil {

			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":           "9cdde993-689a-4b49-b362-9929007425ae",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return tempTestCaseBasicInformationStruct{}, err
		}

	}

	if numberOfRowFromDB > 1 {
		numberOfRowFromDB = 2
	}

	switch numberOfRowFromDB {
	case 0:

		err := errors.New(fmt.Sprintf("expected one row from datavase but got zero rows for testcase: %s", testCaseUuid))
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "91bf8bf6-9c03-433c-8125-e08efe8ccb2d",
			"testCaseUuid": testCaseUuid,
			"sqlToExecute": sqlToExecute,
		}).Error("Expected 1 row but got zero rows")

		return tempTestCaseBasicInformationStruct{}, err

	case 1:

	case 2:
		err := errors.New(fmt.Sprintf("expected exactly one row from database but got more then one rows for testcase: %s", testCaseUuid))
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "f91489e5-78fe-4cca-95b4-f1a102eaf6cc",
			"testCaseUuid": testCaseUuid,
			"sqlToExecute": sqlToExecute,
		}).Error("Expected 1 row but got more then 1 rows")

		return tempTestCaseBasicInformationStruct{}, err

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
