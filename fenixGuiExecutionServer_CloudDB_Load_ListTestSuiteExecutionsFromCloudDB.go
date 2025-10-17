package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) listTestSuiteExecutionsFromCloudDB(
	listTestSuiteExecutionsRequest *fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsRequest) (
	listTestSuiteExecutionsResponse *fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "6ceed49d-5ec4-4fa4-b60b-297cf8dfe093",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin' in 'listTestSuiteExecutionsFromCloudDB'")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		listTestSuiteExecutionsResponse = &fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:  false,
				Comments: "Problem to do 'DbPool.Begin' in 'listTestSuiteExecutionsFromCloudDB'",
				ErrorCodes: []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{
					fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
					common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestSuiteExecutionsList:                     nil,
			LatestUniqueTestSuiteExecutionDatabaseRowId: 0,
			MoreRowsExists:                              false,
		}

		return listTestSuiteExecutionsResponse, nil
	}

	// Close db-transaction when leaving this function
	defer txn.Commit(context.Background())

	// Load Domains that User has access to
	var domainAndAuthorizations []DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixGuiExecutionServerObject.PrepareLoadUsersDomains(
		txn,
		listTestSuiteExecutionsRequest.GetUserAndApplicationRunTimeIdentification().GetGCPAuthenticatedUser())

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "490a30d8-81da-4c35-b7ae-97fc73db9c92",
			"gCPAuthenticatedUser": listTestSuiteExecutionsRequest.GetUserAndApplicationRunTimeIdentification().GCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		// Create response message
		var ackNackResponse *fenixExecutionServerGuiGrpcApi.AckNackResponse
		ackNackResponse = &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack: false,
			Comments: fmt.Sprintf("User %s doesn't have access to any domains",
				listTestSuiteExecutionsRequest.GetUserAndApplicationRunTimeIdentification().GCPAuthenticatedUser),
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
				common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		}

		listTestSuiteExecutionsResponse = &fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse{
			AckNackResponse:                             ackNackResponse,
			TestSuiteExecutionsList:                     nil,
			LatestUniqueTestSuiteExecutionDatabaseRowId: 0,
			MoreRowsExists:                              false,
		}

		return listTestSuiteExecutionsResponse, nil

	}

	// Get 'raw' TestSuite Executions, with or without TestInstructionExecutions
	var rawTestSuiteExecutionsList []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionsListMessage
	var moreRowsExistInDatabase bool

	rawTestSuiteExecutionsList, moreRowsExistInDatabase, err = loadRawTestSuiteExecutionsList(
		txn,
		listTestSuiteExecutionsRequest.GetLatestUniqueTestSuiteExecutionDatabaseRowId(),
		listTestSuiteExecutionsRequest.GetOnlyRetrieveLimitedSizedBatch(),
		listTestSuiteExecutionsRequest.GetBatchSize(),
		listTestSuiteExecutionsRequest.GetTestSuiteExecutionFromTimeStamp(),
		listTestSuiteExecutionsRequest.GetTestSuiteExecutionToTimeStamp(),
		domainAndAuthorizations,
		listTestSuiteExecutionsRequest.GetRetrieveAllExecutionsForSpecificTestSuiteUuid(),
		listTestSuiteExecutionsRequest.GetSpecificTestSuiteUuid(),
		[]string{})

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "b2727c62-2469-401e-ad82-d2623f97e357",
			"Error": err,
		}).Error("Something went wrong when loading raw TestSuiteExecutionsList")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		listTestSuiteExecutionsResponse = &fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:  false,
				Comments: "Something went wrong when loading raw TestSuiteExecutionsList",
				ErrorCodes: []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{
					fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
					common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestSuiteExecutionsList:                     nil,
			LatestUniqueTestSuiteExecutionDatabaseRowId: 0,
			MoreRowsExists:                              false,
		}

		return nil, err
	}

	// Loop TestSuiteExecutions and get all TestCaseExecutions the one's that doesn't have end status for the TestSuiteExecution
	var maxUniqueExecutionCounter int32
	for suiteExecutionIndex, tempRawTestSuiteExecution := range rawTestSuiteExecutionsList {

		// If TestSuiteExecutionStatus is NOT an "End status" then Load all TestCaseExecutions
		if hasTestSuiteAnEndStatus(int32(tempRawTestSuiteExecution.TestSuiteExecutionStatus)) == false {

			// Load ExecutionStatuses for the TestCaseExecutions
			var testCasesExecutionStatusMap map[string]testCaseOrTestSuiteExecutionsForLoadTestCasesExecutionStatusStruct // Key is 'TestCaseExecutionUuid' + 'TestCaseExecutionVersion'
			testCasesExecutionStatusMap, err = fenixGuiExecutionServerObject.loadTestCasesExecutionStatus(
				txn,
				baseSqlWhereOnTestSuiteExecutionUuid,
				[]testCaseOrTestSuiteExecutionsForLoadTestCasesExecutionStatusStruct{
					{
						executionUuid:    tempRawTestSuiteExecution.TestSuiteExecutionUuid,
						executionVersion: tempRawTestSuiteExecution.TestSuiteExecutionVersion,
					},
				})

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":                        "ebb9b34d-d825-49b7-9590-b8d4b01a0610",
					"Error":                     err,
					"tempRawTestSuiteExecution": tempRawTestSuiteExecution,
				}).Error("Something went wrong when loading TestCasesExecutionStatus")

				return listTestSuiteExecutionsResponse, err

			}

			// Load all TestInstructionExecutions for TestCaseExecutions in TestSuiteExecution
			var testInstructionsExecutionStatusPreviewValuesMap map[string][]*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusPreviewValueMessage // Key is 'TestCaseExecutionUuid' + 'TestCaseExecutionVersion'
			testInstructionsExecutionStatusPreviewValuesMap, err = fenixGuiExecutionServerObject.
				loadTestInstructionsExecutionStatusPreviewValues(
					txn,
					baseSqlWhereOnTestSuiteExecutionUuid,
					[]testCaseOrTestSuiteExecutionsForLoadTestCasesExecutionStatusStruct{
						{
							executionUuid:    tempRawTestSuiteExecution.TestSuiteExecutionUuid,
							executionVersion: tempRawTestSuiteExecution.TestSuiteExecutionVersion,
						},
					})

			// Exit when there was a problem reading the database
			if err != nil {
				return nil, err
			}

			// Loop TestCaseExecutions and add TestInstructionsExecutions for the one that doesn't have end status for the TestCaseExecution
			for tempTestCasesExecutionMapKey, tempTestCasesExecutionStatus := range testCasesExecutionStatusMap {

				// If TestCaseExecutionStatus is NOT an "End status" then Load all TestInstructionExecutions
				if hasTestCaseAnEndStatus(tempTestCasesExecutionStatus.testCaseExecutionStatus) == false {

					// Add "TestInstructionsExecutionStatusPreviewValues" to 'Raw TestSuiteExecution'
					var testInstructionsExecutionStatusPreviewValuesMessage *fenixExecutionServerGuiGrpcApi.TestInstructionsExecutionStatusPreviewValuesMessage
					testInstructionsExecutionStatusPreviewValuesMessage = &fenixExecutionServerGuiGrpcApi.
						TestInstructionsExecutionStatusPreviewValuesMessage{
						TestCaseExecutionUuid:                       tempTestCasesExecutionStatus.executionUuid,
						TestCaseExecutionVersion:                    tempTestCasesExecutionStatus.executionVersion,
						TestCaseExecutionStatus:                     fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum(tempTestCasesExecutionStatus.testCaseExecutionStatus),
						TestInstructionExecutionStatusPreviewValues: testInstructionsExecutionStatusPreviewValuesMap[tempTestCasesExecutionMapKey],
					}

					tempRawTestSuiteExecution.TestInstructionsExecutionStatusPreviewValues = testInstructionsExecutionStatusPreviewValuesMessage

					// Exit when there was a problem updating the database
					if err != nil {
						return nil, err
					}

					// Store back the updated TestSuiteExecution in the slice
					rawTestSuiteExecutionsList[suiteExecutionIndex] = tempRawTestSuiteExecution

				} else {

				}

			}
		}

		// Extract 'maxUniqueExecutionCounter'
		maxUniqueExecutionCounter = tempRawTestSuiteExecution.UniqueExecutionCounter
	}

	// Create the 'ListTestSuiteExecutionsResponse'
	listTestSuiteExecutionsResponse = &fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
		TestSuiteExecutionsList:                     rawTestSuiteExecutionsList,
		LatestUniqueTestSuiteExecutionDatabaseRowId: maxUniqueExecutionCounter,
		MoreRowsExists:                              moreRowsExistInDatabase,
	}

	return listTestSuiteExecutionsResponse, nil
}

