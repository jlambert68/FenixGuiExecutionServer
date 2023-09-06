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

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) listTestCasesOnExecutionQueueLoadFromCloudDB(
	dbTransaction pgx.Tx,
	userID string,
	domainList []string) (
	testCaseExecutionBasicInformationMessage []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage,
	err error) {

	usedDBSchema := "FenixExecution" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TCEQ.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestCaseExecutionQueue\" TCEQ "

	// if domainList has domains then add that as Where-statement
	if domainList != nil {
		sqlToExecute = sqlToExecute + "WHERE TCEQ.\"DomainUuid\" IN " +
			fenixGuiTestCaseBuilderServerObject.generateSQLINArray(domainList)
		sqlToExecute = sqlToExecute + " "
	}

	sqlToExecute = sqlToExecute + "ORDER BY TCEQ.\"QueueTimeStamp\" ASC, TCEQ.\"DomainName\" ASC, TCEQ.\"TestSuiteName\" ASC, TCEQ.\"TestCaseName\" ASC; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "0102874a-1195-433c-b7fb-4788e32ff832",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'listTestCasesOnExecutionQueueLoadFromCloudDB'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "e935a5a5-bed1-445c-8115-1150d59a6301",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

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

			// ReportLevel
			&tempExecutionStatusReportLevelEnum,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "1600d7c4-ce72-430c-9a2f-3b530e5c0f83",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert temp-variables into gRPC-variables
		testCaseExecutionBasicInformation.PlacedOnTestExecutionQueueTimeStamp = timestamppb.New(tempPlacedOnTestExecutionQueueTimeStamp)
		testCaseExecutionBasicInformation.ExecutionPriority = fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum(tempExecutionPriority)

		// Add 'testCaseExecutionBasicInformation' to 'testCaseExecutionBasicInformationMessageFromOnQueue'
		testCaseExecutionBasicInformationMessage = append(testCaseExecutionBasicInformationMessage, &testCaseExecutionBasicInformation)

	}

	return testCaseExecutionBasicInformationMessage, err
}
