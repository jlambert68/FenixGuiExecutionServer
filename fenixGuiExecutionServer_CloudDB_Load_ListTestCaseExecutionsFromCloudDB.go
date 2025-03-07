package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"database/sql"
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

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) listTestCaseExecutionsFromCloudDB(
	listTestCaseExecutionsRequest *fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsRequest) (
	listTestCaseExecutionsResponse *fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "7f7cc69e-9f8a-4bfb-957d-29e1e0dbac55",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin' in 'listTestCaseExecutionsFromCloudDB'")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		listTestCaseExecutionsResponse = &fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:  false,
				Comments: "Problem to do 'DbPool.Begin' in 'ListTestCaseExecutions'",
				ErrorCodes: []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{
					fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
					common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCaseExecutionsList:                     nil,
			LatestUniqueTestCaseExecutionDatabaseRowId: 0,
			MoreRowsExists:                             false,
		}

		return listTestCaseExecutionsResponse, nil
	}

	// Close db-transaction when leaving this function
	defer txn.Commit(context.Background())

	// Load Domains that User has access to
	var domainAndAuthorizations []DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixGuiTestCaseBuilderServerObject.PrepareLoadUsersDomains(
		txn,
		listTestCaseExecutionsRequest.GetUserAndApplicationRunTimeIdentification().GetGCPAuthenticatedUser())

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "f2e2a3dd-c985-4b31-91cd-aad81a3414a0",
			"gCPAuthenticatedUser": listTestCaseExecutionsRequest.GetUserAndApplicationRunTimeIdentification().GCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		// Create response message
		var ackNackResponse *fenixExecutionServerGuiGrpcApi.AckNackResponse
		ackNackResponse = &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack: false,
			Comments: fmt.Sprintf("User %s doesn't have access to any domains",
				listTestCaseExecutionsRequest.GetUserAndApplicationRunTimeIdentification().GCPAuthenticatedUser),
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
				common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		}

		listTestCaseExecutionsResponse = &fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse{
			AckNackResponse:                            ackNackResponse,
			TestCaseExecutionsList:                     nil,
			LatestUniqueTestCaseExecutionDatabaseRowId: 0,
			MoreRowsExists:                             false,
		}

		return listTestCaseExecutionsResponse, nil

	}

	// Get 'raw' TestCase Executions, with or without TestInstructionExecutions
	var rawTestCaseExecutionsList []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionsListMessage
	var moreRowsExistInDatabase bool

	rawTestCaseExecutionsList, moreRowsExistInDatabase, err = loadRawTestCaseExecutionsList(
		txn,
		listTestCaseExecutionsRequest.GetLatestUniqueTestCaseExecutionDatabaseRowId(),
		listTestCaseExecutionsRequest.GetOnlyRetrieveLimitedSizedBatch(),
		listTestCaseExecutionsRequest.GetBatchSize(),
		listTestCaseExecutionsRequest.GetTestCaseExecutionFromTimeStamp(),
		listTestCaseExecutionsRequest.GetTestCaseExecutionToTimeStamp(),
		domainAndAuthorizations,
		listTestCaseExecutionsRequest.GetRetrieveAllExecutionsForSpecificTestCaseUuid(),
		listTestCaseExecutionsRequest.GetSpecificTestCaseUuid())

	// Loop TestCaseExecutions and add TestInstructionsExecutions for the one that doesn't have end status for the TestCaseExecution
	var maxUniqueExecutionCounter int32
	for index, tempRawTestCaseExecution := range rawTestCaseExecutionsList {

		// If TestCaseExecutionStatus is NOT an "End status" then Load all TestInstructionExecutions
		if hasTestCaseAnEndStatus(int32(tempRawTestCaseExecution.TestCaseExecutionStatus)) == false {

			var testInstructionsExecutionStatusPreviewValuesMessage *fenixExecutionServerGuiGrpcApi.TestInstructionsExecutionStatusPreviewValuesMessage

			// Load all TestInstructionExecutions for TestCase
			testInstructionsExecutionStatusPreviewValuesMessage, err = fenixGuiTestCaseBuilderServerObject.
				loadTestInstructionsExecutionStatusPreviewValues(
					txn, tempRawTestCaseExecution)

			// Exit when there was a problem reading the database
			if err != nil {
				return nil, err
			}

			// Load TestCaseExecution-status
			var testCaseExecutionStatus int32
			testCaseExecutionStatus, err = fenixGuiTestCaseBuilderServerObject.
				loadTestCaseExecutionStatus(
					txn, tempRawTestCaseExecution)

			// Exit when there was a problem reading the database
			if err != nil {
				return nil, err
			}

			// Add "ExecutionStatusPreviewValues" to 'Raw TestCaseExecution'
			tempRawTestCaseExecution.TestInstructionsExecutionStatusPreviewValues = testInstructionsExecutionStatusPreviewValuesMessage

			// Add TestCaseExecution-status to 'Raw TestCaseExecution'
			tempRawTestCaseExecution.TestCaseExecutionStatus = fenixExecutionServerGuiGrpcApi.
				TestCaseExecutionStatusEnum(testCaseExecutionStatus)

			// Exit when there was a problem updating the database
			if err != nil {
				return nil, err
			}

			// Store back the updated TestCaseExecution in the slice
			rawTestCaseExecutionsList[index] = tempRawTestCaseExecution

		}

		// Extract 'maxUniqueExecutionCounter'
		maxUniqueExecutionCounter = tempRawTestCaseExecution.UniqueExecutionCounter
	}

	// Create the 'ListTestCaseExecutionsResponse'
	listTestCaseExecutionsResponse = &fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
		TestCaseExecutionsList:                     rawTestCaseExecutionsList,
		LatestUniqueTestCaseExecutionDatabaseRowId: maxUniqueExecutionCounter,
		MoreRowsExists:                             moreRowsExistInDatabase,
	}

	return listTestCaseExecutionsResponse, nil
}

