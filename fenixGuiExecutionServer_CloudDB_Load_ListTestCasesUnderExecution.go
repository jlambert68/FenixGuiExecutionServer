package main

import (
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) listTestCasesUnderExecutionLoadFromCloudDB(userID string, domainList []string) (testCaseUnderExecutionMessage []*fenixExecutionServerGuiGrpcApi.TestCaseUnderExecutionMessage, err error) {

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TCUE.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCasesUnderExecution\" TCUE "
	sqlToExecute = sqlToExecute + "WHERE TCUE.\"ExecutionHasFinished\" = false "

	// if domainList has domains then add that as Where-statement
	if domainList != nil {
		sqlToExecute = sqlToExecute + "AND TCUE.\"DomainUuid\" IN " +
			fenixGuiTestCaseBuilderServerObject.generateSQLINArray(domainList)
		sqlToExecute = sqlToExecute + " "
	}

	sqlToExecute = sqlToExecute + "ORDER BY TCUE.\"ExecutionStartTimeStamp\" ASC, TCUE.\"DomainName\" ASC, TCUE.\"TestSuiteName\" ASC, TCUE.\"TestCaseName\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "f5b88d14-dbdb-4cb3-a33a-2b71ab1eeda9",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

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
		testCaseUnderExecution := fenixExecutionServerGuiGrpcApi.TestCaseUnderExecutionMessage{}
		testCaseExecutionBasicInformation := fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage{}
		testCaseExecutionDetails := fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage{}

		err := rows.Scan(
			// TestCaseExecutionBasicInformationMessage
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

			// TestCaseExecutionDetailsMessage
			&tempExecutionStartTimeStamp,
			&tempExecutionStopTimeStamp,
			&tempTestCaseExecutionStatus,
			&testCaseExecutionDetails.ExecutionHasFinished,
			&tempUniqueCounter,
			&tempExecutionStatusUpdateTimeStamp,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "a1833e32-5ad5-468d-8e5f-ca6280d58099",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert temp-variables into gRPC-variables
		testCaseExecutionBasicInformation.PlacedOnTestExecutionQueueTimeStamp = timestamppb.New(tempPlacedOnTestExecutionQueueTimeStamp)
		testCaseExecutionBasicInformation.ExecutionPriority = fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(tempExecutionPriority)

		testCaseExecutionDetails.ExecutionStartTimeStamp = timestamppb.New(tempExecutionStartTimeStamp)
		testCaseExecutionDetails.ExecutionStopTimeStamp = timestamppb.New(tempExecutionStopTimeStamp)
		testCaseExecutionDetails.TestCaseExecutionStatus = fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum(tempTestCaseExecutionStatus)
		testCaseExecutionDetails.ExecutionStatusUpdateTimeStamp = timestamppb.New(tempExecutionStatusUpdateTimeStamp)

		// Build 'TestCaseUnderExecutionMessage'
		testCaseUnderExecution.TestCaseExecutionBasicInformation = &testCaseExecutionBasicInformation
		testCaseUnderExecution.TestCaseExecutionDetails = &testCaseExecutionDetails

		// Add 'TestCaseUnderExecutionMessage' to slice of all 'TestCaseUnderExecutionMessage's
		testCaseUnderExecutionMessage = append(testCaseUnderExecutionMessage, &testCaseUnderExecution)

	}

	return testCaseUnderExecutionMessage, err
}
