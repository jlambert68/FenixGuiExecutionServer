package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

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

	// Convert 'TestCaseExecutionKeys' into slice with 'UniqueCounter' for table 'TestCaseExecutionQueue'
	var uniqueCountersForTableTestCaseExecutionQueue []int
	uniqueCountersForTableTestCaseExecutionQueue, err = fenixGuiTestCaseBuilderServerObject.loadUniqueCountersBasedOnTestCaseExecutionKeys(
		txn, testCaseExecutionKeys, "TestCaseExecutionQueue")

	var temptestCaseExecutionResponseMessagesMap map[string]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage // map[TestCaseExecutionKey]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
	temptestCaseExecutionResponseMessagesMap = make(map[string]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage)

	// Load TestCaseExecutions from table 'TestCaseExecutionQueue'
	err = fenixGuiTestCaseBuilderServerObject.loadTestCasesExecutionsFromOnExecutionQueue(
		txn,
		uniqueCountersForTableTestCaseExecutionQueue,
		&temptestCaseExecutionResponseMessagesMap)

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
	sqlToExecute = sqlToExecute + "SELECT UniqueCounter "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"" + databaseTableName + "\" "

	// if TestCaseExecutionKeysList has 'TestCaseExecutionKeys' then add that as Where-statement
	if TestCaseExecutionKeys != nil {
		for TestCaseExecutionKeyCounter, TestCaseExecutionKey := range TestCaseExecutionKeys {
			if TestCaseExecutionKeyCounter == 0 {
				// Add 'Where' for the first TestCaseExecutionKey, otherwise add an 'ADD'
				sqlToExecute = sqlToExecute + "WHERE "
			} else {
				sqlToExecute = sqlToExecute + "AND "
			}

			sqlToExecute = sqlToExecute + "\"TestCaseExecutionUuid\" = '" + TestCaseExecutionKey.TestCaseExecutionUuid + "' "
			sqlToExecute = sqlToExecute + "AND "
			sqlToExecute = sqlToExecute + "\"TestCaseExecutionVersion\" = " + strconv.FormatUint(uint64(TestCaseExecutionKey.TestCaseExecutionVersion), 10)
			sqlToExecute = sqlToExecute + " "
		}
	}

	sqlToExecute = sqlToExecute + "; "

	// Query DB
	rows, err := dbTransaction.Query(context.Background(), sqlToExecute)

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
	temptestCaseExecutionResponseMessagesMapReference *map[string]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage) (
	err error) {

	// Convert reference into variable to use
	temptestCaseExecutionResponseMessagesMap := *temptestCaseExecutionResponseMessagesMapReference

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
	rows, err := dbTransaction.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b041cb41-8e3b-4f87-922a-09f23fbb253e",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables used for 'temptestCaseExecutionResponseMessagesMap'
	var testCaseExecutionMapKey string
	var existsInMap bool

	// Variables to used when extract data from result set
	var tempPlacedOnTestExecutionQueueTimeStamp time.Time
	var tempExecutionPriority int
	var tempUniqueCounter int

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
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "030eeab7-5bd0-4013-83f4-3a36d9267c64",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert temp-variables into gRPC-variables
		testCaseExecutionBasicInformation.PlacedOnTestExecutionQueueTimeStamp = timestamppb.New(tempPlacedOnTestExecutionQueueTimeStamp)
		testCaseExecutionBasicInformation.ExecutionPriority = fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(tempExecutionPriority)

		// Create 'testCaseExecutionMapKey'
		testCaseExecutionMapKey = testCaseExecutionBasicInformation.TestCaseExecutionUuid + strconv.FormatUint(uint64(testCaseExecutionBasicInformation.TestCaseExecutionVersion), 10)

		// Check if data exist for testCaseExecutionMapKey
		var tempTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
		tempTestCaseExecutionResponseMessage, existsInMap = temptestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

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

			// Initiate object to be stored in 'temptestCaseExecutionResponseMessagesMap'
			tempTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage{
				TestCaseExecutionBasicInformation: &testCaseExecutionBasicInformation,
				TestCaseExecutionDetails:          nil,
				TestInstructionExecutions:         nil}

			// Add 'tempTestCaseExecutionResponseMessage' to 'temptestCaseExecutionResponseMessagesMap'
			temptestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey] = tempTestCaseExecutionResponseMessage

		} else {
			// Add to existing 'tempTestCaseExecutionResponseMessage'
			tempTestCaseExecutionResponseMessage.TestCaseExecutionBasicInformation = &testCaseExecutionBasicInformation
		}

	}

	return err
}

// The pure TestCaseExecution-information
type testCasesExecutionInformationMessageStruct struct {
	testCaseExecutionBasicInformationMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage
	testCaseExecutionDetailsMessage          *fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestCasesExecutionsFromUnderExecutions(
	dbTransaction pgx.Tx,
	uniqueCounters []int,
	temptestCaseExecutionResponseMessagesMapReference *map[string]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage) (
	err error) {

	// Convert reference into variable to use
	temptestCaseExecutionResponseMessagesMap := *temptestCaseExecutionResponseMessagesMapReference

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TCUE.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCasesUnderExecution\" TCUE "

	// if uniqueCounters has values then add that as Where-statement
	if uniqueCounters != nil {
		sqlToExecute = sqlToExecute + "WHERE TCUQ.\"UniqueCounter\" IN " +
			fenixGuiTestCaseBuilderServerObject.generateSQLINArrayForIntegerSlice(uniqueCounters)

	}
	sqlToExecute = sqlToExecute + "; "

	// Query DB
	rows, err := dbTransaction.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "98b552ed-1031-42da-a5a9-287e542abfb1",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables used for 'temptestCaseExecutionResponseMessagesMap'
	var testCaseExecutionMapKey string
	var existsInMap bool

	// Variables to used when extract data from result set
	var tempPlacedOnTestExecutionQueueTimeStamp time.Time
	var tempExecutionPriority int

	var tempExecutionStartTimeStamp time.Time
	var tempExecutionStopTimeStamp time.Time
	var tempTestCaseExecutionStatus int
	var tempExecutionStatusUpdateTimeStamp time.Time

	var tempUniqueCounter int

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
			&tempUniqueCounter,
			&tempExecutionStatusUpdateTimeStamp,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "61ca3d9d-bc80-4702-873f-48f62bfcadb1",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
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
		var tempTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
		tempTestCaseExecutionResponseMessage, existsInMap = temptestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey]

		if existsInMap == false {
			// Initiate object to be stored in 'temptestCaseExecutionResponseMessagesMap'
			tempTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage{
				TestCaseExecutionBasicInformation: &tempTestCaseExecutionBasicInformationMessage,
				TestCaseExecutionDetails:          &tempTestCaseExecutionDetailsMessage,
				TestInstructionExecutions:         nil}

			// Add 'tempTestCaseExecutionResponseMessage' to 'temptestCaseExecutionResponseMessagesMap'
			temptestCaseExecutionResponseMessagesMap[testCaseExecutionMapKey] = tempTestCaseExecutionResponseMessage

		} else {
			// Add to existing 'tempTestCaseExecutionResponseMessage'
			tempTestCaseExecutionResponseMessage.TestCaseExecutionBasicInformation = &tempTestCaseExecutionBasicInformationMessage
			tempTestCaseExecutionResponseMessage.TestCaseExecutionDetails = &tempTestCaseExecutionDetailsMessage
		}
	}

	return err
}
