package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) listTestCasesUnderExecutionLoadFromCloudDB(dbTransaction pgx.Tx, GCPuserID string, domainList []string) (testCaseUnderExecutionMessage []*fenixExecutionServerGuiGrpcApi.TestCaseUnderExecutionMessage, err error) {

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TCUE.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCasesUnderExecution\" TCUE "
	sqlToExecute = sqlToExecute + "WHERE TCUE.\"ExecutionHasFinished\" = false "

	// if domainList has domains then add that as Where-statement
	if domainList != nil {
		sqlToExecute = sqlToExecute + "AND TCUE.\"DomainUuid\" IN " +
			fenixGuiExecutionServerObject.generateSQLINArray(domainList)
		sqlToExecute = sqlToExecute + " "
	}

	sqlToExecute = sqlToExecute + "ORDER BY TCUE.\"ExecutionStartTimeStamp\" ASC, TCUE.\"DomainName\" ASC, TCUE.\"TestSuiteName\" ASC, TCUE.\"TestCaseName\" ASC; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "1c8148b3-c65d-4b08-95de-37e1a4cb1022",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'listTestCasesUnderExecutionLoadFromCloudDB'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
	var tempExecutionStatusReportLevelEnum int

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

			// ReportLevel
			&tempExecutionStatusReportLevelEnum,
		)

		if err != nil {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
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
