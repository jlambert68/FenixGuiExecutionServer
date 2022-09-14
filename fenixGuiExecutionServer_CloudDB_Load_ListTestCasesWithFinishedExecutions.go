package main

import (
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) listTestCasesWithFinishedExecutionsLoadFromCloudDB(userID string) (testCasesWithFinishedExecutions []*fenixExecutionServerGuiGrpcApi.TestCaseWithFinishedExecutionMessage, err error) {

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TCUE.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCaseExecutionQueue\" TCUE "
	sqlToExecute = sqlToExecute + "WHERE TCUE.\"ExecutionHasFinished\" = true"
	sqlToExecute = sqlToExecute + "ORDER BY TCUE.\"ExecutionStartTimeStamp\" ASC, TCUE.\"DomainName\" ASC, TCUE.\"TestSuiteName\" ASC, TCUE.\"TestInstructionName\" ASC "
	sqlToExecute = sqlToExecute + "UNION "
	sqlToExecute = sqlToExecute + "SELECT TCFE.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCasesFinishedExecution\" TCFE "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b5cf1554-e111-4522-b3f3-5de9e6f02367",
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

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		testCaseUnderExecution := fenixExecutionServerGuiGrpcApi.TestCaseWithFinishedExecutionMessage{}
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
			testCaseExecutionDetails.ExecutionHasFinished,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "686b2224-8ff7-470d-992c-0f6438d4dd40",
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

		// Build 'TestCaseWithFinishedExecutionMessage'
		testCaseUnderExecution.TestCaseExecutionBasicInformation = &testCaseExecutionBasicInformation
		testCaseUnderExecution.TestCaseExecutionDetails = &testCaseExecutionDetails

		// Add 'TestCaseWithFinishedExecutionMessage' to slice of all 'TestCaseWithFinishedExecutionMessage's
		testCasesWithFinishedExecutions = append(testCasesWithFinishedExecutions, &testCaseUnderExecution)

	}

	return testCasesWithFinishedExecutions, err
}
