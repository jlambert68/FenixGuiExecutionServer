package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"errors"
	"fmt"
	uuidGenerator "github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

// After all stuff is done, then Commit or Rollback depending on result
var doCommitNotRoleBackInitiateTestSuiteExecution bool

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) commitOrRoleBackInitiateTestSuiteExecution(dbTransaction pgx.Tx) {
	if doCommitNotRoleBackInitiateTestSuiteExecution == true {
		dbTransaction.Commit(context.Background())
	} else {
		dbTransaction.Rollback(context.Background())
	}
}

// Prepare for Saving the Initiation of a new TestCaseExecution in the CloudDB
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) prepareInitiateTestSuiteExecutionSaveToCloudDB(
	initiateSingleTestSuiteExecutionRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateTestSuiteExecutionWithOneTestDataSetRequestMessage) (
	initiateSingleTestSuiteExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestSuiteExecutionResponseMessage) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "455963ec-77d8-4f99-b279-8a56e644ada1",
		"initiateSingleTestCaseExecutionRequestMessage": initiateSingleTestSuiteExecutionRequestMessage,
	}).Debug("Incoming 'prepareInitiateTestSuiteExecutionSaveToCloudDB'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "d4efc648-3a60-4bbc-8286-f7fcb38d5b6c",
	}).Debug("Outgoing 'prepareInitiateTestSuiteExecutionSaveToCloudDB'")

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":    "627c35c6-235c-4836-82cc-88331a9b7d2f",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareInitiateTestSuiteExecutionSaveToCloudDB'")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		initiateSingleTestSuiteExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestSuiteExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Problem when saving to database",
				ErrorCodes: errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.
					CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
		}

		return initiateSingleTestSuiteExecutionResponseMessage
	}

	// Standard is to do a Rollback
	doCommitNotRoleBackInitiateTestSuiteExecution = false
	defer fenixGuiTestCaseBuilderServerObject.commitOrRoleBackInitiateTestSuiteExecution(txn) //txn.Commit(context.Background())

	// Load TestCases from TestSuite
	var testCasesInTestSuite *fenixTestCaseBuilderServerGrpcApi.TestCasesInTestSuiteMessage
	var tempTestSuiteBasicInformation tempTestSuiteBasicInformationStruct
	testCasesInTestSuite, tempTestSuiteBasicInformation, err = fenixGuiTestCaseBuilderServerObject.loadTestCasesForTestSuite(
		txn,
		initiateSingleTestSuiteExecutionRequestMessage.TestSuiteUuid)

	if err != nil {

		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":    "dbf2a735-7d2e-47ad-8ba6-bc7d40b29ec4",
			"error": err,
		}).Error("Problem when loading TestCases from TestSuite in 'prepareInitiateTestSuiteExecutionSaveToCloudDB'")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		initiateSingleTestSuiteExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestSuiteExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Problem when saving to database",
				ErrorCodes: errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.
					CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
		}

		return initiateSingleTestSuiteExecutionResponseMessage

	}

	// Check if there are no TestCases
	if testCasesInTestSuite == nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id": "dbf2a735-7d2e-47ad-8ba6-bc7d40b29ec4",
			"initiateSingleTestSuiteExecutionRequestMessage.TestSuiteUuid": initiateSingleTestSuiteExecutionRequestMessage.TestSuiteUuid,
		}).Debug("TestSuite has no TestCases")

		// Create Return message
		initiateSingleTestSuiteExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestSuiteExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:    true,
				Comments:   "",
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.
					CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
		}

		return initiateSingleTestSuiteExecutionResponseMessage
	}

	// Create ReturnMessage to be used, later, for when everything goes well
	initiateSingleTestSuiteExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestSuiteExecutionResponseMessage{
		TestCasesInExecutionQueue: nil,
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.
				CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
	}

	// Generate a new TestSuiteExecution-UUID
	testSuiteExecutionUuid := uuidGenerator.New().String()

	// Populate the TestSuite-information to be used when creating the TestCasesExecutions
	var testSuiteInformation testSuiteInformationStruct
	testSuiteInformation = testSuiteInformationStruct{
		suiteUuid:             initiateSingleTestSuiteExecutionRequestMessage.GetTestSuiteUuid(),
		suiteName:             tempTestSuiteBasicInformation.testSuiteName,
		suiteVersion:          uint32(tempTestSuiteBasicInformation.testSuiteVersion),
		suiteExecutionUuid:    testSuiteExecutionUuid,
		suiteExecutionVersion: 1,
	}

	// Loop all TestCases and call 'prepareInitiateTestCaseExecutionSaveToCloudDB'
	for _, tempTestCaseInTestSuite := range testCasesInTestSuite.GetTestCasesInTestSuite() {

		var initiateSingleTestCaseExecutionRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionRequestMessage
		initiateSingleTestCaseExecutionRequestMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionRequestMessage{
			UserAndApplicationRunTimeIdentification: initiateSingleTestSuiteExecutionRequestMessage.UserAndApplicationRunTimeIdentification,
			TestCaseUuid:                            tempTestCaseInTestSuite.GetTestCaseUuid(),
			TestDataSetUuid:                         common_config.ZeroUuid,
			ExecutionStatusReportLevel:              initiateSingleTestSuiteExecutionRequestMessage.GetExecutionStatusReportLevel(),
			TestDataForTestCaseExecution:            initiateSingleTestSuiteExecutionRequestMessage.GetTestDataForTestCaseExecution(),
		}

		var tempInitiateSingleTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage
		tempInitiateSingleTestCaseExecutionResponseMessage = fenixGuiTestCaseBuilderServerObject.prepareInitiateTestCaseExecutionSaveToCloudDB(
			txn,
			initiateSingleTestCaseExecutionRequestMessage,
			fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_HIGH_SINGLE_TESTSUITE,
			testSuiteInformation)

		if tempInitiateSingleTestCaseExecutionResponseMessage.GetAckNackResponse().GetAckNack() == false {

			errMsg := fmt.Sprintf("Problem when saving TestCase: %s in database, from TestSuite: %s",
				tempTestCaseInTestSuite.GetTestCaseUuid(),
				initiateSingleTestSuiteExecutionRequestMessage.GetTestSuiteUuid())

			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"id": "d5ce79ca-0c62-4f4f-977f-7ffaeabd7536",
				"initiateSingleTestSuiteExecutionRequestMessage.TestSuiteUuid": initiateSingleTestSuiteExecutionRequestMessage.TestSuiteUuid,
				"tempTestCaseInTestSuite.TestCaseUuid":                         tempTestCaseInTestSuite.GetTestCaseUuid(),
				"Comments":                                                     tempInitiateSingleTestCaseExecutionResponseMessage.GetAckNackResponse().GetComments(),
			}).Error(errMsg)

			// Create Return message
			initiateSingleTestSuiteExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.InitiateSingleTestSuiteExecutionResponseMessage{
				TestCasesInExecutionQueue: nil,
				AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
					AckNack:    false,
					Comments:   errMsg,
					ErrorCodes: nil,
					ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.
						CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				},
			}

			return initiateSingleTestSuiteExecutionResponseMessage
		}

		// Append single TestCase-response to main TestSuite-response
		initiateSingleTestSuiteExecutionResponseMessage.TestCasesInExecutionQueue = append(
			initiateSingleTestSuiteExecutionResponseMessage.TestCasesInExecutionQueue,
			tempInitiateSingleTestCaseExecutionResponseMessage.GetTestCasesInExecutionQueue())
	}

	// Commit every database change
	doCommitNotRoleBackInitiateTestSuiteExecution = true

	return initiateSingleTestSuiteExecutionResponseMessage
}

