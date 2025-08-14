package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

// Load 'simple' TestData from Database
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadTestDataFromSimpleTestDataAreaFile(
	testDataAreaUuidSlice []string) (
	testDataFromSimpleTestDataAreaFileMessages []*fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage,
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"Id":                    "69b610e2-6448-4459-b650-c64f5d7ee688",
		"testDataAreaUuidSlice": testDataAreaUuidSlice,
	}).Debug("Entering: loadTestDataFromSimpleTestDataAreaFile")

	defer func() {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id": "46ae2549-fac0-48c2-be94-796b8e5ee687",
		}).Debug("Exiting: loadTestDataFromSimpleTestDataAreaFile")
	}()

	// Generate SQLINArray containing DomainUuids
	var sQLINArray string

	// create the IN-array...('sdada', 'adadadf')
	sQLINArray = fenixGuiExecutionServerObject.generateSQLINArray(testDataAreaUuidSlice)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT \"TestDataDomainUuid\", \"TestDataDomainName\", \"TestDataDomainTemplateName\"," +
		" \"TestDataAreaUuid\", \"TestDataAreaName\", " +
		"\"TestDataFileSha256Hash\", \"ImportantDataInFileSha256Hash\", \"InsertedTimeStamp\", " +
		"\"TestDataFromOneSimpleTestDataAreaFileFullMessage\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestDataFromSimpleTestDataAreaFile\" "
	sqlToExecute = sqlToExecute + "WHERE \"TestDataAreaUuid\" IN " + sQLINArray + ""
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "49e2cdeb-5404-441d-b185-040b36fef61c",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadTestDataFromSimpleTestDataAreaFile'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "9174ec6b-f244-45f3-ac7a-5cd4b6cf577e",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var insertedTimeStampAsTimeStamp time.Time
	var tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsString string
	var tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray []byte

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

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "4cda5dd5-94dd-4216-a56a-a377ad08a5a1",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert json-string into byte-arrays
		tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray = []byte(tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsString)

		// Convert json-byte-array into proto-messages
		err = protojson.Unmarshal(tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray, &oneTestDataFromOneSimpleTestDataAreaFileFullMessage)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "89bfcfac-d89d-40ae-8d47-0data127138",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray' into proto-message")

			return nil, err
		}

		// Add TemplateRepositoryConnectionParameters to list
		testDataFromSimpleTestDataAreaFileMessages = append(testDataFromSimpleTestDataAreaFileMessages, &oneTestDataFromOneSimpleTestDataAreaFileFullMessage)

	}

	return testDataFromSimpleTestDataAreaFileMessages, err
}
