package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"encoding/json"
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
	TestInstructionExecutionsInformation     *[]*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage
	ExecutionLogPostsAndValues               *[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
	RunTimeUpdatedAttributes                 *[]*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage
}

// Temporary structure for handling TestCaseExecutions and references to TestInstructionExecutions
type workObjectForTestCaseExecutionResponseMessageStruct struct {
	TestCaseExecutionBasicInformation *fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage
	TestCaseExecutionDetails          *[]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
	TestInstructionExecutionsMap      *map[string]*workObjectForTestInstructionExecutionsMessageStruct // map[TestInstructionExecutionKey]*workObjectForTestInstructionExecutionsMessageStruct
}

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadFullTestCasesExecutionInformation(
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
	uniqueCountersForTableTestCaseExecutionQueue, err = fenixGuiExecutionServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestCaseExecutionQueue")

	// If there are no TestCases under onQueue then ignore this part
	if uniqueCountersForTableTestCaseExecutionQueue != nil {
		// Load TestCaseExecutions from table 'TestCaseExecutionQueue'
		_, err = fenixGuiExecutionServerObject.loadTestCasesExecutionsFromOnExecutionQueue(
			txn,
			uniqueCountersForTableTestCaseExecutionQueue,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Convert 'TestCaseExecutionKeys' into slice with 'UniqueCounter' for table 'TestCasesUnderExecution'
	var uniqueCountersForTableTTestCasesUnderExecution []int
	uniqueCountersForTableTTestCasesUnderExecution, err = fenixGuiExecutionServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestCasesUnderExecution")

	// If there are no TestCases under Execution then ignore this part
	if uniqueCountersForTableTTestCasesUnderExecution != nil {

		// Load TestCaseExecutions from table 'TestCasesUnderExecution'
		_, err = fenixGuiExecutionServerObject.loadTestCasesExecutionsFromUnderExecutions(
			txn,
			uniqueCountersForTableTTestCasesUnderExecution,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Convert 'TestInstructionExecutionKeys' into slice with 'UniqueCounter' for table 'TestInstructionExecutionQueue'
	var uniqueCountersForTableTestInstructionExecutionQueue []int
	uniqueCountersForTableTestInstructionExecutionQueue, err = fenixGuiExecutionServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestInstructionExecutionQueue")

	// Only process when there still are TestInstructionExecution on the ExecutionQueue
	if len(uniqueCountersForTableTestInstructionExecutionQueue) > 0 {

		// Load TestInstructionExecutions from table 'TestInstructionExecutionQueue'
		_, err = fenixGuiExecutionServerObject.loadTestInstructionsExecutionsFromOnExecutionQueue(
			txn,
			uniqueCountersForTableTestInstructionExecutionQueue,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Convert 'TestInstructionExecutionKeys' into slice with 'UniqueCounter' for table 'TestInstructionsUnderExecution'
	var uniqueCountersForTableTestInstructionsUnderExecution []int
	uniqueCountersForTableTestInstructionsUnderExecution, err = fenixGuiExecutionServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestInstructionsUnderExecution")

	// Only process when there still are TestInstructionExecution on the ExecutionQueue
	if len(uniqueCountersForTableTestInstructionsUnderExecution) > 0 {

		// Load TestInstructionExecutions from table 'TestInstructionsUnderExecution'
		_, err = fenixGuiExecutionServerObject.loadTestInstructionsExecutionsUnderExecution(
			txn,
			uniqueCountersForTableTestInstructionsUnderExecution,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Load TestCaseExecution-logs
	//var logPostAndValuesMapPtr *map[string]*[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
	err = fenixGuiExecutionServerObject.loadTestCaseExecutionLogs(
		txn,
		testCaseExecutionKeys,
		&tempTestCaseExecutionResponseMessagesMap)

	if err != nil {

		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":    "1b4ce01f-b7bc-40bf-a87a-99fc2f153543",
			"Error": err,
		}).Error("Something went wrong when 'Loading TestCaseExecution-logs'")

		return nil, err
	}

	// Load TestInstructionExecution-RunTime Updated Attributes
	err = fenixGuiExecutionServerObject.loadRunTimeUpdatedAttribute(
		txn,
		testCaseExecutionKeys,
		&tempTestCaseExecutionResponseMessagesMap)

	if err != nil {

		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":    "06f74681-095b-49ec-938e-2cde1bd15e18",
			"Error": err,
		}).Error("Something went wrong when 'Loading TestInstructionExecution-RunTime Updated Attributes'")

		return nil, err
	}

	// Convert 'tempTestCaseExecutionResponseMessagesMap' into gRPC-response object
	err = fenixGuiExecutionServerObject.convertTestCaseExecutionResponseMessagesMapIntoGrpcResponse(
		&tempTestCaseExecutionResponseMessagesMap,
		&testCaseExecutionResponseMessages)

	if err != nil {
		return nil, err
	}

	return testCaseExecutionResponseMessages, err
}

// Convert 'TestCaseExecutionKeys' (TestCaseExecutionUuid + TestCaseExecutionVersion) into a slice with 'UniqueCounter' which are unique number for every DB-row in table
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadUniqueCountersBasedOnTestCaseExecutionKeys(
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
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadTestCasesExecutionsFromOnExecutionQueue(
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
			fenixGuiExecutionServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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

			// Initiate 'TestInstructionExecutionsMap'
			var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct
			tempTestInstructionExecutionsMap = make(map[string]*workObjectForTestInstructionExecutionsMessageStruct)

			// Initiate object to be stored in 'tempTestCaseExecutionResponseMessagesMap'
			tempWorkObjectForTestCaseExecutionResponseMessage = &workObjectForTestCaseExecutionResponseMessageStruct{
				TestCaseExecutionBasicInformation: &testCaseExecutionBasicInformation,
				TestCaseExecutionDetails:          &tempTestCaseExecutionDetails,
				TestInstructionExecutionsMap:      &tempTestInstructionExecutionsMap,
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

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadTestCasesExecutionsFromUnderExecutions(
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
			fenixGuiExecutionServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
		var tempTestCaseExecutionResponseMessagePtr *workObjectForTestCaseExecutionResponseMessageStruct
		var tempTestCaseExecutionResponseMessage workObjectForTestCaseExecutionResponseMessageStruct

		tempTestCaseExecutionResponseMessagePtr, existsInMap = tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

		if existsInMap == false {
			// Initiate object to be stored in 'tempTestCaseExecutionResponseMessagesMap'

			var tempTestCaseExecutionDetailsMessageSlice []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
			tempTestCaseExecutionDetailsMessageSlice = append(tempTestCaseExecutionDetailsMessageSlice, &tempTestCaseExecutionDetailsMessage)

			// Initiate 'TestInstructionExecutionsMap'
			var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct
			tempTestInstructionExecutionsMap = make(map[string]*workObjectForTestInstructionExecutionsMessageStruct)

			tempTestCaseExecutionResponseMessagePtr = &workObjectForTestCaseExecutionResponseMessageStruct{
				TestCaseExecutionBasicInformation: &tempTestCaseExecutionBasicInformationMessage,
				TestCaseExecutionDetails:          &tempTestCaseExecutionDetailsMessageSlice,
				TestInstructionExecutionsMap:      &tempTestInstructionExecutionsMap}

			// Add 'tempTestCaseExecutionResponseMessagePtr' to 'tempTestCaseExecutionResponseMessagesMap'
			tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey] = tempTestCaseExecutionResponseMessagePtr

		} else {

			// Append to existing 'tempTestCaseExecutionResponseMessage'
			tempTestCaseExecutionResponseMessage = *tempTestCaseExecutionResponseMessagePtr
			var tempTestCaseExecutionDetailsPtr *[]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
			var tempTestCaseExecutionDetails []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage

			tempTestCaseExecutionDetailsPtr = tempTestCaseExecutionResponseMessage.TestCaseExecutionDetails
			tempTestCaseExecutionDetails = *tempTestCaseExecutionDetailsPtr

			tempTestCaseExecutionDetails = append(tempTestCaseExecutionDetails, &tempTestCaseExecutionDetailsMessage)

		}

		// Add to number of rows
		numberOfRows = numberOfRows + 1

	}

	return numberOfRows, err
}

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadTestInstructionsExecutionsFromOnExecutionQueue(
	dbTransaction pgx.Tx,
	uniqueCounters []int,
	tempTestCaseExecutionResponseMessagesMapPtr *map[string]*workObjectForTestCaseExecutionResponseMessageStruct) (
	numberOfRows int,
	err error) {

	// Convert reference into variable to use
	tempTestCaseExecutionResponseMessagesMap := *tempTestCaseExecutionResponseMessagesMapPtr

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TIEQ.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructionExecutionQueue\" TIEQ "

	// if uniqueCounters has values then add that as Where-statement
	if uniqueCounters != nil {
		sqlToExecute = sqlToExecute + "WHERE TIEQ.\"UniqueCounter\" IN " +
			fenixGuiExecutionServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
		var tempWorkObjectForTestCaseExecutionResponseMessagePtr *workObjectForTestCaseExecutionResponseMessageStruct
		var tempWorkObjectForTestCaseExecutionResponseMessage workObjectForTestCaseExecutionResponseMessageStruct
		tempWorkObjectForTestCaseExecutionResponseMessagePtr, existsInMap = tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

		if existsInMap == false {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":                             "6ea5ed57-b015-4fca-bee4-26355b2df789",
				"testInstructionExecutionMapKey": testInstructionExecutionMapKey,
			}).Error("Couldn't find 'testCaseExecutionMapKey' in 'tempTestCaseExecutionResponseMessagesMap'")

			return 0, err
		}

		// Extract the 'tempWorkObjectForTestCaseExecutionResponseMessage' to use
		tempWorkObjectForTestCaseExecutionResponseMessage = *tempWorkObjectForTestCaseExecutionResponseMessagePtr

		var tempTestInstructionExecutionsMapPtr *map[string]*workObjectForTestInstructionExecutionsMessageStruct
		var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct

		tempTestInstructionExecutionsMapPtr = tempWorkObjectForTestCaseExecutionResponseMessage.TestInstructionExecutionsMap
		tempTestInstructionExecutionsMap = *tempTestInstructionExecutionsMapPtr

		// Initiate object to be stored in 'TestInstructionExecutionsMap'
		var tempWorkObjectForTestInstructionExecutionsMessage workObjectForTestInstructionExecutionsMessageStruct
		_, existsInMap = tempTestInstructionExecutionsMap[testInstructionExecutionMapKey]

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

			tempWorkObjectForTestInstructionExecutionsMessage = workObjectForTestInstructionExecutionsMessageStruct{
				TestInstructionExecutionBasicInformation: &testInstructionExecutionBasicInformation,
				TestInstructionExecutionsInformation:     &tempTestInstructionExecutions,
				ExecutionLogPostsAndValues:               nil,
				RunTimeUpdatedAttributes:                 nil,
			}

			// Add 'tempWorkObjectForTestInstructionExecutionsMessage' to 'TestInstructionExecutionsMap'
			tempTestInstructionExecutionsMap[testInstructionExecutionMapKey] =
				&tempWorkObjectForTestInstructionExecutionsMessage
		} else {

			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadTestInstructionsExecutionsUnderExecution(
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
			fenixGuiExecutionServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
		var tempWorkObjectForTestCaseExecutionResponseMessagePtr *workObjectForTestCaseExecutionResponseMessageStruct
		var tempWorkObjectForTestCaseExecutionResponseMessage workObjectForTestCaseExecutionResponseMessageStruct
		tempWorkObjectForTestCaseExecutionResponseMessagePtr, existsInMap =
			tempTestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

		if existsInMap == false {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":                             "55ff90f3-6ac2-4c8a-ae34-ad008ccb02a8",
				"testInstructionExecutionMapKey": testInstructionExecutionMapKey,
				"sqlToExecute":                   sqlToExecute,
			}).Error("Couldn't find 'testCaseExecutionMapKey' in 'tempTestCaseExecutionResponseMessagesMap'")

			return 0, errors.New("couldn't find 'testCaseExecutionMapKey' in 'tempTestCaseExecutionResponseMessagesMap")
		}

		// Create 'testInstructionExecutionMapKey'
		testInstructionExecutionMapKey =
			tempTestInstructionExecutionBasicInformation.TestInstructionExecutionUuid +
				strconv.FormatUint(uint64(tempTestInstructionExecutionBasicInformation.TestInstructionExecutionVersion), 10)

		// Extract map with TestInstructionExecution data for the TestCaseExecution
		tempWorkObjectForTestCaseExecutionResponseMessage = *tempWorkObjectForTestCaseExecutionResponseMessagePtr

		var tempTestInstructionExecutionsMapPtr *map[string]*workObjectForTestInstructionExecutionsMessageStruct
		var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct

		tempTestInstructionExecutionsMapPtr = tempWorkObjectForTestCaseExecutionResponseMessage.TestInstructionExecutionsMap
		tempTestInstructionExecutionsMap = *tempTestInstructionExecutionsMapPtr

		// Initiate object to be stored in 'TestInstructionExecutionsMap'
		var tempWorkObjectForTestInstructionExecutionsMessagePtr *workObjectForTestInstructionExecutionsMessageStruct
		var tempWorkObjectForTestInstructionExecutionsMessage workObjectForTestInstructionExecutionsMessageStruct
		tempWorkObjectForTestInstructionExecutionsMessagePtr, existsInMap = tempTestInstructionExecutionsMap[testInstructionExecutionMapKey]

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
			var tempTestInstructionExecution workObjectForTestInstructionExecutionsMessageStruct
			tempTestInstructionExecution = workObjectForTestInstructionExecutionsMessageStruct{
				TestInstructionExecutionBasicInformation: &tempTestInstructionExecutionBasicInformation,
				TestInstructionExecutionsInformation:     &tempTestInstructionExecutions,
				ExecutionLogPostsAndValues:               nil,
				RunTimeUpdatedAttributes:                 nil,
			}

			// Add back to Map
			tempTestInstructionExecutionsMap[testInstructionExecutionMapKey] = &tempTestInstructionExecution
		} else {

			// Extract slice with existing TestInstructionExecutions
			tempWorkObjectForTestInstructionExecutionsMessage = *tempWorkObjectForTestInstructionExecutionsMessagePtr
			var tempTestInstructionExecutionsInformationPtr *[]*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage
			var tempTestInstructionExecutionsInformation []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage

			tempTestInstructionExecutionsInformationPtr = tempWorkObjectForTestInstructionExecutionsMessage.TestInstructionExecutionsInformation
			tempTestInstructionExecutionsInformation = *tempTestInstructionExecutionsInformationPtr

			// Append to existing data
			tempTestInstructionExecutionsInformation = append(tempTestInstructionExecutionsInformation,
				&tempTestInstructionExecutionsInformationMessage)

		}

		// Add to number of rows
		numberOfRows = numberOfRows + 1

	}

	return numberOfRows, err
}

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadTestCaseExecutionLogs(
	dbTransaction pgx.Tx,
	testCaseExecutionKeys []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionKeyMessage,
	tempTestCaseExecutionResponseMessagesMapPtr *map[string]*workObjectForTestCaseExecutionResponseMessageStruct) (
	err error) {

	// Convert from Ptr to Map
	var tempTestCaseExecutionResponseMessagesMap map[string]*workObjectForTestCaseExecutionResponseMessageStruct // map[TestCaseExecutionKey]*[]*workObjectForTestCaseExecutionResponseMessageStruct.LogPostAndValuesMessage
	tempTestCaseExecutionResponseMessagesMap = make(map[string]*workObjectForTestCaseExecutionResponseMessageStruct)

	tempTestCaseExecutionResponseMessagesMap = *tempTestCaseExecutionResponseMessagesMapPtr

	var logPostAndValuesMap map[string]*[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage // map[TestInstructionExecutionKey]*[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
	logPostAndValuesMap = make(map[string]*[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage)

	var existInMap bool

	// Generate slice with TestCaseExecutions to get logs for
	var testCaseExecutionMapKeys []string
	var testCaseExecutionMapKey string

	for _, testCaseExecutionUuid := range testCaseExecutionKeys {

		testCaseExecutionMapKey = testCaseExecutionUuid.GetTestCaseExecutionUuid() +
			strconv.FormatUint(uint64(testCaseExecutionUuid.GetTestCaseExecutionVersion()), 10)

		// Add TestCaseExecutionUuid to the slice for the SQL
		testCaseExecutionMapKeys = append(testCaseExecutionMapKeys, testCaseExecutionMapKey)

	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT ELP.*, CONCAT(ELP.\"TestCaseExecutionUuid\", ELP.\"TestCaseExecutionVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"ExecutionLogPosts\" ELP "
	sqlToExecute = sqlToExecute + "WHERE  CONCAT(ELP.\"TestCaseExecutionUuid\", " +
		"ELP.\"TestCaseExecutionVersion\") IN " +
		fenixGuiExecutionServerObject.generateSQLINArray(testCaseExecutionMapKeys)
	sqlToExecute = sqlToExecute + "ORDER BY ELP.\"LogPostTimeStamp\",  ELP.\"TestInstructionExecutionUuid\" "
	sqlToExecute = sqlToExecute + "; "

	if common_config.LogAllSQLs == true {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "9a0e1c27-c386-46f7-a05d-c36780ae1953",
			"sqlToExecute": sqlToExecute,
		}).Info("SQL to be executed")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "e174dc5f-bc65-4ead-b9f0-530e1631e565",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	var (
		tempDomainUuid                           string
		tempTestCaseExecutionUuid                string
		tempTestCaseExecutionVerion              int
		tempTestInstructionExecutionStatus       int
		tempLogPostUuid                          string
		tempLogPostTimeStamp                     *time.Time
		tempFoundVsExpectedValuesAsJsonbAsString string
		tempTestCaseExecutionMapKey              string
		numberOfRows                             int
	)

	// FoundVersusExpectedValueStruct within 'LogPostStruct'
	// Holds one variables and its expected value vs found value
	type FoundVersusExpectedValueStruct struct {
		FoundValue    string `json:"FoundValue"`
		ExpectedValue string `json:"ExpectedValue"`
	}

	// FoundVersusExpectedValueForVariableStruct within 'LogPostStruct'
	// Holds one variables and its expected value vs found value
	type FoundVersusExpectedValueForVariableStruct struct {
		VariableName              string                         `json:"VariableName"`
		VariableDescription       string                         `json:"VariableDescription"`
		FoundVersusExpectedValues FoundVersusExpectedValueStruct `json:"FoundVersusExpectedValues"`
	}

	// FoundVersusExpectedValueStruct within 'LogPostStruct'
	// Holds one variables and its expected value vs found value
	type FoundVersusExpectedValuesStruct struct {
		FoundVersusExpectedValue []FoundVersusExpectedValueForVariableStruct `json:"FoundVersusExpectedValue"`
	}

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		var tempLogPostAndValues fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage

		err = rows.Scan(
			&tempDomainUuid,
			&tempTestCaseExecutionUuid,
			&tempTestCaseExecutionVerion,
			&tempLogPostAndValues.TestInstructionExecutionUuid,
			&tempLogPostAndValues.TestInstructionExecutionVersion,
			&tempTestInstructionExecutionStatus,
			&tempLogPostUuid,
			&tempLogPostTimeStamp,
			&tempLogPostAndValues.LogPostStatus,
			&tempLogPostAndValues.LogPostText,
			&tempFoundVsExpectedValuesAsJsonbAsString,
			&tempTestCaseExecutionMapKey,
		)

		if err != nil {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":           "9b7a0d31-cc18-4765-aacf-4f0424a9cccf",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// One more row found in database
		numberOfRows = numberOfRows + 1

		// Convert temp-variables into gRPC-variables - LogPostTimeStamp
		if tempLogPostTimeStamp != nil {
			tempLogPostAndValues.LogPostTimeStamp =
				timestamppb.New(*tempLogPostTimeStamp)
		}

		// Clean 'tempFoundVsExpectedValuesAsJsonbAsString'
		//tempFoundVsExpectedValuesAsJsonbAsString = tempFoundVsExpectedValuesAsJsonbAsString[1 : len(tempFoundVsExpectedValuesAsJsonbAsString)-1]

		var tempFoundVsExpectedValue FoundVersusExpectedValuesStruct

		// Check if this an empty json; "{}" or not
		if len(tempFoundVsExpectedValuesAsJsonbAsString) > 4 {
			// There are Found vs Expected values, so add name 'FoundVersusExpectedValue' to the json
			tempFoundVsExpectedValuesAsJsonbAsString = "{\"FoundVersusExpectedValue\":" + tempFoundVsExpectedValuesAsJsonbAsString + "}"

			// Unmarshal (cast) JSON into the struct.
			err = json.Unmarshal([]byte(tempFoundVsExpectedValuesAsJsonbAsString), &tempFoundVsExpectedValue)
			if err != nil {
				fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
					"Id":    "68d920d6-656a-485e-8d31-9ad6fc4d9507",
					"Error": err,
					"tempFoundVsExpectedValuesAsJsonbAsString": tempFoundVsExpectedValuesAsJsonbAsString,
				}).Error("Couldn't unmarshal 'tempFoundVsExpectedValuesAsJsonbAsString' into proto-structure")

				return err
			}
		}

		// Convert local message, from json, into proto-message
		var tempFoundVersusExpectedValues []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage_FoundVersusExpectedValueMessage

		// Loop all Found vs Expected values and convert to proto-message

		for _, extractedFoundVersusExpectedValue := range tempFoundVsExpectedValue.FoundVersusExpectedValue {
			var tempFoundVersusExpectedValue *fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage_FoundVersusExpectedValueMessage

			tempFoundVersusExpectedValue = &fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage_FoundVersusExpectedValueMessage{
				VariableName:        extractedFoundVersusExpectedValue.VariableName,
				VariableDescription: extractedFoundVersusExpectedValue.VariableDescription,
				FoundValue:          extractedFoundVersusExpectedValue.FoundVersusExpectedValues.FoundValue,
				ExpectedValue:       extractedFoundVersusExpectedValue.FoundVersusExpectedValues.ExpectedValue,
			}

			// Add to slice of Expected vs Found slice
			tempFoundVersusExpectedValues = append(tempFoundVersusExpectedValues, tempFoundVersusExpectedValue)

		}

		// Extract the pure Found vs Expected values array and store in the main Log-object
		tempLogPostAndValues.FoundVersusExpectedValue = tempFoundVersusExpectedValues

		// Extract RunTimeUpdatedAttributeSlice from map for certain
		var logPostAndValuesMessageSlicePtr *[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
		var logPostAndValuesMessageSlice []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage

		// Create 'testInstructionExecutionMapKey'
		var testInstructionExecutionMapKey string

		testInstructionExecutionMapKey = tempLogPostAndValues.TestInstructionExecutionUuid +
			strconv.FormatUint(uint64(tempLogPostAndValues.TestInstructionExecutionVersion), 10)

		// Try to extract existing log-post slice for TestInstructionExecution
		logPostAndValuesMessageSlicePtr, existInMap = logPostAndValuesMap[testInstructionExecutionMapKey]

		if existInMap == true {
			// Slice exist in map, so add to existing slice
			logPostAndValuesMessageSlice = *logPostAndValuesMessageSlicePtr

			logPostAndValuesMessageSlice = append(logPostAndValuesMessageSlice, &tempLogPostAndValues)

		} else {
			// First instance of TestInstructionExecution in map so just add to new slice
			logPostAndValuesMessageSlice = append(logPostAndValuesMessageSlice, &tempLogPostAndValues)
		}

		// Store slice back in Map
		logPostAndValuesMap[testInstructionExecutionMapKey] = &logPostAndValuesMessageSlice

	}

	// Check if any logpost were found
	if numberOfRows == 0 {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":                          "1d2c1775-eb46-40e1-83b3-05d33b9830b2",
			"tempTestCaseExecutionMapKey": tempTestCaseExecutionMapKey,
		}).Debug("No Log-post were found in database")

		return nil

	}

	// Store log-posts and values in overall response object

	// Extract TestCaseExecution-object
	var tempTestCaseExecutionPtr *workObjectForTestCaseExecutionResponseMessageStruct
	var tempTestCaseExecution workObjectForTestCaseExecutionResponseMessageStruct
	tempTestCaseExecutionPtr, existInMap = tempTestCaseExecutionResponseMessagesMap[tempTestCaseExecutionMapKey]

	if numberOfRows > 0 && existInMap == false {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":                          "e0cb4f9d-8a3c-46d7-8aaf-9303a8df5031",
			"tempTestCaseExecutionMapKey": tempTestCaseExecutionMapKey,
		}).Error("Should never happen that TestCaseExecution is missing in map, 'tempTestCaseExecutionResponseMessagesMap'")

		err = errors.New("should never happen that TestCaseExecution is missing in map, 'tempTestCaseExecutionResponseMessagesMap'")

		return err
	}

	// Get the object from the Ptr
	tempTestCaseExecution = *tempTestCaseExecutionPtr

	// Get the TestInstructionExecutionMap
	var tempTestInstructionExecutionsMapPtr *map[string]*workObjectForTestInstructionExecutionsMessageStruct
	var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct

	tempTestInstructionExecutionsMapPtr = tempTestCaseExecution.TestInstructionExecutionsMap
	tempTestInstructionExecutionsMap = *tempTestInstructionExecutionsMapPtr

	// Get the TestInstructionExecution-object
	var tempTestInstructionExecutionObjectPtr *workObjectForTestInstructionExecutionsMessageStruct

	// Loop TestInstructionExecutions in LogObject and store log-info and values in main TestInstructionExecution-object
	for testInstructionExecutionMapKey, logPostAndValueSlicePtr := range logPostAndValuesMap {

		// Get logPostAndValueSlice
		var logPostAndValueSlice []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
		logPostAndValueSlice = *logPostAndValueSlicePtr

		// Extract correct TestInstructionExecution-object to store 'logPostAndValueSlice' in
		tempTestInstructionExecutionObjectPtr, existInMap = tempTestInstructionExecutionsMap[testInstructionExecutionMapKey]

		if existInMap == false {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":                             "aa4eb62d-1132-4aae-9653-6ffecf05c045",
				"testInstructionExecutionMapKey": testInstructionExecutionMapKey,
			}).Error("Should never happen that TestInstructionExecution is missing in map, 'tempTestInstructionExecutionsMap'")

			err = errors.New("should never happen that TestInstructionExecution is missing in map, 'tempTestInstructionExecutionsMap'\"")

			return err
		}

		// Create a new object to store Execution LogPosts And Values
		tempTestInstructionExecutionObjectPtr.ExecutionLogPostsAndValues = &logPostAndValueSlice

	}

	return err
}

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadRunTimeUpdatedAttribute(
	dbTransaction pgx.Tx,
	testCaseExecutionKeys []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionKeyMessage,
	tempTestCaseExecutionResponseMessagesMapPtr *map[string]*workObjectForTestCaseExecutionResponseMessageStruct) (
	err error) {

	var runTimeUpdatedAttributesMap map[string]*[]*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage
	runTimeUpdatedAttributesMap = make(map[string]*[]*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage)

	var existInMap bool

	// Get the Map for the TestCaseExecutionResponseMessages
	var tempTestCaseExecutionResponseMessagesMap map[string]*workObjectForTestCaseExecutionResponseMessageStruct
	tempTestCaseExecutionResponseMessagesMap = make(map[string]*workObjectForTestCaseExecutionResponseMessageStruct)

	tempTestCaseExecutionResponseMessagesMap = *tempTestCaseExecutionResponseMessagesMapPtr

	// Generate slice with TestInstructionExecutions to get RunTime Changed Attributes for
	var testInstructionExecutionUuidList []string

	for _, tempTestCaseExecutionPtr := range tempTestCaseExecutionResponseMessagesMap {

		for tempTestInstructionExecutionKey, _ := range *tempTestCaseExecutionPtr.TestInstructionExecutionsMap {

			// Add TestInstructionExecutionKey to the slice for the SQL
			testInstructionExecutionUuidList = append(testInstructionExecutionUuidList, tempTestInstructionExecutionKey)

		}

	}
	/*
		WITH values AS
		(SELECT
		CONCAT(TIAUECH."TestInstructionExecutionUuid", TIAUECH."TestInstructionExecutionVersion", TIAUECH."TestInstructionAttributeUuid") AS key,
			max(TIAUECH."UniqueId_New") AS id

		FROM "FenixExecution"."TestInstructionAttributesUnderExecutionChangeHistory" TIAUECH

		WHERE CONCAT(TIAUECH."TestInstructionExecutionUuid", TIAUECH."TestInstructionExecutionVersion") IN (
			'fc200d71-375e-48c2-b479-b16abff70f4a1',
			'45a90d6e-fc49-4dff-ba75-4f0966f0c2181',
			'b4251d19-2a8f-4e9a-8e9b-153b6b6dc8351',
			'df77db68-7a94-45b1-869a-4961f29d981e1',
			'915bd3d6-c2ce-4f69-814e-7b375cc5b4fb1',
			'6b484bb4-6a6f-48fc-8de2-f9163a3a8ec31'
		)
		GROUP BY
		TIAUECH."TestInstructionExecutionUuid",
			TIAUECH."TestInstructionExecutionVersion",
			TIAUECH."TestInstructionAttributeUuid"

		HAVING  COUNT(*) > 1 )

		SELECT  TIAUECH.* , CONCAT(TIUE."TestCaseExecutionUuid", TIUE."TestCaseExecutionVersion")
		FROM "FenixExecution"."TestInstructionAttributesUnderExecutionChangeHistory" TIAUECH,  "FenixExecution"."TestInstructionsUnderExecution" TIUE
		WHERE TIAUECH."UniqueId_New" IN (SELECT id FROM values) AND
		TIUE."TestInstructionExecutionUuid" = TIAUECH."TestInstructionExecutionUuid" AND  TIUE."TestInstructionInstructionExecutionVersion" = TIAUECH."TestInstructionExecutionVersion";

		sqlToExecute := ""
		sqlToExecute = sqlToExecute + "SELECT DISTINCT ON (TIAUECH.\"TestInstructionExecutionUuid\") TIAUECH.*, " +
			"CONCAT(TIUE.\"TestCaseExecutionUuid\", TIUE.\"TestCaseExecutionVersion\") "
		sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestInstructionAttributesUnderExecutionChangeHistory\" TIAUECH, " +
			" \"FenixExecution\".\"TestInstructionsUnderExecution\" TIUE "
		sqlToExecute = sqlToExecute + "WHERE CONCAT(TIAUECH.\"TestInstructionExecutionUuid\", " +
			"TIAUECH.\"TestInstructionExecutionVersion\") IN " +
			fenixGuiExecutionServerObject.generateSQLINArray(testInstructionExecutionUuidList)
		sqlToExecute = sqlToExecute + " AND "
		sqlToExecute = sqlToExecute + " TIUE.\"TestInstructionExecutionUuid\" = TIAUECH.\"TestInstructionExecutionUuid\" AND "
		sqlToExecute = sqlToExecute + " TIUE.\"TestInstructionInstructionExecutionVersion\" = TIAUECH.\"TestInstructionExecutionVersion\" "
		sqlToExecute = sqlToExecute + "ORDER BY TIAUECH.\"TestInstructionExecutionUuid\", TIAUECH.\"UniqueId\" DESC "
		sqlToExecute = sqlToExecute + "; "
	*/

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "WITH values AS "
	sqlToExecute = sqlToExecute + "(SELECT " +
		"CONCAT(TIAUECH.\"TestInstructionExecutionUuid\", TIAUECH.\"TestInstructionExecutionVersion\", " +
		"TIAUECH.\"TestInstructionAttributeUuid\") AS key, 	MAX(TIAUECH.\"UniqueId_New\") AS id "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestInstructionAttributesUnderExecutionChangeHistory\" TIAUECH "
	sqlToExecute = sqlToExecute + "WHERE CONCAT(TIAUECH.\"TestInstructionExecutionUuid\", TIAUECH.\"TestInstructionExecutionVersion\") IN " +
		fenixGuiExecutionServerObject.generateSQLINArray(testInstructionExecutionUuidList) + " "
	sqlToExecute = sqlToExecute + "GROUP BY "
	sqlToExecute = sqlToExecute + "TIAUECH.\"TestInstructionExecutionUuid\", TIAUECH.\"TestInstructionExecutionVersion\", " +
		"TIAUECH.\"TestInstructionAttributeUuid\" "
	sqlToExecute = sqlToExecute + "HAVING  COUNT(*) > 1 ) "

	sqlToExecute = sqlToExecute + "SELECT  TIAUECH.* , CONCAT(TIUE.\"TestCaseExecutionUuid\", TIUE.\"TestCaseExecutionVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestInstructionAttributesUnderExecutionChangeHistory\" TIAUECH, " +
		"\"FenixExecution\".\"TestInstructionsUnderExecution\" TIUE "
	sqlToExecute = sqlToExecute + "WHERE TIAUECH.\"UniqueId_New\" IN (SELECT id FROM values) AND "
	sqlToExecute = sqlToExecute + "TIUE.\"TestInstructionExecutionUuid\" = TIAUECH.\"TestInstructionExecutionUuid\" AND  " +
		"TIUE.\"TestInstructionInstructionExecutionVersion\" = TIAUECH.\"TestInstructionExecutionVersion\" "
	sqlToExecute = sqlToExecute + "; "

	if common_config.LogAllSQLs == true {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "0160ec19-ba83-49b7-9892-a56d24895898",
			"sqlToExecute": sqlToExecute,
		}).Info("SQL to be executed")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "973506e6-92f2-4141-86ce-326c642ef5c9",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	var (
		tempUpdatedTimeStamp      *time.Time
		tempUniqueId              int
		tempUniqueIdNew           int
		tempTestCaseExecutionUuid string
		tempTestInstructionName   string

		tempTestCaseExecutionMapKey string
		numberOfRows                int
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		var tempRunTimeUpdatedAttribute fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage

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
			&tempTestCaseExecutionUuid,
			&tempTestInstructionName,
			&tempUniqueId,
			&tempUniqueIdNew,
			&tempTestCaseExecutionMapKey,
		)

		if err != nil {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":           "a8f86a8c-820d-4443-886a-11ca7e9f7dfe",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// One more row found in database
		numberOfRows = numberOfRows + 1

		// Convert temp-variables into gRPC-variables
		if tempUpdatedTimeStamp != nil {
			tempRunTimeUpdatedAttribute.UpdateTimeStamp =
				timestamppb.New(*tempUpdatedTimeStamp)
		}

		// Extract RunTimeUpdatedAttributeSlice from map for certain
		var runTimeUpdatedAttributeSlicePtr *[]*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage

		// Create 'testInstructionExecutionMapKey'
		var testInstructionExecutionMapKey string

		testInstructionExecutionMapKey = tempRunTimeUpdatedAttribute.TestInstructionExecutionUuid +
			strconv.FormatUint(uint64(tempRunTimeUpdatedAttribute.TestInstructionExecutionVersion), 10)

		// Try to extract existing RuntTimeAttributes slice for TestInstructionExecution
		runTimeUpdatedAttributeSlicePtr, existInMap = runTimeUpdatedAttributesMap[testInstructionExecutionMapKey]

		if existInMap == true {

			*runTimeUpdatedAttributeSlicePtr = append(*runTimeUpdatedAttributeSlicePtr, &tempRunTimeUpdatedAttribute)

		} else {
			// First instance of TestInstructionExecution in map so just add to new slice
			runTimeUpdatedAttributeSlicePtr = new([]*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage)
			*runTimeUpdatedAttributeSlicePtr = append(*runTimeUpdatedAttributeSlicePtr, &tempRunTimeUpdatedAttribute)
		}

		// Store slice back in Map
		runTimeUpdatedAttributesMap[testInstructionExecutionMapKey] = runTimeUpdatedAttributeSlicePtr

	}

	// Check if any RunTimeUpdated attributes were found
	if numberOfRows == 0 {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":                          "b71e44be-ed23-4992-9f46-e35cd0be0b1d",
			"tempTestCaseExecutionMapKey": tempTestCaseExecutionMapKey,
		}).Debug("No RunTimeUpdated variables were found in database")

		return nil

	}

	// Store log-posts and values in overall response object

	// Extract TestCaseExecution-object
	var tempTestCaseExecutionPtr *workObjectForTestCaseExecutionResponseMessageStruct
	var tempTestCaseExecution workObjectForTestCaseExecutionResponseMessageStruct
	tempTestCaseExecutionPtr, existInMap = tempTestCaseExecutionResponseMessagesMap[tempTestCaseExecutionMapKey]

	if numberOfRows > 0 && existInMap == false {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":                          "81234875-1951-4eb8-b007-cd6de8725b57",
			"tempTestCaseExecutionMapKey": tempTestCaseExecutionMapKey,
		}).Error("Should never happen that TestCaseExecution is missing in map, 'tempTestCaseExecutionResponseMessagesMap'")

		err = errors.New("should never happen that TestCaseExecution is missing in map, 'tempTestCaseExecutionResponseMessagesMap'")

		return err
	}

	// Get the object from the Ptr
	tempTestCaseExecution = *tempTestCaseExecutionPtr

	// Get the TestInstructionExecutionMap
	var tempTestInstructionExecutionsMapPtr *map[string]*workObjectForTestInstructionExecutionsMessageStruct
	var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct

	tempTestInstructionExecutionsMapPtr = tempTestCaseExecution.TestInstructionExecutionsMap
	tempTestInstructionExecutionsMap = *tempTestInstructionExecutionsMapPtr

	// Get the TestInstructionExecution-object
	var tempTestInstructionExecutionObjectPtr *workObjectForTestInstructionExecutionsMessageStruct

	// Loop TestInstructionExecutions in LogObject and store log-info and values in main TestInstructionExecution-object
	for testInstructionExecutionMapKey, runTimeUpdatedAttributesSlicePtr := range runTimeUpdatedAttributesMap {

		// Get runTimeUpdatedAttributeSlice
		var runTimeUpdatedAttributeSlice []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage
		runTimeUpdatedAttributeSlice = *runTimeUpdatedAttributesSlicePtr

		// Extract correct TestInstructionExecution-object to store 'runTimeUpdatedAttributeSlice' in
		tempTestInstructionExecutionObjectPtr, existInMap = tempTestInstructionExecutionsMap[testInstructionExecutionMapKey]

		if existInMap == false {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":                             "07c4adb5-6be3-4e1a-9b14-71f7b940fc37",
				"testInstructionExecutionMapKey": testInstructionExecutionMapKey,
			}).Error("Should never happen that TestInstructionExecution is missing in map, 'tempTestInstructionExecutionsMap'")

			err = errors.New("should never happen that TestInstructionExecution is missing in map, 'tempTestInstructionExecutionsMap'\"")

			return err
		}

		// Create a new object to store Execution LogPosts And Values
		tempTestInstructionExecutionObjectPtr.RunTimeUpdatedAttributes = &runTimeUpdatedAttributeSlice

	}

	return err

}

