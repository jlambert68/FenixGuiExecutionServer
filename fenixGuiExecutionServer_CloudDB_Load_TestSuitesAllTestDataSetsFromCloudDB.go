package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) initiateLoadTestSuitesAllTestDataSetsFromCloudDB(
	testSuiteUuid string) (
	testDataForTesSuiteExecutions []*fenixExecutionServerGuiGrpcApi.TestDataForTestCaseExecutionMessage,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "0ec799f3-19db-4d5e-afa9-248d6e7a53ed",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin' in 'initiateLoadTestSuitesAllTestDataSetsFromCloudDB'")

		errId := "9a20a62a-e717-4b5e-839e-4f53dbaf80d0"

		err = errors.New(fmt.Sprintf("problem to do 'DbPool.Begin' in 'initiateLoadTestSuitesAllTestDataSetsFromCloudDB' [ErrorId: %s]", errId))

		return nil, err
	}

	// Close db-transaction when leaving this function
	defer txn.Commit(context.Background())

	// Load the TestDataSet from the database
	var usersChosenTestDataForTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.UsersChosenTestDataForTestSuiteMessage
	testDataForTesSuiteExecutions, err = fenixGuiTestCaseBuilderServerObject.loadTestSuitesAllTestDataSetsFromCloudDB(
		txn,
		testSuiteUuid)

	if err != nil {
		return nil, err
	}

	// Convert TestData structure to structure used when initiating TestSuiteExecution
	for _, testDataForTestCaseExecution := range usersChosenTestDataForTestSuiteMessage.GetUsersSelectedTestDataPointRow() {
	}

	return testDataForTesSuiteExecutions, err
}

// Get 'raw' TestCase Executions, with or without TestInstructionExecutions
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestSuitesAllTestDataSetsFromCloudDB(
	dbTransaction pgx.Tx,
	testSuiteUuid string) (
	_ *fenixTestCaseBuilderServerGrpcApi.UsersChosenTestDataForTestSuiteMessage,
	err error) {

	/*
		SELECT ts."TestSuiteUuid", ts."TestSuiteTestData"
		FROM "FenixBuilder"."TestSuites" ts
		WHERE ts."TestSuiteUuid" = '975364d5-157b-4926-a2b7-b5260b7826b1'
		ORDER BY ts."TestSuiteVersion" DESC
		LIMIT 1;

	*/

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT ts.\"TestSuiteUuid\", ts.\"TestSuiteTestData\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" ts "
	sqlToExecute = sqlToExecute + "WHERE ts.\"TestSuiteUuid\" = '" + testSuiteUuid + "' "
	sqlToExecute = sqlToExecute + "ORDER BY ts.\"TestSuiteVersion\" DESC"
	sqlToExecute = sqlToExecute + "LIMIT 1 "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "022d79b5-4811-49bd-beb7-2d7b8e2f5205",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadTestSuitesAllTestDataSetsFromCloudDB'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "9c09b9fb-8702-4aac-8436-519218ef5892",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Temp variables to used when extract data from result set
	var tempTestSuiteUuid string
	var tempTestSuiteTestDataAsJson string
	var tempTestSuiteTestDataAsByteArray []byte
	var tempTestSuiteTestDataAsGrpc fenixTestCaseBuilderServerGrpcApi.UsersChosenTestDataForTestSuiteMessage
	var rowFound bool

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempTestSuiteUuid,
			&tempTestSuiteTestDataAsJson,
		)

		rowFound = true

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "3a7c3158-c9a6-4b29-81ba-af7b56ee1b7f",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		if rowFound == false {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "93114a4e-378c-4522-b3b4-f4447bb4ca71",
				"sqlToExecute": sqlToExecute,
			}).Error("Didn't find any row in database, should have found one")

			errId := "cc36d846-5873-4345-a872-99e82637c3a2"

			return nil, errors.New(fmt.Sprintf("Didn't find any row in database, should have found one. [ErrorId: %s]", errId))
		}

		// Convert json-strings into byte-arrays
		tempTestSuiteTestDataAsByteArray = []byte(tempTestSuiteTestDataAsJson)

		// Convert json-byte-arrays into proto-messages
		err = protojson.Unmarshal(tempTestSuiteTestDataAsByteArray, &tempTestSuiteTestDataAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "9546d885-308b-4b72-aa8a-1f19ebe1131e",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestSuiteTestDataAsByteArray' into proto-message")

			return nil, err
		}

		// Max one row can be retrieved
		break

	}

	return &tempTestSuiteTestDataAsGrpc, err
}