// The maximum number of TestCaseExecutions to retrieve in one batch, when asked for
const numberOfTestCaseExecutionsToRetrieveWhenNotSpecified = 10

// Get 'raw' TestCase Executions, with or without TestInstructionExecutions
func loadRawTestCaseExecutionsList(
	dbTransaction pgx.Tx,
	latestUniqueTestCaseExecutionDatabaseRowId int32,
	onlyRetrieveLimitedSizedBatch bool,
	batchSize int32,
	testCaseExecutionFromTimeStamp *timestamppb.Timestamp,
	testCaseExecutionToTimeStamp *timestamppb.Timestamp,
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
	retrieveAllExecutionsForSpecificTestCaseUuid bool,
	specificTestCaseUuid string) (
	rawTestCaseExecutionsList []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionsListMessage,
	moreRowsExistInDatabase bool,
	err error) {

	// Generate a Domains list and Calculate the Authorization requirements
	var tempCalculatedDomainAndAuthorizations DomainAndAuthorizationsStruct
	var domainList []string
	for _, domainAndAuthorization := range domainAndAuthorizations {
		// Add to DomainList
		domainList = append(domainList, domainAndAuthorization.DomainUuid)

		// Calculate the Authorization requirements for...
		// TestCaseAuthorizationLevelOwnedByDomain
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
			"Id":    "223630cb-946b-4c48-9a18-3b01ec98f6b6",
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
			"Id":    "52eff4c9-4b32-43c0-9646-0d1bc13fd39a",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, false, err
	}
	/*
		WITH EC AS (
		    SELECT "TestCaseUuid", COUNT(*) AS "ExecutionCount"
		    FROM "FenixExecution"."TestCasesExecutionsForListings"
		    GROUP BY "TestCaseUuid"
		)
		SELECT *
		FROM (
		    --SELECT DISTINCT ON (TCEQL."TestCaseUuid") TCEQL.*, EC."ExecutionCount"
		    SELECT TCEQL.*, EC."ExecutionCount"
		    FROM "FenixExecution"."TestCasesExecutionsForListings" TCEQL
		    JOIN EC ON TCEQL."TestCaseUuid" = EC."TestCaseUuid"  -- ✅ JOIN EC HERE
		    JOIN "FenixBuilder"."TestCases" TC
		        ON TC."TestCaseUuid" = TCEQL."TestCaseUuid"
		        AND TC."TestCaseVersion" = TCEQL."TestCaseVersion"
		    WHERE TCEQL."DomainUuid" IN ('16458c6c-4f4f-4011-8bd6-34750490c8c1',
		                                 '7edf2269-a8d3-472c-aed6-8cdcc4a8b6ae',
		                                 'e81b9734-5dce-43c9-8d77-3368940cf126')
		    AND (TC."CanListAndViewTestCaseAuthorizationLevelOwnedByDomain" & 11) =
		        TC."CanListAndViewTestCaseAuthorizationLevelOwnedByDomain"
		    AND (TC."CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomai" & 11) =
		        TC."CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomai"
		    AND TC."TestCaseUuid" = '1fe0caf6-0bbf-4470-8aa5-34123896c697'
		    AND TCEQL."UniqueExecutionCounter" > 26
		    -- ORDER BY TCEQL."TestCaseUuid", TCEQL."UniqueExecutionCounter" DESC
		) sub
		ORDER BY sub."UniqueExecutionCounter" ASC
		LIMIT 5;
	*/

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "WITH EC AS ( "
	sqlToExecute = sqlToExecute + "SELECT \"TestCaseUuid\", COUNT(*) AS \"ExecutionCount\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestCasesExecutionsForListings\" "
	sqlToExecute = sqlToExecute + "GROUP BY \"TestCaseUuid\""
	sqlToExecute = sqlToExecute + ") "

	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM ( "

	// Should we retrieve one execution per TestCaseUuid or should we retrieve all executions for one TestCaseUuid
	if retrieveAllExecutionsForSpecificTestCaseUuid == false {
		// Retrieve one execution TestCaseUuid
		sqlToExecute = sqlToExecute + "SELECT DISTINCT ON (TCEQL.\"TestCaseUuid\") TCEQL.* , EC.\"ExecutionCount\" "
	} else {
		// Retrieve all executions for one TestCaseUuid
		sqlToExecute = sqlToExecute + "SELECT TCEQL.* , EC.\"ExecutionCount\" "

	}

	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestCasesExecutionsForListings\" TCEQL "
	sqlToExecute = sqlToExecute + "JOIN EC ON TCEQL.\"TestCaseUuid\" = EC.\"TestCaseUuid\" "
	sqlToExecute = sqlToExecute + "JOIN \"FenixBuilder\".\"TestCases\" TC " +
		"ON TC.\"TestCaseUuid\" = TCEQL.\"TestCaseUuid\" AND " +
		"TC.\"TestCaseVersion\" = TCEQL.\"TestCaseVersion\" "

	// if domainList has domains then add that as Where-statement
	if domainList != nil {
		sqlToExecute = sqlToExecute + "WHERE TCEQL.\"DomainUuid\" IN " +
			common_config.GenerateSQLINArray(domainList)
		sqlToExecute = sqlToExecute + " AND "
	} else {

		// Else exit the SQL
		return nil, false, err
	}

	sqlToExecute = sqlToExecute + "(TC.\"CanListAndViewTestCaseAuthorizationLevelOwnedByDomain\" & " + tempCanListAndViewTestCaseOwnedByThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TC.\"CanListAndViewTestCaseAuthorizationLevelOwnedByDomain\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(TC.\"CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomai\" & " + tempCanListAndViewTestCaseHavingTIandTICfromThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TC.\"CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomai\" "

	// Should we retrieve one execution per TestCaseUuid or should we retrieve all executions for one TestCaseUuid
	if retrieveAllExecutionsForSpecificTestCaseUuid == false {
		// Retrieve one execution TestCaseUuid

	} else {

		// Retrieve all executions for one TestCaseUuid
		if len(specificTestCaseUuid) != 36 {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                   "5c832e3f-f9de-442f-86c6-50652d7af977",
				"specificTestCaseUuid": specificTestCaseUuid,
			}).Error("The TestCaseUuid doesn't seem to be a UUID")

			err = errors.New(fmt.Sprintf("the TestCaseUuid (%s) doesn't seem to be a UUID ", specificTestCaseUuid))

			return nil, false, err
		}

		sqlToExecute = sqlToExecute + "AND TC.\"TestCaseUuid\" = '" + specificTestCaseUuid + "' "

	}

	// Add filter criteria in SQL: 'latestUniqueTestCaseExecutionDatabaseRowId'
	if latestUniqueTestCaseExecutionDatabaseRowId > 0 {
		sqlToExecute = sqlToExecute + fmt.Sprintf(" AND TCEQL.\"UniqueExecutionCounter\" > %d ",
			latestUniqueTestCaseExecutionDatabaseRowId)
	}

	// TimeStamp is NULL
	var nullTimeStamp time.Time
	nullTimeStamp = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	// Add filter criteria in SQL: 'testCaseExecutionFromTimeStamp'
	if !testCaseExecutionFromTimeStamp.AsTime().Equal(nullTimeStamp) {

		sqlToExecute = sqlToExecute + fmt.Sprintf(" AND TCEQL.\"ExecutionStartTimeStamp\" > '%s' ",
			testCaseExecutionFromTimeStamp.String())
	}

	// Add filter criteria in SQL: 'testCaseExecutionToTimeStamp'
	if !testCaseExecutionFromTimeStamp.AsTime().Equal(nullTimeStamp) {

		sqlToExecute = sqlToExecute + fmt.Sprintf(" AND TCEQL.\"ExecutionStopTimeStamp\" < '%s' ",
			testCaseExecutionToTimeStamp.String())
	}

	// Add Ordering for inner SQL
	// Should we retrieve one execution per TestCaseUuid or should we retrieve all executions for one TestCaseUuid
	if retrieveAllExecutionsForSpecificTestCaseUuid == false {
		// Retrieve one execution TestCaseUuid
		sqlToExecute = sqlToExecute + "ORDER BY TCEQL.\"TestCaseUuid\", TCEQL.\"UniqueExecutionCounter\" DESC "
	} else {
		// Retrieve all executions for one TestCaseUuid

	}

	sqlToExecute = sqlToExecute + ") sub "

	// Add Ordering for outer SQL
	sqlToExecute = sqlToExecute + "ORDER BY sub.\"UniqueExecutionCounter\" ASC "

	// Add Limit number of rows if requested
	if onlyRetrieveLimitedSizedBatch == true {
		if batchSize < 1 {
			sqlToExecute = sqlToExecute + fmt.Sprintf("LIMIT %d;", numberOfTestCaseExecutionsToRetrieveWhenNotSpecified+1)
		} else {
			sqlToExecute = sqlToExecute + fmt.Sprintf("LIMIT %d;", batchSize+1)
		}

	} else {
		sqlToExecute = sqlToExecute + ";"
	}

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "a70c2265-b486-4e5d-b991-0f5b38ccb349",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadRawTestCaseExecutionsList'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "cf9f7307-ba77-402c-bff6-253eceee06d6",
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
	var tempTestCaseExecutionStatus int
	var tempExecutionStatusUpdateTimeStamp time.Time
	var tempExecutionStatusReportLevel int
	var tempTestCasePreviewAsString string
	var tempTestInstructionsExecutionStatusPreviewValuesAsString string

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		var rawTestCaseExecutionsListItem fenixExecutionServerGuiGrpcApi.TestCaseExecutionsListMessage

		err := rows.Scan(
			&rawTestCaseExecutionsListItem.DomainUUID,
			&rawTestCaseExecutionsListItem.DomainName,
			&rawTestCaseExecutionsListItem.TestSuiteUuid,
			&rawTestCaseExecutionsListItem.TestSuiteName,
			&rawTestCaseExecutionsListItem.TestSuiteVersion,
			&rawTestCaseExecutionsListItem.TestSuiteExecutionUuid,
			&rawTestCaseExecutionsListItem.TestSuiteExecutionVersion,
			&rawTestCaseExecutionsListItem.TestCaseUuid,
			&rawTestCaseExecutionsListItem.TestCaseName,
			&rawTestCaseExecutionsListItem.TestCaseVersion,
			&rawTestCaseExecutionsListItem.TestCaseExecutionUuid,
			&rawTestCaseExecutionsListItem.TestCaseExecutionVersion,
			&tempQueueTimeStamp,
			&rawTestCaseExecutionsListItem.TestDataSetUuid,
			&tempExecutionPriority,
			&tempExecutionStartTimeStamp,
			&tempExecutionStopTimeStamp,
			&tempTestCaseExecutionStatus,
			&rawTestCaseExecutionsListItem.ExecutionHasFinished,
			&rawTestCaseExecutionsListItem.UniqueCounter,
			&tempExecutionStatusUpdateTimeStamp,
			&tempExecutionStatusReportLevel,
			&tempTestCasePreviewAsString,
			&tempTestInstructionsExecutionStatusPreviewValuesAsString,
			&rawTestCaseExecutionsListItem.UniqueExecutionCounter,
			&rawTestCaseExecutionsListItem.NumberOfTestCaseExecutionForTestCase,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "299bb9a9-01d6-491b-9cd4-a4d5341d55bc",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, false, err
		}

		// Convert temp-variables into gRPC-variables
		rawTestCaseExecutionsListItem.QueueTimeStamp = timestamppb.New(tempQueueTimeStamp)
		rawTestCaseExecutionsListItem.ExecutionPriority = fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(tempExecutionPriority)
		rawTestCaseExecutionsListItem.ExecutionStartTimeStamp = timestamppb.New(tempExecutionStartTimeStamp)
		rawTestCaseExecutionsListItem.ExecutionStopTimeStamp = timestamppb.New(tempExecutionStopTimeStamp)
		rawTestCaseExecutionsListItem.TestCaseExecutionStatus = fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum(tempTestCaseExecutionStatus)
		rawTestCaseExecutionsListItem.ExecutionStatusUpdateTimeStamp = timestamppb.New(tempExecutionStatusUpdateTimeStamp)
		rawTestCaseExecutionsListItem.ExecutionStatusReportLevel = fenixExecutionServerGuiGrpcApi.ExecutionStatusReportLevelEnum(tempExecutionStatusReportLevel)

		var tempTestCasePreviewStructureMessage fenixExecutionServerGuiGrpcApi.TestCaseExecutionsListMessage
		err = protojson.Unmarshal([]byte(tempTestCasePreviewAsString), &tempTestCasePreviewStructureMessage)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                          "929056f6-0cd0-417a-aab2-3e1caab22baf",
				"Error":                       err,
				"tempTestCasePreviewAsString": tempTestCasePreviewAsString,
			}).Error("Something went wrong when converting 'tempTestCasePreviewAsString' into proto-message")

			// Drop this message and continue with next message
			return nil, false, err
		}

		rawTestCaseExecutionsListItem.TestCasePreview = tempTestCasePreviewStructureMessage.TestCasePreview

		if tempTestInstructionsExecutionStatusPreviewValuesAsString != "{}" {
			var tempTestInstructionsExecutionStatusPreviewValuesMessage fenixExecutionServerGuiGrpcApi.
				TestInstructionsExecutionStatusPreviewValuesMessage

			err = protojson.Unmarshal([]byte(tempTestInstructionsExecutionStatusPreviewValuesAsString), &tempTestInstructionsExecutionStatusPreviewValuesMessage)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":    "b366c05c-1d84-42fb-b0a6-5a864950a857",
					"Error": err,
					"tempTestInstructionsExecutionStatusPreviewValuesAsString": tempTestInstructionsExecutionStatusPreviewValuesAsString,
				}).Error("Something went wrong when converting 'tempTestInstructionsExecutionStatusPreviewValuesAsString' into proto-message")

				// Drop this message and continue with next message
				return nil, false, err
			}

			rawTestCaseExecutionsListItem.TestInstructionsExecutionStatusPreviewValues = &tempTestInstructionsExecutionStatusPreviewValuesMessage
		}

		// Add 'rawTestCaseExecutionsListItem' to 'rawTestCaseExecutionsList'
		rawTestCaseExecutionsList = append(rawTestCaseExecutionsList, &rawTestCaseExecutionsListItem)

	}

	// Check if batch size should be applied when checking if there are more data for client to retrieve
	if onlyRetrieveLimitedSizedBatch == true {
		if batchSize < 1 {

			// Yes, so check if max batch size was achieved
			if len(rawTestCaseExecutionsList) > numberOfTestCaseExecutionsToRetrieveWhenNotSpecified {

				// More rows exists
				moreRowsExistInDatabase = true
				rawTestCaseExecutionsList = rawTestCaseExecutionsList[:numberOfTestCaseExecutionsToRetrieveWhenNotSpecified]
			}
		} else {

			// Yes, so check if max batch size was achieved
			if len(rawTestCaseExecutionsList) > int(batchSize) {

				// More rows exists
				moreRowsExistInDatabase = true
				rawTestCaseExecutionsList = rawTestCaseExecutionsList[:batchSize]
			}
		}
	}

	return rawTestCaseExecutionsList, moreRowsExistInDatabase, err
}