/*
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) convertTestCaseExecutionResponseMessagesMapIntoGrpcResponse(
	tempTestCaseExecutionResponseMessagesMapReference *map[string]*workObjectForTestCaseExecutionResponseMessageStruct,
	testCaseExecutionResponseMessagesReference *[]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage) (
	err error) {

	var executionLogPostsAndValues *map[string]*[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
	var runTimeUpdatedAttributes *map[string]*[]*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage

	var existInMap bool

	// Loop over TestCaseExecutions in Map
	for _, testCaseExecution := range *tempTestCaseExecutionResponseMessagesMapReference {

		// Create slice for the TestInstructionExecutions within this TestCaseExecution
		var tempTestInstructionExecutions []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage

		// Extract TestInstructionMap
		var tempTestInstructionsExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct
		tempTestInstructionsExecutionsMap = make(map[string]*workObjectForTestInstructionExecutionsMessageStruct)

		tempTestInstructionsExecutionsMap = testCaseExecution.TestInstructionExecutionsMap

		// Loop over TestInstructionExecutions
		var tempLogPostAndValuesMessage []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
		RunTimeUpdatedAttributes[] * RunTimeUpdatedAttributeMessage
		for testInstructionExecutionMapKey, testInstructionExecutionsSlice := range tempTestInstructionsExecutionsMap {

			for _, testInstructionExecution := range testInstructionExecutionsSlice.{

			}

			// Get correct 'ExecutionLogPostsAndValues-object'
			var tempLogPostAndValuesMessagePtr []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
			var tempLogPostAndValuesMessage []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
			tempLogPostAndValuesMessage, existInMap = executionLogPostsAndValues[testInstructionExecution.]

// Get correct 'RunTimeUpdatedAttributes-object'

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

*/

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) convertTestCaseExecutionResponseMessagesMapIntoGrpcResponse(
	tempTestCaseExecutionResponseMessagesMapReference *map[string]*workObjectForTestCaseExecutionResponseMessageStruct,
	testCaseExecutionResponseMessagesReference *[]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage) (
	err error) {

	// Loop over TestCaseExecutions in Map
	for _, testCaseExecution := range *tempTestCaseExecutionResponseMessagesMapReference {

		// Create slice for the TestInstructionExecutions within this TestCaseExecution
		var tempTestInstructionExecutions []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage

		// Extract TestInstructionMap
		var tempTestInstructionExecutionsMapPtr *map[string]*workObjectForTestInstructionExecutionsMessageStruct
		var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct

		tempTestInstructionExecutionsMapPtr = testCaseExecution.TestInstructionExecutionsMap
		tempTestInstructionExecutionsMap = *tempTestInstructionExecutionsMapPtr

		// Loop over TestInstructionExecutions
		for _, testInstructionExecution := range tempTestInstructionExecutionsMap {

			// Create the TestInstructionExecution object to be added
			var tempTestInstructionExecutionsMessage *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage
			tempTestInstructionExecutionsMessage = &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage{
				TestInstructionExecutionBasicInformation: testInstructionExecution.TestInstructionExecutionBasicInformation,
				TestInstructionExecutionsInformation:     *testInstructionExecution.TestInstructionExecutionsInformation,
				ExecutionLogPostsAndValues:               nil,
				RunTimeUpdatedAttributes:                 nil,
			}

			// Check if there is an initiated value for 'Logposts'
			if testInstructionExecution.ExecutionLogPostsAndValues == nil {
				// No initiated value
				var tempExecutionLogPostsAndValues []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
				tempExecutionLogPostsAndValues = []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage{}
				tempTestInstructionExecutionsMessage.ExecutionLogPostsAndValues = tempExecutionLogPostsAndValues
			} else {
				// Initiated value exists
				tempTestInstructionExecutionsMessage.ExecutionLogPostsAndValues = *testInstructionExecution.ExecutionLogPostsAndValues
			}

			// Check if there is an initiated value for 'RunTimeUpdated variables'
			if testInstructionExecution.RunTimeUpdatedAttributes == nil {
				// No initiated value
				var tempRunTimeUpdatedAttributes []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage
				tempRunTimeUpdatedAttributes = []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage{}
				tempTestInstructionExecutionsMessage.RunTimeUpdatedAttributes = tempRunTimeUpdatedAttributes
			} else {
				// Initiated value exists
				tempTestInstructionExecutionsMessage.RunTimeUpdatedAttributes = *testInstructionExecution.RunTimeUpdatedAttributes
			}

			// Append TestInstructionExecution to Slice of all TestInstructionExecutions fur current TestCaseExecution
			tempTestInstructionExecutions = append(tempTestInstructionExecutions, tempTestInstructionExecutionsMessage)

		}
		// Create TestCaseExecution object to be added
		var tempTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
		tempTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage{
			TestCaseExecutionBasicInformation: testCaseExecution.TestCaseExecutionBasicInformation,
			TestCaseExecutionDetails:          *testCaseExecution.TestCaseExecutionDetails,
			TestInstructionExecutions:         tempTestInstructionExecutions,
		}

		// Append TestCaseExecution to Slice of all TestCaseExecutions for current gRPC-response object
		*testCaseExecutionResponseMessagesReference = append(*testCaseExecutionResponseMessagesReference, tempTestCaseExecutionResponseMessage)
	}

	return err
}
