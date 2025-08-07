package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) initiateLoadTestSuitesAllTestDataSetsFromCloudDB(
	testSuiteUuid string) (
	testDataFromSimpleTestDataAreaFileMessages []*fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage,
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
	usersChosenTestDataForTestSuiteMessage, err = fenixGuiTestCaseBuilderServerObject.loadTestSuitesAllTestDataSetsFromCloudDB(
		txn,
		testSuiteUuid)

	if err != nil {
		return nil, err
	}

	// Get TestSuites all TestDataSetsValues
	testDataFromSimpleTestDataAreaFileMessages, err = fenixGuiTestCaseBuilderServerObject.loadTestSuitesAllTestDataSetValuesFromCloudDB(
		txn,
		usersChosenTestDataForTestSuiteMessage)

	if err != nil {
		return nil, err
	}

	return testDataFromSimpleTestDataAreaFileMessages, err
}

// Get TestSuites all TestData
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
	sqlToExecute = sqlToExecute + "ORDER BY ts.\"TestSuiteVersion\" DESC "
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

// Get TestSuites all TestDataSetsValues
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadTestSuitesAllTestDataSetValuesFromCloudDB(
	dbTransaction pgx.Tx,
	tempTestSuiteTestDataAsGrpc *fenixTestCaseBuilderServerGrpcApi.UsersChosenTestDataForTestSuiteMessage) (
	testDataFromSimpleTestDataAreaFileMessages []*fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage,
	err error) {

	var sqlWhereClause string

	// Get all TestDataFromOneSimpleTestDataAreaFileFullMessages for each group of (TestDataDomainUuid, TestDataAreaUuid)
	for _, testDataPointNameMapPtr := range tempTestSuiteTestDataAsGrpc.ChosenTestDataPointsPerGroupMap {

		// Get the TestDataRowMap from Ptr
		var testDataPointNameMap fenixTestCaseBuilderServerGrpcApi.TestDataPointNameMapMessage
		testDataPointNameMap = *testDataPointNameMapPtr

		// Loop all TestDataRows
		for _, testDataRowsPerTestDataPoint := range testDataPointNameMap.ChosenTestDataRowsPerTestDataPointMap {

			// Loop all TestDataRows
			for _, testDataRow := range testDataRowsPerTestDataPoint.TestDataRows {

				if sqlWhereClause == "" {
					// First combination
					sqlWhereClause = "(\"TestDataDomainUuid\" = '" + testDataRow.GetTestDataDomainUuid() + "' "
					sqlWhereClause = sqlWhereClause + "AND "
					sqlWhereClause = sqlWhereClause + "\"TestDataAreaUuid\" = '" + testDataRow.GetTestDataAreaUuid() + "') "

				} else {
					// Not the first combination
					sqlWhereClause = sqlWhereClause + "OR "
					sqlWhereClause = sqlWhereClause + "(\"TestDataDomainUuid\" = '" + testDataRow.GetTestDataDomainUuid() + "' "
					sqlWhereClause = sqlWhereClause + "AND "
					sqlWhereClause = sqlWhereClause + "\"TestDataAreaUuid\" = '" + testDataRow.GetTestDataAreaUuid() + "') "

				}
			}
		}
	}

	if len(sqlWhereClause) == 0 {
		// No data to retrieve from database

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                          "72c7b3b8-35d5-4186-92bc-0056bf8e9888",
			"tempTestSuiteTestDataAsGrpc": tempTestSuiteTestDataAsGrpc,
		}).Warning("No data to retrieve from database, shouldn't happen.")

		return nil, nil
	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT \"TestDataDomainUuid\", \"TestDataDomainName\", \"TestDataDomainTemplateName\"," +
		" \"TestDataAreaUuid\", \"TestDataAreaName\", " +
		"\"TestDataFileSha256Hash\", \"ImportantDataInFileSha256Hash\", \"InsertedTimeStamp\", " +
		"\"TestDataFromOneSimpleTestDataAreaFileFullMessage\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestDataFromSimpleTestDataAreaFile\" "
	sqlToExecute = sqlToExecute + "WHERE "
	sqlToExecute = sqlToExecute + sqlWhereClause
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
	var insertedTimeStampAsTimeStamp time.Time
	var tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsString string
	var tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray []byte

	var rowFound bool

	// Extract data from DB result set
	for rows.Next() {

		var oneTestDataFromOneSimpleTestDataAreaFileMessage fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage
		var oneTestDataFromOneSimpleTestDataAreaFileFullMessage fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage

		err = rows.Scan(
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataDomainUuid,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataDomainName,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataDomainTemplateName,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataAreaUuid,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataAreaName,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataFileSha256Hash,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.ImportantDataInFileSha256Hash,
			&insertedTimeStampAsTimeStamp,
			&tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsString,
		)

		rowFound = true

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "a2cd9113-18c9-4f35-b682-24190cb3db94",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		if rowFound == false {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "3859fd8d-047c-43a2-a553-8004d948ab62",
				"sqlToExecute": sqlToExecute,
			}).Error("Didn't find any row in database, should have found one")

			errId := "f46ed657-ece2-41cc-b6e7-c020d9de4327"

			return nil, errors.New(fmt.Sprintf("Didn't find any row in database, should have found one or more. [ErrorId: %s]", errId))
		}

		// Convert json-string into byte-arrays
		tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray = []byte(tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsString)

		// Convert json-byte-array into proto-messages
		err = protojson.Unmarshal(tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray, &oneTestDataFromOneSimpleTestDataAreaFileFullMessage)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "5f68c073-a66c-48a1-b2cf-3cfa4be3b28d",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray' into proto-message")

			return nil, err
		}

		// Add TemplateRepositoryConnectionParameters to list
		testDataFromSimpleTestDataAreaFileMessages = append(testDataFromSimpleTestDataAreaFileMessages, &oneTestDataFromOneSimpleTestDataAreaFileFullMessage)

		// Add TemplateRepositoryConnectionParameters to list
		testDataFromSimpleTestDataAreaFileMessages = append(testDataFromSimpleTestDataAreaFileMessages, &oneTestDataFromOneSimpleTestDataAreaFileFullMessage)

	}

	return testDataFromSimpleTestDataAreaFileMessages, err
}