func hasTestCaseAnEndStatus(testCaseExecutionStatus int32) (isTestCaseEndStatus bool) {

	var testCaseExecutionStatusProto fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum
	testCaseExecutionStatusProto = fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum(testCaseExecutionStatus)

	switch testCaseExecutionStatusProto {

	// Is an End status
	case fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_CONTROLLED_INTERRUPTION,
		fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_CONTROLLED_INTERRUPTION_CAN_BE_RERUN,
		fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_FINISHED_OK,
		fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_FINISHED_OK_CAN_BE_RERUN,
		fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_FINISHED_NOT_OK,
		fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_FINISHED_NOT_OK_CAN_BE_RERUN,
		fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_UNEXPECTED_INTERRUPTION,
		fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_TCE_UNEXPECTED_INTERRUPTION_CAN_BE_RERUN:

		isTestCaseEndStatus = true

	// Is not an End status
	default:
		isTestCaseEndStatus = false

	}

	return isTestCaseEndStatus
}

// Retrieve "ExecutionStatusPreviewValues" for all TestInstructions for one TestCaseExecution
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestInstructionsExecutionStatusPreviewValues(
	dbTransaction pgx.Tx,
	testCaseExecutionsListMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionsListMessage) (
	testInstructionsExecutionStatusPreviewValuesMessage *fenixExecutionServerGuiGrpcApi.TestInstructionsExecutionStatusPreviewValuesMessage,
	err error) {

	// Load 'ExecutionStatusPreviewValues'

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TIUE.\"TestCaseExecutionUuid\", TIUE.\"TestCaseExecutionVersion\", "
	sqlToExecute = sqlToExecute + "TIUE.\"TestInstructionExecutionUuid\", TIUE.\"TestInstructionInstructionExecutionVersion\", "
	sqlToExecute = sqlToExecute + "TIUE.\"MatureTestInstructionUuid\", TIUE.\"TestInstructionName\", "
	sqlToExecute = sqlToExecute + "TIUE.\"SentTimeStamp\", TIUE.\"TestInstructionExecutionEndTimeStamp\", "
	sqlToExecute = sqlToExecute + "TIUE.\"TestInstructionExecutionStatus\", "
	sqlToExecute = sqlToExecute + "TIUE.\"ExecutionDomainUuid\", TIUE.\"ExecutionDomainName\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestInstructionsUnderExecution\" TIUE "
	sqlToExecute = sqlToExecute + "WHERE "
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"TestCaseExecutionUuid\" = '%s' AND \"TestCaseExecutionVersion\" = %d ",
		testCaseExecutionsListMessage.TestCaseExecutionUuid,
		testCaseExecutionsListMessage.TestCaseExecutionVersion)
	sqlToExecute = sqlToExecute + "ORDER BY TIUE.\"SentTimeStamp\" ASC"
	sqlToExecute = sqlToExecute + ";"

	// Query DB
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "0e402c36-1468-459a-b11d-1c43e6995304",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Number of rows
	var numberOfRowFromDB int32
	numberOfRowFromDB = 0

	var testCasePreviewAndExecutionStatusPreviewValues []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusPreviewValueMessage
	var sentTimeStampAsTimeStamp time.Time
	var testInstructionExecutionEndTimeStampAsTimeStamp time.Time
	var nullableTestInstructionExecutionEndTimeStampAsTimeStamp sql.NullTime

	// Extract data from DB result set
	for rows.Next() {

		var testCasePreviewAndExecutionStatusPreviewValue fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusPreviewValueMessage
		numberOfRowFromDB = numberOfRowFromDB + 1

		err := rows.Scan(
			&testCasePreviewAndExecutionStatusPreviewValue.TestCaseExecutionUuid,
			&testCasePreviewAndExecutionStatusPreviewValue.TestCaseExecutionVersion,
			&testCasePreviewAndExecutionStatusPreviewValue.TestInstructionExecutionUuid,
			&testCasePreviewAndExecutionStatusPreviewValue.TestInstructionInstructionExecutionVersion,
			&testCasePreviewAndExecutionStatusPreviewValue.MatureTestInstructionUuid,
			&testCasePreviewAndExecutionStatusPreviewValue.TestInstructionName,
			&sentTimeStampAsTimeStamp,
			&nullableTestInstructionExecutionEndTimeStampAsTimeStamp,
			&testCasePreviewAndExecutionStatusPreviewValue.TestInstructionExecutionStatus,
			&testCasePreviewAndExecutionStatusPreviewValue.ExecutionDomainUuid,
			&testCasePreviewAndExecutionStatusPreviewValue.ExecutionDomainName,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":                "8ab8a4ba-a743-4705-a0b0-1ebae7f24063",
				"Error":             err,
				"sqlToExecute":      sqlToExecute,
				"numberOfRowFromDB": numberOfRowFromDB,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Check if the timestamp is valid or NULL
		if nullableTestInstructionExecutionEndTimeStampAsTimeStamp.Valid {
			// Timestamp is not NULL
			testInstructionExecutionEndTimeStampAsTimeStamp = nullableTestInstructionExecutionEndTimeStampAsTimeStamp.Time
		} else {
			// TimeStamp is NULL
			testInstructionExecutionEndTimeStampAsTimeStamp = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		}

		// Convert DataTime into gRPC-version
		testCasePreviewAndExecutionStatusPreviewValue.SentTimeStamp = timestamppb.New(sentTimeStampAsTimeStamp)
		testCasePreviewAndExecutionStatusPreviewValue.TestInstructionExecutionEndTimeStamp = timestamppb.
			New(testInstructionExecutionEndTimeStampAsTimeStamp)

		// Add value to slice of values
		testCasePreviewAndExecutionStatusPreviewValues = append(testCasePreviewAndExecutionStatusPreviewValues,
			&testCasePreviewAndExecutionStatusPreviewValue)

	}

	testInstructionsExecutionStatusPreviewValuesMessage = &fenixExecutionServerGuiGrpcApi.
		TestInstructionsExecutionStatusPreviewValuesMessage{
		TestInstructionExecutionStatusPreviewValues: testCasePreviewAndExecutionStatusPreviewValues}

	return testInstructionsExecutionStatusPreviewValuesMessage, err

}