// Temporary variable for storing temp result from database
type tempTestSuiteBasicInformationStruct struct {
	testSuiteName    string
	testSuiteVersion int
}

// Load BasicInformation for TestSuite to be able to populate the TestCaseExecution
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestCasesForTestSuite(
	dbTransaction pgx.Tx,
	testSuiteUuid string) (
	_ *fenixTestCaseBuilderServerGrpcApi.TestCasesInTestSuiteMessage,
	_ tempTestSuiteBasicInformationStruct,
	err error) {

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TS.\"TestCasesInTestSuite\", " +
		"TS.\"TestSuiteName\", TS.\"TestSuiteVersion\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" TS "
	sqlToExecute = sqlToExecute + "WHERE TS.\"TestSuiteUuid\" = '" + testSuiteUuid + "' AND "
	sqlToExecute = sqlToExecute + "TS.\"TestSuiteVersion\" = (SELECT MAX(TS2.\"TestSuiteVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" TS2 "
	sqlToExecute = sqlToExecute + "WHERE TS2.\"TestSuiteUuid\" = '" + testSuiteUuid + "');"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "db442978-b183-47d3-9aa7-1792ad21efd1",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadTestCasesForTestSuite'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "8eb6dc50-3261-4321-b295-928e6d36beb1",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, tempTestSuiteBasicInformationStruct{}, err
	}

	// USed to secure that exactly one row was found
	numberOfRowFromDB := 0

	var (
		tempTestSuiteBasicInformation tempTestSuiteBasicInformationStruct

		tempTestCasesInTestSuiteAsJson          string
		tempTestCasesInTestSuiteAsJsonByteArray []byte
		testCasesInTestSuite                    fenixTestCaseBuilderServerGrpcApi.TestCasesInTestSuiteMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		numberOfRowFromDB = numberOfRowFromDB + 1

		err = rows.Scan(
			&tempTestCasesInTestSuiteAsJson,
			&tempTestSuiteBasicInformation.testSuiteName,
			&tempTestSuiteBasicInformation.testSuiteVersion,
		)

		if err != nil {

			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "9cdde993-689a-4b49-b362-9929007425ae",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, tempTestSuiteBasicInformationStruct{}, err
		}

	}

	if numberOfRowFromDB > 1 {
		numberOfRowFromDB = 2
	}

	switch numberOfRowFromDB {
	case 0:
		// TestSuite doesn't have any TestCases

		return nil, tempTestSuiteBasicInformationStruct{}, err

	case 1:

	case 2:
		err = errors.New(fmt.Sprintf("expected exactly one row from database but got more then one rows for TestSuite: %s",
			testSuiteUuid))

		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":            "fc7c2ff1-f6ce-4916-a534-220aae4a3391",
			"testSuiteUuid": testSuiteUuid,
			"sqlToExecute":  sqlToExecute,
		}).Error("Expected 0 or 1 row but got more then 1 rows")

		return nil, tempTestSuiteBasicInformationStruct{}, err

	}

	// Convert json-strings into byte-arrays
	tempTestCasesInTestSuiteAsJsonByteArray = []byte(tempTestCasesInTestSuiteAsJson)

	// Convert json-byte-arrays into proto-messages
	err = protojson.Unmarshal(tempTestCasesInTestSuiteAsJsonByteArray, &testCasesInTestSuite)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "f92e1deb-7754-439a-99c0-8407c61bb6a1",
			"Error": err,
		}).Error("Something went wrong when converting 'tempTestCasesInTestSuiteAsJsonByteArray' into proto-message")

		return nil, tempTestSuiteBasicInformationStruct{}, err
	}

	return &testCasesInTestSuite, tempTestSuiteBasicInformation, err

}
