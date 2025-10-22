package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadFullTestSuitesExecutionInformation(
	testSuiteExecutionKeys []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionKeyMessage,
	getGCPAuthenticatedUser string) (
	testSuiteExecutionResponseMessages []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionResponseMessage,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "75545f47-d6f7-4764-8aa7-215307ca1245",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin' in 'loadFullTestSuitesExecutionInformation'")

		return testSuiteExecutionResponseMessages, nil
	}

	// Close db-transaction when leaving this function
	defer txn.Commit(context.Background())

	// Map for keep track of all response messages, but in Map-format instead of slice-format
	var tempTestCaseExecutionResponseMessagesMap map[string]*workObjectForTestCaseExecutionResponseMessageStruct // Key = 'TestCaseExecutionKey'
	tempTestCaseExecutionResponseMessagesMap = make(map[string]*workObjectForTestCaseExecutionResponseMessageStruct)

	// Convert 'TestCaseExecutionKeys' into slice with 'UniqueCounter' for table 'TestCaseExecutionQueue'
	var uniqueCountersForTableTestSuiteExecutionQueue []int
	uniqueCountersForTableTestSuiteExecutionQueue, err = fenixGuiExecutionServerObject.loadUniqueCountersBasedOnTestSuiteExecutionKeys(
		txn, testSuiteExecutionKeys, "TestCaseExecutionQueue")

	// If there are no TestCases under onQueue then ignore this part
	if uniqueCountersForTableTestSuiteExecutionQueue != nil {
		// Load TestCaseExecutions from table 'TestCaseExecutionQueue'
		_, err = fenixGuiExecutionServerObject.loadTestCasesExecutionsFromOnExecutionQueue(
			txn,
			uniqueCountersForTableTestSuiteExecutionQueue,
			&tempTestCaseExecutionResponseMessagesMap)

		if err != nil {
			return nil, err
		}
	}

	// Convert 'TestSuiteExecutionKeys' into slice with 'UniqueCounter' for table 'TestCasesUnderExecution'
	var uniqueCountersForTableTTestCasesUnderExecution []int
	uniqueCountersForTableTTestCasesUnderExecution, err = fenixGuiExecutionServerObject.loadUniqueCountersBasedOnTestSuiteExecutionKeys(
		txn, testSuiteExecutionKeys, "TestCasesUnderExecution")

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

	// Create a slice with all 'TestCaseExecutionKeys' used for the TestSuites
	var testCaseExecutionKeys []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionKeyMessage
	var uuidExample = "587b09bf-0f0d-4f2c-b0e9-edd28c482441"
	var testCaseExecutionUuid string
	var testCaseExecutionVersion int64

	for testCaseExecutionKey, _ := range tempTestCaseExecutionResponseMessagesMap {

		// Extract 'TestCaseExecutionUuid' and 'TestCaseExecutionVersion'
		testCaseExecutionUuid = testCaseExecutionKey[:len(uuidExample)]
		testCaseExecutionVersion, err = strconv.ParseInt(testCaseExecutionKey[len(uuidExample):], 10, 64)

		var tempTestCaseExecutionKey *fenixExecutionServerGuiGrpcApi.TestCaseExecutionKeyMessage
		tempTestCaseExecutionKey = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionKeyMessage{
			TestCaseExecutionUuid:    testCaseExecutionUuid,
			TestCaseExecutionVersion: uint32(testCaseExecutionVersion),
		}

		// Add 'testCaseExecutionKey' slice of keys
		testCaseExecutionKeys = append(testCaseExecutionKeys, tempTestCaseExecutionKey)

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
			"Id":    "3181caef-88ed-4962-a5ff-e47210fd6cde",
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
			"Id":    "fe2505a2-045e-43b2-9292-4a3e73cbe8e3",
			"Error": err,
		}).Error("Something went wrong when 'Loading TestInstructionExecution-RunTime Updated Attributes'")

		return nil, err
	}

	// Convert 'tempTestCaseExecutionResponseMessagesMap' into object used in gRPC-response
	var testCaseExecutionResponseMessages []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
	var testCaseExecutionKeysForEachTestSuiteExecutionKeyMap testCaseExecutionKeysForEachTestSuiteExecutionKeyMapType
	testCaseExecutionKeysForEachTestSuiteExecutionKeyMap, err = fenixGuiExecutionServerObject.
		convertTestCaseExecutionResponseMessagesMapIntoGrpcResponse(
			&tempTestCaseExecutionResponseMessagesMap,
			&testCaseExecutionResponseMessages)

	if err != nil {
		return nil, err
	}

	// Load Domains that User has access to
	var domainAndAuthorizations []DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixGuiExecutionServerObject.PrepareLoadUsersDomains(
		txn,
		getGCPAuthenticatedUser)

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "e1f0bcdf-60d5-4a7c-97da-17225b30e62b",
			"gCPAuthenticatedUser": getGCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		err = errors.New(fmt.Sprintf("User %s doesn't have access to any domains",
			getGCPAuthenticatedUser))

		return nil, err

	}

	// Generate 'testSuiteExecutionKeys' as []string
	var testSuiteExecutionKeyAsString string
	var testSuiteExecutionKeysAsString []string
	for _, testSuiteExecutionKey := range testSuiteExecutionKeys {

		testSuiteExecutionKeyAsString = testSuiteExecutionKey.GetTestSuiteExecutionUuid() +
			strconv.FormatUint(uint64(testSuiteExecutionKey.GetTestSuiteExecutionVersion()), 10)

		// Add to slice of keys
		testSuiteExecutionKeysAsString = append(testSuiteExecutionKeysAsString, testSuiteExecutionKeyAsString)

	}

	// Load TestSuite data to be used in Response
	var rawTestSuiteExecutionsList []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionsListMessage
	rawTestSuiteExecutionsList, _, err = loadRawTestSuiteExecutionsList(
		txn,
		0,
		false,
		0,
		nil,
		nil,
		domainAndAuthorizations,
		false,
		"",
		testSuiteExecutionKeysAsString)

	if err != nil {
		return nil, err
	}

	// Build response object
	var testCaseExecutionKeysForEachTestSuiteExecution []testCaseExecutionKeyType
	var testSuiteExecutionKey testSuiteExecutionKeyType

	// Build TestSuiteExecutionKey
	for _, rawTestSuiteExecution := range rawTestSuiteExecutionsList {

		testSuiteExecutionKey = testSuiteExecutionKeyType(rawTestSuiteExecution.TestSuiteExecutionUuid +
			strconv.FormatUint(uint64(rawTestSuiteExecution.TestSuiteExecutionVersion), 10))

		// Extract the TestCaseExecutionKeys based on the TestSuiteExecutionKey
		testCaseExecutionKeysForEachTestSuiteExecution = testCaseExecutionKeysForEachTestSuiteExecutionKeyMap[testSuiteExecutionKey]

		// Get all TestCaseExecutions for the TestSuiteExecution
		var tempTestCaseExecutionResponseMessages []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
		var testCaseExecutionKeyAsString string
		for _, testCaseExecutionKey := range testCaseExecutionKeysForEachTestSuiteExecution {

			testCaseExecutionKeyAsString = string(testCaseExecutionKey)

			// Extract 'TestCaseExecutionBasicInformation'
			var tempTestCaseExecutionBasicInformation *fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage
			tempTestCaseExecutionBasicInformation = tempTestCaseExecutionResponseMessagesMap[testCaseExecutionKeyAsString].TestCaseExecutionBasicInformation

			// Extract 'TestCaseExecutionDetails'
			var tempTestCaseExecutionDetails []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage
			tempTestCaseExecutionDetails = *tempTestCaseExecutionResponseMessagesMap[testCaseExecutionKeyAsString].TestCaseExecutionDetails

			// Extract 'TestInstructionExecutions'
			var tempTestInstructionExecutions []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage
			var tempTestInstructionExecutionsMapPtr *map[string]*workObjectForTestInstructionExecutionsMessageStruct
			var tempTestInstructionExecutionsMap map[string]*workObjectForTestInstructionExecutionsMessageStruct
			tempTestInstructionExecutionsMapPtr = tempTestCaseExecutionResponseMessagesMap[testCaseExecutionKeyAsString].TestInstructionExecutionsMap
			tempTestInstructionExecutionsMap = *tempTestInstructionExecutionsMapPtr

			// Loop 'tempTestInstructionExecutionsMap' and produce a slice with the data
			for _, tempTestInstructionExecutionsMapEntry := range tempTestInstructionExecutionsMap {

				// 'LogPostAndValuesMessage'
				var tempExecutionLogPostsAndValues []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
				if tempTestInstructionExecutionsMapEntry.ExecutionLogPostsAndValues == nil {
					tempExecutionLogPostsAndValues = []*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage{}
				} else {
					tempExecutionLogPostsAndValues = *tempTestInstructionExecutionsMapEntry.ExecutionLogPostsAndValues
				}

				// Extract 'RunTimeUpdatedAttributes'
				var tempRunTimeUpdatedAttributes []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage
				if tempTestInstructionExecutionsMapEntry.RunTimeUpdatedAttributes == nil {
					tempRunTimeUpdatedAttributes = []*fenixExecutionServerGuiGrpcApi.RunTimeUpdatedAttributeMessage{}
				} else {
					tempRunTimeUpdatedAttributes = *tempTestInstructionExecutionsMapEntry.RunTimeUpdatedAttributes
				}

				// Create 'TestInstructionExecutionsMessage'-structure
				var tempTestInstructionExecution *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage
				tempTestInstructionExecution = &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsMessage{
					TestInstructionExecutionBasicInformation: tempTestInstructionExecutionsMapEntry.TestInstructionExecutionBasicInformation,
					TestInstructionExecutionsInformation:     *tempTestInstructionExecutionsMapEntry.TestInstructionExecutionsInformation,
					ExecutionLogPostsAndValues:               tempExecutionLogPostsAndValues,
					RunTimeUpdatedAttributes:                 tempRunTimeUpdatedAttributes,
				}

				// Add 'tempTestInstructionExecution' to the slice
				tempTestInstructionExecutions = append(tempTestInstructionExecutions, tempTestInstructionExecution)

			}

			// Create the 'tempTestCaseExecutionResponseMessage'
			var tempTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage
			tempTestCaseExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage{
				TestCaseExecutionBasicInformation: tempTestCaseExecutionBasicInformation,
				TestCaseExecutionDetails:          tempTestCaseExecutionDetails,
				TestInstructionExecutions:         tempTestInstructionExecutions,
			}

			// Add 'tempTestCaseExecutionResponseMessage' to slice of messages
			tempTestCaseExecutionResponseMessages = append(tempTestCaseExecutionResponseMessages, tempTestCaseExecutionResponseMessage)

		}

		// Make response object
		var testSuiteExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.TestSuiteExecutionResponseMessage
		testSuiteExecutionResponseMessage = &fenixExecutionServerGuiGrpcApi.TestSuiteExecutionResponseMessage{
			TestSuiteExecutionBasicInformation: &fenixExecutionServerGuiGrpcApi.TestSuiteExecutionBasicInformationMessage{
				DomainUuid:                                   rawTestSuiteExecution.GetDomainUUID(),
				DomainName:                                   rawTestSuiteExecution.GetDomainName(),
				TestSuiteUuid:                                rawTestSuiteExecution.GetTestSuiteUuid(),
				TestSuiteName:                                rawTestSuiteExecution.GetTestSuiteName(),
				TestSuiteVersion:                             uint32(rawTestSuiteExecution.GetTestSuiteVersion()),
				TestSuiteExecutionUuid:                       rawTestSuiteExecution.GetTestSuiteExecutionUuid(),
				TestSuiteExecutionVersion:                    uint32(rawTestSuiteExecution.GetTestSuiteExecutionVersion()),
				UpdatingTestCaseUuid:                         rawTestSuiteExecution.GetUpdatingTestCaseUuid(),
				UpdatingTestCaseName:                         rawTestSuiteExecution.GetUpdatingTestCaseName(),
				UpdatingTestCaseVersion:                      uint32(rawTestSuiteExecution.GetUpdatingTestCaseVersion()),
				UpdatingTestCaseExecutionUuid:                rawTestSuiteExecution.GetUpdatingTestCaseExecutionUuid(),
				UpdatingTestCaseExecutionVersion:             uint32(rawTestSuiteExecution.GetUpdatingTestCaseExecutionVersion()),
				PlacedOnTestExecutionQueueTimeStamp:          rawTestSuiteExecution.GetQueueTimeStamp(),
				TestDataSetUuid:                              rawTestSuiteExecution.GetTestDataSetUuid(),
				ExecutionPriority:                            rawTestSuiteExecution.GetExecutionPriority(),
				ExecutionStatusReportLevel:                   rawTestSuiteExecution.GetExecutionStatusReportLevel(),
				TestSuitePreview:                             rawTestSuiteExecution.GetTestSuitePreview(),
				TestCasesPreviews:                            rawTestSuiteExecution.GetTestCasesPreviews(),
				TestInstructionsExecutionStatusPreviewValues: rawTestSuiteExecution.GetTestInstructionsExecutionStatusPreviewValues(),
			},
			TestSuiteExecutionDetails: []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionDetailsMessage{
				&fenixExecutionServerGuiGrpcApi.TestSuiteExecutionDetailsMessage{
					ExecutionStartTimeStamp:        rawTestSuiteExecution.GetExecutionStartTimeStamp(),
					ExecutionStopTimeStamp:         rawTestSuiteExecution.GetExecutionStopTimeStamp(),
					TestSuiteExecutionStatus:       rawTestSuiteExecution.GetTestSuiteExecutionStatus(),
					ExecutionHasFinished:           rawTestSuiteExecution.GetExecutionHasFinished(),
					ExecutionStatusUpdateTimeStamp: rawTestSuiteExecution.GetExecutionStatusUpdateTimeStamp(),
					UniqueDatabaseRowCounter:       uint64(rawTestSuiteExecution.GetUniqueExecutionCounter()),
				},
			},
			TestCaseExecutions: tempTestCaseExecutionResponseMessages,
		}

		// Add 'testSuiteExecutionResponseMessage' to slice of messages
		testSuiteExecutionResponseMessages = append(testSuiteExecutionResponseMessages, testSuiteExecutionResponseMessage)

	}

	return testSuiteExecutionResponseMessages, err
}