// Retrieve "TestCaseExecutionStatus" for one TestCaseExecution
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestCaseExecutionStatus(
	dbTransaction pgx.Tx,
	rawTestCaseExecution *fenixExecutionServerGuiGrpcApi.TestCaseExecutionsListMessage) (
	testCaseExecutionStatus int32,
	err error) {

	// Load 'ExecutionStatusPreviewValues'

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TCUE.\"TestCaseExecutionStatus\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestCasesUnderExecution\" TCUE "
	sqlToExecute = sqlToExecute + "WHERE "
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"TestCaseExecutionUuid\" = '%s' AND \"TestCaseExecutionVersion\" = %d ",
		rawTestCaseExecution.TestCaseExecutionUuid,
		rawTestCaseExecution.TestCaseExecutionVersion)
	sqlToExecute = sqlToExecute + ";"

	// Query DB
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "fc0b2088-798c-48f1-8dc3-373c0d58c636",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return 0, err
	}

	// Number of rows
	var numberOfRowFromDB int32
	numberOfRowFromDB = 0

	// Extract data from DB result set
	for rows.Next() {

		numberOfRowFromDB = numberOfRowFromDB + 1

		err := rows.Scan(
			&testCaseExecutionStatus,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":                "41ee1994-5a0a-43a6-ac3c-d3eb4ed46137",
				"Error":             err,
				"sqlToExecute":      sqlToExecute,
				"numberOfRowFromDB": numberOfRowFromDB,
			}).Error("Something went wrong when processing result from database")

			return 0, err
		}

	}

	// If number of rows <> 1 then there is a problem
	if numberOfRowFromDB != 1 {

		err = errors.New("number of rows in database response was not exact 1 row")

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                "a0a0b95e-42c6-4b64-a671-668ccdc0de52",
			"Error":             err,
			"sqlToExecute":      sqlToExecute,
			"numberOfRowFromDB": numberOfRowFromDB,
		}).Error("Something went wrong when processing result from database")

		return 0, err
	}

	return testCaseExecutionStatus, err

}