// The maximum number of TestSuiteExecutions to retrieve in one batch, when asked for
const numberOfTestSuiteExecutionsToRetrieveWhenNotSpecified = 10

// Get 'raw' TestCase Executions, with or without TestInstructionExecutions
func loadRawTestSuiteExecutionsList(
	dbTransaction pgx.Tx,
	latestUniqueTestSuiteExecutionDatabaseRowId int32,
	onlyRetrieveLimitedSizedBatch bool,
	batchSize int32,
	testSuiteExecutionFromTimeStamp *timestamppb.Timestamp,
	testSuiteExecutionToTimeStamp *timestamppb.Timestamp,
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
	retrieveAllExecutionsForSpecificTestSuiteUuid bool,
	specificTestSuiteUuid string,
	specificTestSuiteExecutionsKeys []string) (
	rawTestSuiteExecutionsList []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionsListMessage,
	moreRowsExistInDatabase bool,
	err error) {

	// Generate a Domains list and Calculate the Authorization requirements
	var tempCalculatedDomainAndAuthorizations DomainAndAuthorizationsStruct
	var domainList []string
	for _, domainAndAuthorization := range domainAndAuthorizations {
		// Add to DomainList
		domainList = append(domainList, domainAndAuthorization.DomainUuid)

		// Calculate the Authorization requirements for...
		// TestSuiteAuthorizationLevelOwnedByDomain
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOwnedByThisDomain =
			tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOwnedByThisDomain +
				domainAndAuthorization.CanListAndViewTestCaseOwnedByThisDomain

		// TestCaseAuthorizationLevelHavingTiAndTicWithDomain
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain =
			tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain +
				domainAndAuthorization.CanListAndViewTestCaseHavingTIandTICFromThisDomain
	}

	// Convert Values into string for TestCaseAuthorizationLevelOwnedByDomain
	var tempCanListAndViewTestCaseOwnedByThisDomainAsString string
	tempCanListAndViewTestCaseOwnedByThisDomainAsString = strconv.FormatInt(
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain, 10)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "90934050-a972-4b92-9bc9-c99f9811e28b",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, false, err
	}

	// Convert Values into string for TestCaseAuthorizationLevelOwnedByDomain
	var tempCanListAndViewTestCaseHavingTIandTICfromThisDomainAsString string
	tempCanListAndViewTestCaseHavingTIandTICfromThisDomainAsString = strconv.FormatInt(
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain, 10)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "cceef719-eb91-4aa0-a1a6-1b20e334b21d",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, false, err
	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "WITH EC AS ( "
	sqlToExecute = sqlToExecute + "SELECT \"TestSuiteUuid\", COUNT(*) AS \"ExecutionCount\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestSuitesExecutionsForListings\" "
	sqlToExecute = sqlToExecute + "GROUP BY \"TestSuiteUuid\""
	sqlToExecute = sqlToExecute + ") "

	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM ( "

	// Should we retrieve one execution per TestSuiteUuid or should we retrieve all executions for one TestSuiteUuid
	if retrieveAllExecutionsForSpecificTestSuiteUuid == false {
		// Retrieve one execution TestSuiteUuid
		sqlToExecute = sqlToExecute + "SELECT DISTINCT ON (TCEQL.\"TestSuiteUuid\") TCEQL.* , EC.\"ExecutionCount\" "
	} else {
		// Retrieve all executions for one TestSuiteUuid
		sqlToExecute = sqlToExecute + "SELECT TCEQL.* , EC.\"ExecutionCount\" "

	}

	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestSuitesExecutionsForListings\" TCEQL "
	sqlToExecute = sqlToExecute + "JOIN EC ON TCEQL.\"TestSuiteUuid\" = EC.\"TestSuiteUuid\" "
	sqlToExecute = sqlToExecute + "JOIN \"FenixBuilder\".\"TestSuites\" TC " +
		"ON TC.\"TestSuiteUuid\" = TCEQL.\"TestSuiteUuid\" AND " +
		"TC.\"TestSuiteVersion\" = TCEQL.\"TestSuiteVersion\" "

	// if domainList has domains then add that as Where-statement
	if domainList != nil {
		sqlToExecute = sqlToExecute + "WHERE TCEQL.\"DomainUuid\" IN " +
			common_config.GenerateSQLINArray(domainList)
		sqlToExecute = sqlToExecute + " AND "
	} else {

		// Else exit the SQL
		return nil, false, err
	}

	sqlToExecute = sqlToExecute + "(TC.\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\" & " + tempCanListAndViewTestCaseOwnedByThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TC.\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(TC.\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\" & " + tempCanListAndViewTestCaseHavingTIandTICfromThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TC.\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\" "

	// Should we retrieve one execution per TestSuiteUuid or should we retrieve all executions for one TestSuiteUuid
	if retrieveAllExecutionsForSpecificTestSuiteUuid == false {
		// Retrieve one execution TestSuiteUuid

	} else {

		// Retrieve all executions for one TestSuiteUuid
		if len(specificTestSuiteUuid) != 36 {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                    "a0f9deb1-a28a-441b-a900-14513f7ae8b8",
				"specificTestSuiteUuid": specificTestSuiteUuid,
			}).Error("The TestSuiteUuid doesn't seem to be a UUID")

			err = errors.New(fmt.Sprintf("the TestSuiteUuid (%s) doesn't seem to be a UUID ", specificTestSuiteUuid))

			return nil, false, err
		}

		sqlToExecute = sqlToExecute + "AND TC.\"TestSuiteUuid\" = '" + specificTestSuiteUuid + "' "

	}

	// Add filter criteria in SQL: 'latestUniqueTestSuiteExecutionDatabaseRowId'
	if latestUniqueTestSuiteExecutionDatabaseRowId > 0 {
		sqlToExecute = sqlToExecute + fmt.Sprintf(" AND TCEQL.\"UniqueExecutionCounter\" > %d ",
			latestUniqueTestSuiteExecutionDatabaseRowId)
	}

	// TimeStamp is NULL
	//var nullTimeStampV1 time.Time
	var nullTimeStampV2 time.Time
	//nullTimeStampV1 = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	nullTimeStampV2 = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)

	// if

	// Add filter criteria in SQL: 'testSuiteExecutionFromTimeStamp'
	if !(testSuiteExecutionFromTimeStamp.AsTime().UTC().Nanosecond() == nullTimeStampV2.UTC().Nanosecond()) {

		sqlToExecute = sqlToExecute + fmt.Sprintf(" AND TCEQL.\"ExecutionStartTimeStamp\" > '%s' ",
			testSuiteExecutionFromTimeStamp.String())
	}

	// Add filter criteria in SQL: 'testSuiteExecutionToTimeStamp'
	if !(testSuiteExecutionToTimeStamp.AsTime().UTC().Nanosecond() == nullTimeStampV2.UTC().Nanosecond()) {

		sqlToExecute = sqlToExecute + fmt.Sprintf(" AND TCEQL.\"ExecutionStopTimeStamp\" < '%s' ",
			testSuiteExecutionToTimeStamp.String())
	}

	// Check if specific TestCasesExecutions should be fetched from DB
	if specificTestSuiteExecutionsKeys != nil && len(specificTestSuiteExecutionsKeys) > 0 {
		sqlToExecute = sqlToExecute + "AND TCEQL.\"TestSuiteExecutionUuid\" IN " +
			common_config.GenerateSQLINArray(specificTestSuiteExecutionsKeys)
		sqlToExecute = sqlToExecute + " "

	}

	// Add Ordering for inner SQL
	// Should we retrieve one execution per TestSuiteUuid or should we retrieve all executions for one TestSuiteUuid
	if retrieveAllExecutionsForSpecificTestSuiteUuid == false {
		// Retrieve one execution TestSuiteUuid
		sqlToExecute = sqlToExecute + "ORDER BY TCEQL.\"TestSuiteUuid\", TCEQL.\"UniqueExecutionCounter\" DESC "
	} else {
		// Retrieve all executions for one TestSuiteUuid

	}

	sqlToExecute = sqlToExecute + ") sub "

	// Add Ordering for outer SQL
	sqlToExecute = sqlToExecute + "ORDER BY sub.\"UniqueExecutionCounter\" ASC "

	// Add Limit number of rows if requested
	if onlyRetrieveLimitedSizedBatch == true {
		if batchSize < 1 {
			sqlToExecute = sqlToExecute + fmt.Sprintf("LIMIT %d;", numberOfTestSuiteExecutionsToRetrieveWhenNotSpecified+1)
		} else {
			sqlToExecute = sqlToExecute + fmt.Sprintf("LIMIT %d;", batchSize+1)
		}

	} else {
		sqlToExecute = sqlToExecute + ";"
	}

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "8311e299-c7da-4114-ae43-5dd9be72e98e",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadRawTestSuiteExecutionsList'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "3b4cc6c8-f948-4b00-8452-6cad6dedece6",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, false, err
	}

	// Temp variables to used when extract data from result set
	var tempQueueTimeStamp time.Time
	var tempExecutionPriority int
	var tempExecutionStartTimeStamp time.Time
	var tempExecutionStopTimeStamp time.Time
	var tempTestSuiteExecutionStatus int
	var tempExecutionStatusUpdateTimeStamp time.Time
	var tempExecutionStatusReportLevel int
	var tempTestSuitePreviewAsString string
	var tempTestInstructionsExecutionStatusPreviewValuesAsString string
	var tempTestCasesPreviewAsString string

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		var rawTestSuiteExecutionsListItem fenixExecutionServerGuiGrpcApi.TestSuiteExecutionsListMessage

		err := rows.Scan(
			&rawTestSuiteExecutionsListItem.DomainUUID,
			&rawTestSuiteExecutionsListItem.DomainName,
			&rawTestSuiteExecutionsListItem.TestSuiteUuid,
			&rawTestSuiteExecutionsListItem.TestSuiteName,
			&rawTestSuiteExecutionsListItem.TestSuiteVersion,
			&rawTestSuiteExecutionsListItem.TestSuiteExecutionUuid,
			&rawTestSuiteExecutionsListItem.TestSuiteExecutionVersion,
			&rawTestSuiteExecutionsListItem.UpdatingTestCaseUuid,
			&rawTestSuiteExecutionsListItem.UpdatingTestCaseName,
			&rawTestSuiteExecutionsListItem.UpdatingTestCaseVersion,
			&rawTestSuiteExecutionsListItem.UpdatingTestCaseExecutionUuid,
			&rawTestSuiteExecutionsListItem.UpdatingTestCaseExecutionVersion,
			&tempQueueTimeStamp,
			&rawTestSuiteExecutionsListItem.TestDataSetUuid,
			&tempExecutionPriority,
			&tempExecutionStartTimeStamp,
			&tempExecutionStopTimeStamp,
			&tempTestSuiteExecutionStatus,
			&rawTestSuiteExecutionsListItem.ExecutionHasFinished,
			&rawTestSuiteExecutionsListItem.UniqueCounter,
			&tempExecutionStatusUpdateTimeStamp,
			&tempExecutionStatusReportLevel,
			&tempTestSuitePreviewAsString,
			&tempTestInstructionsExecutionStatusPreviewValuesAsString,
			&rawTestSuiteExecutionsListItem.UniqueExecutionCounter,
			&tempTestCasesPreviewAsString,
			&rawTestSuiteExecutionsListItem.NumberOfTestSuiteExecutionForTestSuite,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "a5343556-649a-4639-aa6e-1d7520cf50ba",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, false, err
		}

		// Convert temp-variables into gRPC-variables
		rawTestSuiteExecutionsListItem.QueueTimeStamp = timestamppb.New(tempQueueTimeStamp)
		rawTestSuiteExecutionsListItem.ExecutionPriority = fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(tempExecutionPriority)
		rawTestSuiteExecutionsListItem.ExecutionStartTimeStamp = timestamppb.New(tempExecutionStartTimeStamp)
		rawTestSuiteExecutionsListItem.ExecutionStopTimeStamp = timestamppb.New(tempExecutionStopTimeStamp)
		rawTestSuiteExecutionsListItem.TestSuiteExecutionStatus = fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum(tempTestSuiteExecutionStatus)
		rawTestSuiteExecutionsListItem.ExecutionStatusUpdateTimeStamp = timestamppb.New(tempExecutionStatusUpdateTimeStamp)
		rawTestSuiteExecutionsListItem.ExecutionStatusReportLevel = fenixExecutionServerGuiGrpcApi.ExecutionStatusReportLevelEnum(tempExecutionStatusReportLevel)

		var tempTestSuitePreviewStructureMessage fenixExecutionServerGuiGrpcApi.TestSuitePreviewMessage //TestSuiteExecutionsListMessage
		//var tempTestSuitePreviewAsGrpc fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage
		err = protojson.Unmarshal([]byte(tempTestSuitePreviewAsString), &tempTestSuitePreviewStructureMessage)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                           "bc96e2bb-827f-4bc4-9ab8-153c0b9e297f",
				"Error":                        err,
				"tempTestSuitePreviewAsString": tempTestSuitePreviewAsString,
			}).Error("Something went wrong when converting 'tempTestSuitePreviewAsString' into proto-message")

			// Drop this message and continue with next message
			return nil, false, err
		}

		rawTestSuiteExecutionsListItem.TestSuitePreview = &tempTestSuitePreviewStructureMessage

		if tempTestInstructionsExecutionStatusPreviewValuesAsString != "{}" {
			var tempTestInstructionsExecutionStatusPreviewValuesMessage fenixExecutionServerGuiGrpcApi.
				TestInstructionsExecutionStatusPreviewValuesMessage

			err = protojson.Unmarshal([]byte(tempTestInstructionsExecutionStatusPreviewValuesAsString), &tempTestInstructionsExecutionStatusPreviewValuesMessage)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":    "5d1f620b-8a29-43ef-8b36-c6a5a930d7f9",
					"Error": err,
					"tempTestInstructionsExecutionStatusPreviewValuesAsString": tempTestInstructionsExecutionStatusPreviewValuesAsString,
				}).Error("Something went wrong when converting 'tempTestInstructionsExecutionStatusPreviewValuesAsString' into proto-message")

				// Drop this message and continue with next message
				return nil, false, err
			}

			rawTestSuiteExecutionsListItem.TestInstructionsExecutionStatusPreviewValues = &tempTestInstructionsExecutionStatusPreviewValuesMessage
		}

		if tempTestCasesPreviewAsString != "{}" {
			var tempTestCasesPreviewMessage fenixExecutionServerGuiGrpcApi.
				TestCasePreviews

			err = protojson.Unmarshal([]byte(tempTestCasesPreviewAsString), &tempTestCasesPreviewMessage)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":                           "1c8f48ed-9171-4b5f-b088-6e1439761d86",
					"Error":                        err,
					"tempTestCasesPreviewAsString": tempTestCasesPreviewAsString,
				}).Error("Something went wrong when converting 'tempTestCasesPreviewAsString' into proto-message")

				// Drop this message and continue with next message
				return nil, false, err
			}

			rawTestSuiteExecutionsListItem.TestCasesPreviews = &tempTestCasesPreviewMessage
		}

		// Add 'rawTestSuiteExecutionsListItem' to 'rawTestSuiteExecutionsList'
		rawTestSuiteExecutionsList = append(rawTestSuiteExecutionsList, &rawTestSuiteExecutionsListItem)

	}

	// Check if batch size should be applied when checking if there are more data for client to retrieve
	if onlyRetrieveLimitedSizedBatch == true {
		if batchSize < 1 {

			// Yes, so check if max batch size was achieved
			if len(rawTestSuiteExecutionsList) > numberOfTestSuiteExecutionsToRetrieveWhenNotSpecified {

				// More rows exists
				moreRowsExistInDatabase = true
				rawTestSuiteExecutionsList = rawTestSuiteExecutionsList[:numberOfTestSuiteExecutionsToRetrieveWhenNotSpecified]
			}
		} else {

			// Yes, so check if max batch size was achieved
			if len(rawTestSuiteExecutionsList) > int(batchSize) {

				// More rows exists
				moreRowsExistInDatabase = true
				rawTestSuiteExecutionsList = rawTestSuiteExecutionsList[:batchSize]
			}
		}
	}

	return rawTestSuiteExecutionsList, moreRowsExistInDatabase, err
}

func hasTestSuiteAnEndStatus(testSuiteExecutionStatus int32) (isTestSuiteEndStatus bool) {

	var testSuiteExecutionStatusProto fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum
	testSuiteExecutionStatusProto = fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum(testSuiteExecutionStatus)

	switch testSuiteExecutionStatusProto {

	// Is an End status
	case fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum_TSE_CONTROLLED_INTERRUPTION,
		fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum_TSE_CONTROLLED_INTERRUPTION_CAN_BE_RERUN,
		fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum_TSE_FINISHED_OK,
		fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum_TSE_FINISHED_OK_CAN_BE_RERUN,
		fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum_TSE_FINISHED_NOT_OK,
		fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum_TSE_FINISHED_NOT_OK_CAN_BE_RERUN,
		fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum_TSE_UNEXPECTED_INTERRUPTION,
		fenixExecutionServerGuiGrpcApi.TestSuiteExecutionStatusEnum_TSE_UNEXPECTED_INTERRUPTION_CAN_BE_RERUN:

		isTestSuiteEndStatus = true

	// Is not an End status
	default:
		isTestSuiteEndStatus = false

	}

	return isTestSuiteEndStatus
}