// Convert 'TestSuiteExecutionKeys' (TestSuiteExecutionUuid + TestSuiteExecutionVersion) into a slice with 'UniqueCounter' which are unique number for every DB-row in table
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadUniqueCountersBasedOnTestSuiteExecutionKeys(
	dbTransaction pgx.Tx,
	TestSuiteExecutionKeys []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionKeyMessage,
	databaseTableName string) (
	uniqueCounters []int,
	err error) {

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT \"UniqueCounter\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"" + databaseTableName + "\" "

	// if TestSuiteExecutionKeysList has 'TestSuiteExecutionKeys' then add that as Where-statement
	if TestSuiteExecutionKeys != nil {
		for TestSuiteExecutionKeyCounter, TestSuiteExecutionKey := range TestSuiteExecutionKeys {
			if TestSuiteExecutionKeyCounter == 0 {
				// Add 'Where' for the first TestSuiteExecutionKey, otherwise add an 'ADD'
				sqlToExecute = sqlToExecute + "WHERE "
			} else {
				sqlToExecute = sqlToExecute + "OR "
			}

			sqlToExecute = sqlToExecute + "\"TestSuiteExecutionUuid\" = '" + TestSuiteExecutionKey.TestSuiteExecutionUuid + "' "
			sqlToExecute = sqlToExecute + "AND "
			sqlToExecute = sqlToExecute + "\"TestSuiteExecutionVersion\" = " + strconv.FormatUint(uint64(TestSuiteExecutionKey.TestSuiteExecutionVersion), 10)
			sqlToExecute = sqlToExecute + " "
		}
	}

	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "db7ed25e-e97a-4168-ae86-2b2e529c1f3d",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadUniqueCountersBasedOnTestSuiteExecutionKeys'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "cdc6de82-4376-4670-93e0-0c41c7f2dcf8",
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
				"Id":           "7a505f91-13d1-4ea4-b5c4-30f4f0aceedd",
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

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadTestCaseExecutionLogsBasedOnTestSuiteExecutionUuids(
	dbTransaction pgx.Tx,
	testSuiteExecutionKeys []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionKeyMessage,
	tempTestCaseExecutionResponseMessagesMapPtr *map[string]*workObjectForTestCaseExecutionResponseMessageStruct) (
	err error) {

	// Convert from Ptr to Map
	var tempTestCaseExecutionResponseMessagesMap map[string]*workObjectForTestCaseExecutionResponseMessageStruct // map[TestCaseExecutionKey]*[]*workObjectForTestCaseExecutionResponseMessageStruct.LogPostAndValuesMessage
	tempTestCaseExecutionResponseMessagesMap = make(map[string]*workObjectForTestCaseExecutionResponseMessageStruct)

	tempTestCaseExecutionResponseMessagesMap = *tempTestCaseExecutionResponseMessagesMapPtr

	var logPostAndValuesMap map[string]*[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage // map[TestInstructionExecutionKey]*[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage
	logPostAndValuesMap = make(map[string]*[]*fenixExecutionServerGuiGrpcApi.LogPostAndValuesMessage)

	var existInMap bool

	// Generate slice with TestSuiteExecutions to get logs for
	var testSuiteExecutionMapKeys []string
	var testSuiteExecutionMapKey string

	for _, testSuiteExecutionUuid := range testSuiteExecutionKeys {

		testSuiteExecutionMapKey = testSuiteExecutionUuid.GetTestSuiteExecutionUuid() +
			strconv.FormatUint(uint64(testSuiteExecutionUuid.GetTestSuiteExecutionVersion()), 10)

		// Add TestSuiteExecutionUuid to the slice for the SQL
		testSuiteExecutionMapKeys = append(testSuiteExecutionMapKeys, testSuiteExecutionMapKey)

	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT ELP.*, CONCAT(ELP.\"TestCaseExecutionUuid\", ELP.\"TestCaseExecutionVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"ExecutionLogPosts\" ELP "
	sqlToExecute = sqlToExecute + "WHERE  CONCAT(ELP.\"TestSuiteExecutionUuid\", " +
		"ELP.\"TestSuiteExecutionVersion\") IN " +
		fenixGuiExecutionServerObject.generateSQLINArray(testSuiteExecutionMapKeys)
	sqlToExecute = sqlToExecute + "ORDER BY ELP.\"LogPostTimeStamp\",  ELP.\"TestInstructionExecutionUuid\" "
	sqlToExecute = sqlToExecute + "; "

	if common_config.LogAllSQLs == true {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "1c00d844-588c-4931-9889-88ee66430677",
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
			"Id":           "46a9a39c-6459-402d-857c-ef622133642a",
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
				"Id":           "e73f08f5-659b-40f2-aa97-029d8fe5b013",
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
					"Id":    "c5147a7e-e8a9-4b63-b519-a854db175228",
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
			"Id":                          "986d08cb-7276-4929-a5cd-f0a3d8e53f00",
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
			"Id":                          "79e27f61-6aac-4fa5-9aa8-9cfd7168c7c2",
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
				"Id":                             "ef248361-bfe4-41b3-8bad-ed5d749d0343",
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
