package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

// Temporary structure for handling TestInstructionExecutions and their LogPosts and expected and found values
type workObjectForTestInstructionExecutionsMessageStruct struct {
	TestInstructionExecutionBasicInformation *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionBasicInformationMessage
	TestInstructionExecutionsInformation     []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage
	ExecutionLogPostsAndValues               []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
	RunTimeUpdatedAttributes                 []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage
}

// Temporary structure for handling TestCaseExecutions and references to TestInstructionExecutions
type workObjectForTestCaseExecutionResponseMessageStruct struct {
	TestCaseExecutionBasicInformation *fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage
	TestCaseExecutionDetails          []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
	TestInstructionExecutionsMap      map[string]*workObjectForTestInstructionExecutionsMessageStruct // map[TestInstructionExecutionKey]*workObjectForTestInstructionExecutionsMessageStruct
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadFullTestCasesExecutionInformation(
	testCaseExecutionKeys []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionKeyMessage) (
	testCaseExecutionResponseMessages []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "abd88f6a-e916-45ed-97a0-2c3a02eef6f5",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin' in 'loadFullTestCasesExecutionInformation'")

		return testCaseExecutionResponseMessages, nil
	}

	// Close db-transaction when leaving this function
	defer txn.Commit(context.Background())

	// Map for keep track of all response messages, but in Map-format instead of slice-format
	// map[TestCaseExecutionKey]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
	var tempTestCaseExecutionResponseMessagesMap map[string]*workObjectForTestCaseExecutionResponseMessageStruct
	tempTestCaseExecutionResponseMessagesMap = make(map[string]*workObjectForTestCaseExecutionResponseMessageStruct)

	// Convert 'TestCaseExecutionKeys' into slice with 'UniqueCounter' for table 'TestCaseExecutionQueue'
	var uniqueCountersForTableTestCaseExecutionQueue []int
	uniqueCountersForTableTestCaseExecutionQueue, err = fenixGuiTestCaseBuilderServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestCaseExecutionQueue")

	// If there are no TestCases under onQueue then ignore this part
	if uniqueCountersForTableTestCaseExecutionQueue != nil {
		// Load TestCaseExecutions from table 'TestCaseExecutionQueue'
		_, err = fenixGuiTestCaseBuilderServerObject.loadTestCasesExecutionsFromOnExecutionQueue(
			txn,
			uniqueCountersForTableTestCaseExecutionQueue,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Convert 'TestCaseExecutionKeys' into slice with 'UniqueCounter' for table 'TestCasesUnderExecution'
	var uniqueCountersForTableTTestCasesUnderExecution []int
	uniqueCountersForTableTTestCasesUnderExecution, err = fenixGuiTestCaseBuilderServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestCasesUnderExecution")

	// If there are no TestCases under Execution then ignore this part
	if uniqueCountersForTableTTestCasesUnderExecution != nil {

		// Load TestCaseExecutions from table 'TestCasesUnderExecution'
		_, err = fenixGuiTestCaseBuilderServerObject.loadTestCasesExecutionsFromUnderExecutions(
			txn,
			uniqueCountersForTableTTestCasesUnderExecution,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Convert 'TestInstructionExecutionKeys' into slice with 'UniqueCounter' for table 'TestInstructionExecutionQueue'
	var uniqueCountersForTableTestInstructionExecutionQueue []int
	uniqueCountersForTableTestInstructionExecutionQueue, err = fenixGuiTestCaseBuilderServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestInstructionExecutionQueue")

	// Only process when there still are TestInstructionExecution on the ExecutionQueue
	if len(uniqueCountersForTableTestInstructionExecutionQueue) > 0 {

		// Load TestInstructionExecutions from table 'TestInstructionExecutionQueue'
		_, err = fenixGuiTestCaseBuilderServerObject.loadTestInstructionsExecutionsFromOnExecutionQueue(
			txn,
			uniqueCountersForTableTestInstructionExecutionQueue,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Convert 'TestInstructionExecutionKeys' into slice with 'UniqueCounter' for table 'TestInstructionsUnderExecution'
	var uniqueCountersForTableTestInstructionsUnderExecution []int
	uniqueCountersForTableTestInstructionsUnderExecution, err = fenixGuiTestCaseBuilderServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestInstructionsUnderExecution")

	// Only process when there still are TestInstructionExecution on the ExecutionQueue
	if len(uniqueCountersForTableTestInstructionsUnderExecution) > 0 {

		// Load TestInstructionExecutions from table 'TestInstructionsUnderExecution'
		_, err = fenixGuiTestCaseBuilderServerObject.loadTestInstructionsExecutionsUnderExecution(
			txn,
			uniqueCountersForTableTestInstructionsUnderExecution,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Convert 'tempTestCaseExecutionResponseMessagesMap' into gRPC-response object
	err = fenixGuiTestCaseBuilderServerObject.convertTestCaseExecutionResponseMessagesMapIntoGrpcResponse(
		&tempTestCaseExecutionResponseMessagesMap,
		&testCaseExecutionResponseMessages)

	if err != nil {
		return nil, err
	}

	return testCaseExecutionResponseMessages, err
}

// Convert 'TestCaseExecutionKeys' (TestCaseExecutionUuid + TestCaseExecutionVersion) into a slice with 'UniqueCounter' which are unique number for every DB-row in table
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadUniqueCountersBasedOnTestCaseExecutionKeys(
	dbTransaction pgx.Tx,
	TestCaseExecutionKeys []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionKeyMessage,
	databaseTableName string) (
	uniqueCounters []int,
	err error) {

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT \"UniqueCounter\" "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"" + databaseTableName + "\" "

	// if TestCaseExecutionKeysList has 'TestCaseExecutionKeys' then add that as Where-statement
	if TestCaseExecutionKeys != nil {
		for TestCaseExecutionKeyCounter, TestCaseExecutionKey := range TestCaseExecutionKeys {
			if TestCaseExecutionKeyCounter == 0 {
				// Add 'Where' for the first TestCaseExecutionKey, otherwise add an 'ADD'
				sqlToExecute = sqlToExecute + "WHERE "
			} else {
				sqlToExecute = sqlToExecute + "OR "
			}

			sqlToExecute = sqlToExecute + "\"TestCaseExecutionUuid\" = '" + TestCaseExecutionKey.TestCaseExecutionUuid + "' "
			sqlToExecute = sqlToExecute + "AND "
			sqlToExecute = sqlToExecute + "\"TestCaseExecutionVersion\" = " + strconv.FormatUint(uint64(TestCaseExecutionKey.TestCaseExecutionVersion), 10)
			sqlToExecute = sqlToExecute + " "
		}
	}

	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "5c072bd9-da0d-457d-81fa-f6437a6fd81c",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Variables to used when extract data from result set

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		var tempUniqueCounter int

		err := rows.Scan(
			&tempUniqueCounter,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "6edc8e52-0411-4c22-b93f-f608784b85cb",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Add 'tempUniqueCounter' to  slice of UniqueCounters
		uniqueCounters = append(uniqueCounters, tempUniqueCounter)

	}

	return uniqueCounters, err
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestCasesExecutionsFromOnExecutionQueue(
	dbTransaction pgx.Tx,
	uniqueCounters []int,
	tempTestCaseExecutionResponseMessagesMapReference *map[string]*workObjectForTestCaseExecutionResponseMessageStruct) (
	numberOfRows int,
	err error) {

	// Convert reference into variable to use
	tempTestCaseExecutionResponseMessagesMap := *tempTestCaseExecutionResponseMessagesMapReference

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TCEQ.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCaseExecutionQueue\" TCEQ "

	// if uniqueCounters has values then add that as Where-statement
	if uniqueCounters != nil {
		sqlToExecute = sqlToExecute + "WHERE TCEQ.\"UniqueCounter\" IN " +
			fenixGuiTestCaseBuilderServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b041cb41-8e3b-4f87-922a-09f23fbb253e",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return 0, err
	}

	// Variables used for 'tempTestCaseExecutionResponseMessagesMap'
	var testCaseExecutionMapKey string
	var existsInMap bool

	// Variables to used when extract data from result set
	var tempPlacedOnTestExecutionQueueTimeStamp time.Time
	var tempExecutionPriority int
	var tempUniqueCounter int
	var tempExecutionStatusReportLevelEnum int

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		testCaseExecutionBasicInformation := fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage{}

		err := rows.Scan(
			&testCaseExecutionBasicInformation.DomainUuid,
			&testCaseExecutionBasicInformation.DomainName,
			&testCaseExecutionBasicInformation.TestSuiteUuid,
			&testCaseExecutionBasicInformation.TestSuiteName,
			&testCaseExecutionBasicInformation.TestSuiteVersion,
			&testCaseExecutionBasicInformation.TestSuiteExecutionUuid,
			&testCaseExecutionBasicInformation.TestSuiteExecutionVersion,
			&testCaseExecutionBasicInformation.TestCaseUuid,
			&testCaseExecutionBasicInformation.TestCaseName,
			&testCaseExecutionBasicInformation.TestCaseVersion,
			&testCaseExecutionBasicInformation.TestCaseExecutionUuid,
			&testCaseExecutionBasicInformation.TestCaseExecutionVersion,
			&tempPlacedOnTestExecutionQueueTimeStamp,
			&testCaseExecutionBasicInformation.TestDataSetUuid,
			&tempExecutionPriority,
			&tempUniqueCounter,
			&tempExecutionStatusReportLevelEnum,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "030eeab7-5bd0-4013-83f4-3a36d9267c64",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return 0, err
		}

		// Convert temp-variables into gRPC-variables
		testCaseExecutionBasicInformation.PlacedOnTestExecutionQueueTimeStamp = timestamppb.New(tempPlacedOnTestExecutionQueueTimeStamp)
		testCaseExecutionBasicInformation.ExecutionPriority = fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(tempExecutionPriority)

		// Create 'testCaseExecutionMapKey'
		testCaseExecutionMapKey = testCaseExecutionBasicInformation.TestCaseExecutionUuid + strconv.FormatUint(uint64(testCaseExecutionBasicInformation.TestCaseExecutionVersion), 10)

		// Check if data exist for testCaseExecutionMapKey
		var tempWorkObjectForTestCaseExecutionResponseMessage *workObjectForTestCaseExecutionResponseMessageStruct
		tempWorkObjectForTestCaseExecutionResponseMessage, existsInMap = tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

		if existsInMap == false {
			// Initiate all variables
			/*
				var tempFoundVersusExpectedValue *fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage_FoundVersusExpectedValueMessage
				tempFoundVersusExpectedValue = &fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage_FoundVersusExpectedValueMessage{
					FoundValue:    "",
					ExpectedValue: "",
				}

				var  tempLogPostAndValuesMessage  *fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
				tempLogPostAndValuesMessage  = &fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage{
					TestInstructionExecutionUuid:    "",
					TestInstructionExecutionVersion: 0,
					LogPostTimeStamp:                nil,
					LogPostStatus:                   0,
					FoundVersusExpectedValue:        tempFoundVersusExpectedValue,
				}

				var tempExecutionLogPostsAndValues *fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
				tempExecutionLogPostsAndValues = &fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage{
					TestInstructionExecutionUuid:    "",
					TestInstructionExecutionVersion: 0,
					LogPostTimeStamp:                nil,
					LogPostStatus:                   0,
					FoundVersusExpectedValue:        nil,
				}



				var tempTestInstructionExecution *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage
				tempTestInstructionExecution = &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage{
					TestInstructionExecutionBasicInformation: nil,
					TestInstructionExecutionsInformation:     nil,
					ExecutionLogPostsAndValues:               nil,
				}
				var tempTestInstructionExecutions []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage
				tempTestInstructionExecutions = append()

			*/

			// Create a fictive TestCaseExecutionStatus-message to represent that it TestCaseExecution is on ExecutionsQueue
			var tempTestCaseExecutionsInformationMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
			tempTestCaseExecutionsInformationMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage{
				ExecutionStartTimeStamp:        nil,
				ExecutionStopTimeStamp:         nil,
				TestCaseExecutionStatus:        fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_INITIATED,
				ExecutionHasFinished:           false,
				ExecutionStatusUpdateTimeStamp: testCaseExecutionBasicInformation.PlacedOnTestExecutionQueueTimeStamp,
				UniqueDatabaseRowCounter:       0,
			}

			var tempTestCaseExecutionDetails []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
			tempTestCaseExecutionDetails = []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage{
				tempTestCaseExecutionsInformationMessage}

			// Initiate object to be stored in 'tempTestCaseExecutionResponseMessagesMap'
			tempWorkObjectForTestCaseExecutionResponseMessage = &workObjectForTestCaseExecutionResponseMessageStruct{
				TestCaseExecutionBasicInformation: &testCaseExecutionBasicInformation,
				TestCaseExecutionDetails:          tempTestCaseExecutionDetails,
				TestInstructionExecutionsMap:      make(map[string]*workObjectForTestInstructionExecutionsMessageStruct),
			}

			// Add 'tempWorkObjectForTestCaseExecutionResponseMessage' to 'tempTestCaseExecutionResponseMessagesMap'
			tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey] = tempWorkObjectForTestCaseExecutionResponseMessage

		} else {

			// Add to existing 'tempTestCaseExecutionResponseMessage'
			tempWorkObjectForTestCaseExecutionResponseMessage.TestCaseExecutionBasicInformation = &testCaseExecutionBasicInformation
		}

		// Add to number of rows
		numberOfRows = numberOfRows + 1

	}

	return numberOfRows, err
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestCasesExecutionsFromUnderExecutions(
	dbTransaction pgx.Tx,
	uniqueCounters []int,
	tempTestCaseExecutionResponseMessagesMapReference *map[string]*workObjectForTestCaseExecutionResponseMessageStruct) (
	numberOfRows int,
	err error) {

	// Convert reference into variable to use
	tempTestCaseExecutionResponseMessagesMap := *tempTestCaseExecutionResponseMessagesMapReference

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TCUE.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCasesUnderExecution\" TCUE "

	// if uniqueCounters has values then add that as Where-statement
	if uniqueCounters != nil {
		sqlToExecute = sqlToExecute + "WHERE TCUE.\"UniqueCounter\" IN " +
			fenixGuiTestCaseBuilderServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "98b552ed-1031-42da-a5a9-287e542abfb1",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return 0, err
	}

	// Variables used for 'tempTestCaseExecutionResponseMessagesMap'
	var testCaseExecutionMapKey string
	var existsInMap bool

	// Variables to used when extract data from result set
	var tempPlacedOnTestExecutionQueueTimeStamp time.Time
	var tempExecutionPriority int

	var tempExecutionStartTimeStamp time.Time
	var tempExecutionStopTimeStamp time.Time
	var tempTestCaseExecutionStatus int
	var tempExecutionStatusUpdateTimeStamp time.Time
	var tempExecutionStatusReportLevelEnum int

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		var tempTestCaseExecutionBasicInformationMessage fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage
		var tempTestCaseExecutionDetailsMessage fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage

		err := rows.Scan(
			// TestCaseExecutionBasicInformationMessage
			&tempTestCaseExecutionBasicInformationMessage.DomainUuid,
			&tempTestCaseExecutionBasicInformationMessage.DomainName,
			&tempTestCaseExecutionBasicInformationMessage.TestSuiteUuid,
			&tempTestCaseExecutionBasicInformationMessage.TestSuiteName,
			&tempTestCaseExecutionBasicInformationMessage.TestSuiteVersion,
			&tempTestCaseExecutionBasicInformationMessage.TestSuiteExecutionUuid,
			&tempTestCaseExecutionBasicInformationMessage.TestSuiteExecutionVersion,
			&tempTestCaseExecutionBasicInformationMessage.TestCaseUuid,
			&tempTestCaseExecutionBasicInformationMessage.TestCaseName,
			&tempTestCaseExecutionBasicInformationMessage.TestCaseVersion,
			&tempTestCaseExecutionBasicInformationMessage.TestCaseExecutionUuid,
			&tempTestCaseExecutionBasicInformationMessage.TestCaseExecutionVersion,
			&tempPlacedOnTestExecutionQueueTimeStamp,
			&tempTestCaseExecutionBasicInformationMessage.TestDataSetUuid,
			&tempExecutionPriority,

			// TestCaseExecutionDetailsMessage
			&tempExecutionStartTimeStamp,
			&tempExecutionStopTimeStamp,
			&tempTestCaseExecutionStatus,
			&tempTestCaseExecutionDetailsMessage.ExecutionHasFinished,
			&tempTestCaseExecutionDetailsMessage.UniqueDatabaseRowCounter,
			&tempExecutionStatusUpdateTimeStamp,

			// ReportLevel
			&tempExecutionStatusReportLevelEnum,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "61ca3d9d-bc80-4702-873f-48f62bfcadb1",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return 0, err
		}

		// Convert temp-variables into gRPC-variables
		tempTestCaseExecutionBasicInformationMessage.PlacedOnTestExecutionQueueTimeStamp = timestamppb.New(tempPlacedOnTestExecutionQueueTimeStamp)
		tempTestCaseExecutionBasicInformationMessage.ExecutionPriority = fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(tempExecutionPriority)

		tempTestCaseExecutionDetailsMessage.ExecutionStartTimeStamp = timestamppb.New(tempExecutionStartTimeStamp)
		tempTestCaseExecutionDetailsMessage.ExecutionStopTimeStamp = timestamppb.New(tempExecutionStopTimeStamp)
		tempTestCaseExecutionDetailsMessage.TestCaseExecutionStatus = fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum(tempTestCaseExecutionStatus)
		tempTestCaseExecutionDetailsMessage.ExecutionStatusUpdateTimeStamp = timestamppb.New(tempExecutionStatusUpdateTimeStamp)

		// Create 'testCaseExecutionMapKey'
		testCaseExecutionMapKey = tempTestCaseExecutionBasicInformationMessage.TestCaseExecutionUuid + strconv.FormatUint(uint64(tempTestCaseExecutionBasicInformationMessage.TestCaseExecutionVersion), 10)

		// Check if data exist for testCaseExecutionMapKey
		var tempTestCaseExecutionResponseMessage *workObjectForTestCaseExecutionResponseMessageStruct
		tempTestCaseExecutionResponseMessage, existsInMap = tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

		if existsInMap == false {
			// Initiate object to be stored in 'tempTestCaseExecutionResponseMessagesMap'

			var tempTestCaseExecutionDetailsMessageSlice []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
			tempTestCaseExecutionDetailsMessageSlice = append(tempTestCaseExecutionDetailsMessageSlice, &tempTestCaseExecutionDetailsMessage)

			tempTestCaseExecutionResponseMessage = &workObjectForTestCaseExecutionResponseMessageStruct{
				TestCaseExecutionBasicInformation: &tempTestCaseExecutionBasicInformationMessage,
				TestCaseExecutionDetails:          tempTestCaseExecutionDetailsMessageSlice,
				TestInstructionExecutionsMap:      make(map[string]*workObjectForTestInstructionExecutionsMessageStruct)}

			// Add 'tempTestCaseExecutionResponseMessage' to 'tempTestCaseExecutionResponseMessagesMap'
			tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey] = tempTestCaseExecutionResponseMessage

		} else {
			// Append to existing 'tempTestCaseExecutionResponseMessage'
			tempTestCaseExecutionResponseMessage.TestCaseExecutionDetails = append(
				tempTestCaseExecutionResponseMessage.TestCaseExecutionDetails, &tempTestCaseExecutionDetailsMessage)

		}

		// Add to number of rows
		numberOfRows = numberOfRows + 1

	}

	return numberOfRows, err
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestInstructionsExecutionsFromOnExecutionQueue(
	dbTransaction pgx.Tx,
	uniqueCounters []int,
	tempTestCaseExecutionResponseMessagesMapReference *map[string]*workObjectForTestCaseExecutionResponseMessageStruct) (
	numberOfRows int,
	err error) {

	// Convert reference into variable to use
	tempTestCaseExecutionResponseMessagesMap := *tempTestCaseExecutionResponseMessagesMapReference

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TIEQ.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructionExecutionQueue\" TIEQ "

	// if uniqueCounters has values then add that as Where-statement
	if uniqueCounters != nil {
		sqlToExecute = sqlToExecute + "WHERE TIEQ.\"UniqueCounter\" IN " +
			fenixGuiTestCaseBuilderServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "4ac2c057-1a37-47d1-88ad-a37aa7b1153b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return 0, err
	}

	// Variables used for 'tempTestCaseExecutionResponseMessagesMap'
	var testCaseExecutionMapKey string
	var testInstructionExecutionMapKey string
	var existsInMap bool

	// Variables to used when extract data from result set
	var tempPlacedOnTestInstructionExecutionQueueTimeStamp time.Time
	var tempExecutionPriority int
	var tempUniqueCounter int
	var tempExecutionStatusReportLevel fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		testInstructionExecutionBasicInformation := fenixExecutionServerGuiGrpcApi.TestInstructionExecutionBasicInformationMessage{}

		err = rows.Scan(
			&testInstructionExecutionBasicInformation.DomainUuid,
			&testInstructionExecutionBasicInformation.DomainName,
			&testInstructionExecutionBasicInformation.TestInstructionExecutionUuid,
			&testInstructionExecutionBasicInformation.TestInstructionUuid,
			&testInstructionExecutionBasicInformation.TestInstructionName,
			&testInstructionExecutionBasicInformation.TestInstructionMajorVersionNumber,
			&testInstructionExecutionBasicInformation.TestInstructionMinorVersionNumber,
			&tempPlacedOnTestInstructionExecutionQueueTimeStamp,
			&tempExecutionPriority,
			&testInstructionExecutionBasicInformation.TestCaseExecutionUuid,
			&testInstructionExecutionBasicInformation.TestDataSetUuid,
			&testInstructionExecutionBasicInformation.TestCaseExecutionVersion,
			&testInstructionExecutionBasicInformation.TestInstructionExecutionVersion,
			&testInstructionExecutionBasicInformation.TestInstructionExecutionOrder,
			&tempUniqueCounter,
			&testInstructionExecutionBasicInformation.TestInstructionOriginalUuid,
			&tempExecutionStatusReportLevel,
			&testInstructionExecutionBasicInformation.ExecutionDomainUuid,
			&testInstructionExecutionBasicInformation.ExecutionDomainName,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "dc06a877-53d6-4ef1-bffd-af17f27137e7",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return 0, err
		}

		// Convert temp-variables into gRPC-variables
		testInstructionExecutionBasicInformation.QueueTimeStamp = timestamppb.New(tempPlacedOnTestInstructionExecutionQueueTimeStamp)
		testInstructionExecutionBasicInformation.ExecutionPriority = fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(tempExecutionPriority)

		// Create 'testCaseExecutionMapKey'
		testCaseExecutionMapKey = testInstructionExecutionBasicInformation.TestCaseExecutionUuid + strconv.FormatUint(uint64(testInstructionExecutionBasicInformation.TestCaseExecutionVersion), 10)

		// Create 'testInstructionExecutionMapKey'
		testInstructionExecutionMapKey = testInstructionExecutionBasicInformation.TestInstructionExecutionUuid + strconv.FormatUint(uint64(testInstructionExecutionBasicInformation.TestInstructionExecutionVersion), 10)

		// Check if data exist for 'testInstructionExecutionMapKey'
		var tempWorkObjectForTestCaseExecutionResponseMessage *workObjectForTestCaseExecutionResponseMessageStruct
		tempWorkObjectForTestCaseExecutionResponseMessage, existsInMap = tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

		if existsInMap == false {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                             "6ea5ed57-b015-4fca-bee4-26355b2df789",
				"testInstructionExecutionMapKey": testInstructionExecutionMapKey,
			}).Error("Couldn't find 'testCaseExecutionMapKey' in 'tempTestCaseExecutionResponseMessagesMap'")

			return 0, err
		}

		// Initiate object to be stored in 'TestInstructionExecutionsMap'
		var tempWorkObjectForTestInstructionExecutionsMessage *workObjectForTestInstructionExecutionsMessageStruct
		tempWorkObjectForTestInstructionExecutionsMessage, existsInMap =
			tempWorkObjectForTestCaseExecutionResponseMessage.TestInstructionExecutionsMap[testInstructionExecutionMapKey]

		// If 'testInstructionExecutionMapKey' doesn't exist then create the object
		if existsInMap == false {

			// Create a fictive TestInstructionExecution-message to represent that it TestInstructionExecution is on ExecutionsQueue
			var tempTestInstructionExecutionsInformationMessage *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage
			tempTestInstructionExecutionsInformationMessage = &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage{
				SentTimeStamp:                        nil,
				ExpectedExecutionEndTimeStamp:        nil,
				TestInstructionExecutionStatus:       fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum_TIE_INITIATED,
				TestInstructionExecutionEndTimeStamp: nil,
				TestInstructionExecutionHasFinished:  false,
				UniqueDatabaseRowCounter:             0,
				TestInstructionCanBeReExecuted:       false,
				ExecutionStatusUpdateTimeStamp:       testInstructionExecutionBasicInformation.QueueTimeStamp,
			}
			var tempTestInstructionExecutions []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage
			tempTestInstructionExecutions = []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage{
				tempTestInstructionExecutionsInformationMessage}

			tempWorkObjectForTestInstructionExecutionsMessage = &workObjectForTestInstructionExecutionsMessageStruct{
				TestInstructionExecutionBasicInformation: &testInstructionExecutionBasicInformation,
				TestInstructionExecutionsInformation:     tempTestInstructionExecutions,
				ExecutionLogPostsAndValues:               nil,
				RunTimeUpdatedAttributes:                 nil,
			}

			// Add 'tempWorkObjectForTestInstructionExecutionsMessage' to 'TestInstructionExecutionsMap'
			tempWorkObjectForTestCaseExecutionResponseMessage.TestInstructionExecutionsMap[testInstructionExecutionMapKey] =
				tempWorkObjectForTestInstructionExecutionsMessage
		} else {

			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                             "1f02fa15-200e-4cb9-8248-3a57f27242dc",
				"testInstructionExecutionMapKey": testInstructionExecutionMapKey,
				"tempWorkObjectForTestCaseExecutionResponseMessage.TestInstructionExecutionsMap": tempWorkObjectForTestCaseExecutionResponseMessage.TestInstructionExecutionsMap,
			}).Fatalln("We shouldn't come here")
		}

		// Add to number of rows
		numberOfRows = numberOfRows + 1

	}

	return numberOfRows, err
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestInstructionsExecutionsUnderExecution(
	dbTransaction pgx.Tx,
	uniqueCounters []int,
	tempTestCaseExecutionResponseMessagesMapReference *map[string]*workObjectForTestCaseExecutionResponseMessageStruct) (
	numberOfRows int,
	err error) {

	// Convert reference into variable to use
	tempTestCaseExecutionResponseMessagesMap := *tempTestCaseExecutionResponseMessagesMapReference

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TIUE.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructionsUnderExecution\" TIUE "

	// if uniqueCounters has values then add that as Where-statement
	if uniqueCounters != nil {
		sqlToExecute = sqlToExecute + "WHERE TIUE.\"UniqueCounter\" IN " +
			fenixGuiTestCaseBuilderServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "4ceaee78-77b3-4da1-9e30-32543989403c",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return 0, err
	}

	// Variables used for 'tempTestCaseExecutionResponseMessagesMap'
	var testCaseExecutionMapKey string
	var testInstructionExecutionMapKey string
	var existsInMap bool
	var tempTestInstructionSendCounter int

	// Variables to used when extract data from result set
	var (
		tempSentTimeStamp                        *time.Time
		tempExpectedExecutionDuration            *time.Time
		tempExpectedExecutionEndTimeStamp        *time.Time
		tempTestInstructionExecutionStatus       int
		tempExecutionStatusUpdateTimeStamp       *time.Time
		tempTestInstructionExecutionEndTimeStamp *time.Time
		tempQueueTimeStamp                       *time.Time
		tempExecutionPriority                    *int
		tempExecutionStatusReportLevel           fenixExecutionServerGuiGrpcApi.ExecutionStatusReportLevelEnum
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		tempTestInstructionExecutionBasicInformation := fenixExecutionServerGuiGrpcApi.TestInstructionExecutionBasicInformationMessage{}
		var tempTestInstructionExecutionsInformationMessage fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage

		err = rows.Scan(
			&tempTestInstructionExecutionBasicInformation.DomainUuid,
			&tempTestInstructionExecutionBasicInformation.DomainName,
			&tempTestInstructionExecutionBasicInformation.TestInstructionExecutionUuid,
			&tempTestInstructionExecutionBasicInformation.TestInstructionUuid,
			&tempTestInstructionExecutionBasicInformation.TestInstructionName,
			&tempTestInstructionExecutionBasicInformation.TestInstructionMajorVersionNumber,
			&tempTestInstructionExecutionBasicInformation.TestInstructionMinorVersionNumber,
			&tempSentTimeStamp,
			&tempExpectedExecutionDuration,
			&tempExpectedExecutionEndTimeStamp,
			&tempTestInstructionExecutionStatus,
			&tempExecutionStatusUpdateTimeStamp,
			&tempTestInstructionExecutionBasicInformation.TestDataSetUuid,
			&tempTestInstructionExecutionBasicInformation.TestCaseExecutionUuid,
			&tempTestInstructionExecutionBasicInformation.TestCaseExecutionVersion,
			&tempTestInstructionExecutionBasicInformation.TestInstructionExecutionVersion,
			&tempTestInstructionExecutionsInformationMessage.TestInstructionCanBeReExecuted,
			&tempTestInstructionExecutionBasicInformation.TestInstructionExecutionOrder,
			&tempTestInstructionExecutionsInformationMessage.UniqueDatabaseRowCounter,
			&tempTestInstructionExecutionBasicInformation.TestInstructionOriginalUuid,
			&tempTestInstructionExecutionEndTimeStamp,
			&tempTestInstructionExecutionsInformationMessage.TestInstructionExecutionHasFinished,
			&tempQueueTimeStamp,
			&tempExecutionPriority,
			&tempTestInstructionSendCounter,
			&tempExecutionStatusReportLevel,
			&tempTestInstructionExecutionBasicInformation.ExecutionDomainUuid,
			&tempTestInstructionExecutionBasicInformation.ExecutionDomainName,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "828b82fb-cce0-42e2-883c-b2011543fb96",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return 0, err
		}

		// Convert temp-variables into gRPC-variables
		if tempSentTimeStamp != nil {
			tempTestInstructionExecutionsInformationMessage.SentTimeStamp =
				timestamppb.New(*tempSentTimeStamp)
		}
		if tempExpectedExecutionEndTimeStamp != nil {
			tempTestInstructionExecutionsInformationMessage.ExpectedExecutionEndTimeStamp =
				timestamppb.New(*tempExpectedExecutionEndTimeStamp)
		}
		tempTestInstructionExecutionsInformationMessage.TestInstructionExecutionStatus =
			fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum(tempTestInstructionExecutionStatus)
		if tempExecutionStatusUpdateTimeStamp != nil {
			tempTestInstructionExecutionsInformationMessage.ExecutionStatusUpdateTimeStamp =
				timestamppb.New(*tempExecutionStatusUpdateTimeStamp)
		}
		if tempTestInstructionExecutionEndTimeStamp != nil {
			tempTestInstructionExecutionsInformationMessage.TestInstructionExecutionEndTimeStamp =
				timestamppb.New(*tempTestInstructionExecutionEndTimeStamp)
		}
		if tempQueueTimeStamp != nil {
			tempTestInstructionExecutionBasicInformation.QueueTimeStamp =
				timestamppb.New(*tempQueueTimeStamp)
		}
		if tempExecutionPriority != nil {
			tempTestInstructionExecutionBasicInformation.ExecutionPriority =
				fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(*tempExecutionPriority)
		}

		// Create 'testCaseExecutionMapKey'
		testCaseExecutionMapKey =
			tempTestInstructionExecutionBasicInformation.TestCaseExecutionUuid +
				strconv.FormatUint(uint64(tempTestInstructionExecutionBasicInformation.TestCaseExecutionVersion), 10)

		// Check if data exist for 'testInstructionExecutionMapKey'
		var tempWorkObjectForTestCaseExecutionResponseMessage *workObjectForTestCaseExecutionResponseMessageStruct
		tempWorkObjectForTestCaseExecutionResponseMessage, existsInMap =
			tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

		if existsInMap == false {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                             "55ff90f3-6ac2-4c8a-ae34-ad008ccb02a8",
				"testInstructionExecutionMapKey": testInstructionExecutionMapKey,
				"sqlToExecute":                   sqlToExecute,
			}).Error("Couldn't find 'testCaseExecutionMapKey' in 'tempTestCaseExecutionResponseMessagesMap'")

			return 0, errors.New("couldn't find 'testCaseExecutionMapKey' in 'tempTestCaseExecutionResponseMessagesMap")
		}
		// Initiate all variables
		/*
			var tempFoundVersusExpectedValue *fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage_FoundVersusExpectedValueMessage
			tempFoundVersusExpectedValue = &fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage_FoundVersusExpectedValueMessage{
				FoundValue:    "",
				ExpectedValue: "",
			}

			var  tempLogPostAndValuesMessage  *fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
			tempLogPostAndValuesMessage  = &fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage{
				TestInstructionExecutionUuid:    "",
				TestInstructionExecutionVersion: 0,
				LogPostTimeStamp:                nil,
				LogPostStatus:                   0,
				FoundVersusExpectedValue:        tempFoundVersusExpectedValue,
			}

			var tempExecutionLogPostsAndValues *fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
			tempExecutionLogPostsAndValues = &fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage{
				TestInstructionExecutionUuid:    "",
				TestInstructionExecutionVersion: 0,
				LogPostTimeStamp:                nil,
				LogPostStatus:                   0,
				FoundVersusExpectedValue:        nil,
			}
		*/

		// Create 'testInstructionExecutionMapKey'
		testInstructionExecutionMapKey =
			tempTestInstructionExecutionBasicInformation.TestInstructionExecutionUuid +
				strconv.FormatUint(uint64(tempTestInstructionExecutionBasicInformation.TestInstructionExecutionVersion), 10)

		// Initiate object to be stored in 'TestInstructionExecutionsMap'
		var tempWorkObjectForTestInstructionExecutionsMessage *workObjectForTestInstructionExecutionsMessageStruct
		tempWorkObjectForTestInstructionExecutionsMessage, existsInMap =
			tempWorkObjectForTestCaseExecutionResponseMessage.TestInstructionExecutionsMap[testInstructionExecutionMapKey]

		// If 'testInstructionExecutionMapKey' doesn't exist then create it
		if existsInMap == false {

			// Create a fictive TestInstructionExecution-message to represent that it TestInstructionExecution is on ExecutionsQueue
			var tempTestInstructionExecutionsInformationForOnQueMessage *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage
			tempTestInstructionExecutionsInformationForOnQueMessage = &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage{
				SentTimeStamp:                        nil,
				ExpectedExecutionEndTimeStamp:        nil,
				TestInstructionExecutionStatus:       fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum_TIE_INITIATED,
				TestInstructionExecutionEndTimeStamp: nil,
				TestInstructionExecutionHasFinished:  false,
				UniqueDatabaseRowCounter:             0,
				TestInstructionCanBeReExecuted:       false,
				ExecutionStatusUpdateTimeStamp:       tempTestInstructionExecutionBasicInformation.QueueTimeStamp,
			}

			// Create slice for TestInstructionExecutions
			var tempTestInstructionExecutions []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage
			tempTestInstructionExecutions = []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage{
				tempTestInstructionExecutionsInformationForOnQueMessage,
				&tempTestInstructionExecutionsInformationMessage}

			// Add TestInstructionsExecution-information to message to be stored in Map
			var tempTestInstructionExecution *workObjectForTestInstructionExecutionsMessageStruct
			tempTestInstructionExecution = &workObjectForTestInstructionExecutionsMessageStruct{
				TestInstructionExecutionBasicInformation: &tempTestInstructionExecutionBasicInformation,
				TestInstructionExecutionsInformation:     tempTestInstructionExecutions,
				ExecutionLogPostsAndValues:               nil,
				RunTimeUpdatedAttributes:                 nil,
			}

			// Add back to Map
			tempWorkObjectForTestCaseExecutionResponseMessage.TestInstructionExecutionsMap[testInstructionExecutionMapKey] =
				tempTestInstructionExecution
		} else {

			// Append to existing data
			tempWorkObjectForTestInstructionExecutionsMessage.TestInstructionExecutionsInformation = append(
				tempWorkObjectForTestInstructionExecutionsMessage.TestInstructionExecutionsInformation,
				&tempTestInstructionExecutionsInformationMessage)

		}

		// Add to number of rows
		numberOfRows = numberOfRows + 1

	}

	return numberOfRows, err
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestCaseExecutionLogs(
	dbTransaction pgx.Tx,
	testCaseExecutionUuid string) (
	runTimeUpdatedAttributes []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage,
	err error) {

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT ELP.* "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"ExecutionLogPosts\" ELP "
	sqlToExecute = sqlToExecute + "WHERE  ELP.\"TestCaseExecutionUuid\" = '" + testCaseExecutionUuid + "' "
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "e174dc5f-bc65-4ead-b9f0-530e1631e565",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage{}, err
	}

	// Variables to used when extract data from result set
	var (
		tempUpdatedTimeStamp *time.Time
		tempUniqueId         int
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		var tempRunTimeUpdatedAttribute *fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage

		/*
			create table "FenixExecution"."ExecutionLogPosts"
			(
			    "DomainUuid"                      uuid      not null,
			    "TestCaseExecutionUuid"           uuid      not null,
			    "TestCaseExecutionVersion"        integer   not null,
			    "TestInstructionExecutionUuid"    uuid      not null,
			    "TestInstructionExecutionVersion" integer   not null,
			    "TestInstructionExecutionStatus"  integer
			        constraint executionlogposts_testcaseexecutionstatusenum_grpc_id_fk
			            references "FenixExecution"."TestCaseExecutionStatusEnum",
			    "LogPostUuid"                     uuid      not null
			        primary key,
			    "LogPostTimeStamp"                timestamp not null,
			    "LogPostStatus"                   integer
			        constraint executionlogposts_logpoststatusenum_grpc_id_fk
			            references "FenixExecution"."LogPostStatusEnum",
			    "LogPostText"                     varchar,
			    "FoundVsExpectedValuesAsJsonb"    jsonb     not null
			);

		*/

		err = rows.Scan(
			&tempRunTimeUpdatedAttribute.TestInstructionExecutionUuid,
			&tempRunTimeUpdatedAttribute.TestInstructionAttributeType,
			&tempRunTimeUpdatedAttribute.TestInstructionAttributeUuid,
			&tempRunTimeUpdatedAttribute.TestInstructionAttributeName,
			&tempRunTimeUpdatedAttribute.AttributeValueAsString,
			&tempRunTimeUpdatedAttribute.AttributeValueUuid,
			&tempRunTimeUpdatedAttribute.TestInstructionAttributeTypeUuid,
			&tempRunTimeUpdatedAttribute.TestInstructionAttributeTypeName,
			&tempRunTimeUpdatedAttribute.TestInstructionExecutionVersion,
			&tempUpdatedTimeStamp,
			&tempUniqueId,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "b68560d1-3019-4771-84c4-fc24dc208cca",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage{}, err
		}

		// Convert temp-variables into gRPC-variables
		if tempUpdatedTimeStamp != nil {
			tempRunTimeUpdatedAttribute.UpdateTimeStamp =
				timestamppb.New(*tempUpdatedTimeStamp)
		}

		// Add the attribute to the slice of attributes
		runTimeUpdatedAttributes = append(runTimeUpdatedAttributes, tempRunTimeUpdatedAttribute)

	}

	return runTimeUpdatedAttributes, err
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) convertTestCaseExecutionResponseMessagesMapIntoGrpcResponse(
	tempTestCaseExecutionResponseMessagesMapReference *map[string]*workObjectForTestCaseExecutionResponseMessageStruct,
	testCaseExecutionResponseMessagesReference *[]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage) (
	err error) {

	// Loop over TestCaseExecutions in Map
	for _, testCaseExecution := range *tempTestCaseExecutionResponseMessagesMapReference {

		// Create slice for the TestInstructionExecutions within this TestCaseExecution
		var tempTestInstructionExecutions []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage

		// Extract TestInstructionMap
		var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct
		tempTestInstructionExecutionsMap = testCaseExecution.TestInstructionExecutionsMap

		// Loop over TestInstructionExecutions
		for _, testInstructionExecution := range tempTestInstructionExecutionsMap {

			// Create the TestInstructionExecution object to be added
			var tempTestInstructionExecutionsMessage *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage
			tempTestInstructionExecutionsMessage = &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage{
				TestInstructionExecutionBasicInformation: testInstructionExecution.TestInstructionExecutionBasicInformation,
				TestInstructionExecutionsInformation:     testInstructionExecution.TestInstructionExecutionsInformation,
				ExecutionLogPostsAndValues:               nil,
				RunTimeUpdatedAttributes:                 nil,
			}

			// Append TestInstructionExecution to Slice of all TestInstructionExecutions fur current TestCaseExecution
			tempTestInstructionExecutions = append(tempTestInstructionExecutions, tempTestInstructionExecutionsMessage)

		}
		// Create TestCaseExecution object to be added
		var tempTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
		tempTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage{
			TestCaseExecutionBasicInformation: testCaseExecution.TestCaseExecutionBasicInformation,
			TestCaseExecutionDetails:          testCaseExecution.TestCaseExecutionDetails,
			TestInstructionExecutions:         tempTestInstructionExecutions,
		}

		// Append TestCaseExecution to Slice of all TestCaseExecutions for current gRPC-response object
		*testCaseExecutionResponseMessagesReference = append(*testCaseExecutionResponseMessagesReference, tempTestCaseExecutionResponseMessage)
	}

	return err
}
