package main

import (
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// ListAllImmatureTestInstructionAttributes - *********************************************************************
// The TestCase Builder asks for all TestInstructions Attributes that the user must add values to in TestCase
func (s *fenixTestCaseBuilderServerGrpcServicesServer) ListAllImmatureTestInstructionAttributes(ctx context.Context, userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage, error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiExecutionServerObject.isClientUsingCorrectTestDataProtoFileVersion("666", userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage{
			TestInstructionAttributesList: nil,
			AckNackResponse:               returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Current user
	userID := userIdentificationMessage.UserId

	// Define variables to store data from DB in
	var testInstructionAttributesList []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage

	// Get users ImmatureTestInstruction-data from CloudDB
	testInstructionAttributesList, err := fenixGuiExecutionServerObject.loadClientsImmatureTestInstructionAttributesFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage{
			TestInstructionAttributesList: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ImmatureTestInstructionAttributes from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiExecutionServerObject.getHighestFenixTestDataProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage{
		TestInstructionAttributesList: testInstructionAttributesList,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiExecutionServerObject.getHighestFenixTestDataProtoFileVersion()),
		},
	}

	return responseMessage, nil
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadClientsImmatureTestInstructionAttributesFromCloudDB(userID string) (testInstructionAttributesMessage []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage, err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	// **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TIATTR.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructionAttributes\" TIATTR "
	sqlToExecute = sqlToExecute + "ORDER BY TIATTR.\"DomainUuid\" ASC, TIATTR.\"TestInstructionUuid\" ASC, TIATTR.\"TestInstructionAttributeUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "5f769af2-f75a-4ea6-8c3d-2108c9dfb9b7",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Variables to used when extract data from result set
	var tempTestInstructionAttributeInputMask string

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		immatureTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage{}

		err := rows.Scan(
			&immatureTestInstructionAttribute.DomainUuid,
			&immatureTestInstructionAttribute.DomainName,
			&immatureTestInstructionAttribute.TestInstructionUuid,
			&immatureTestInstructionAttribute.TestInstructionName,
			&immatureTestInstructionAttribute.TestInstructionAttributeUuid,
			&immatureTestInstructionAttribute.TestInstructionAttributeName,
			&immatureTestInstructionAttribute.TestInstructionAttributeDescription,
			&immatureTestInstructionAttribute.TestInstructionAttributeMouseOver,
			&immatureTestInstructionAttribute.TestInstructionAttributeTypeUuid,
			&immatureTestInstructionAttribute.TestInstructionAttributeTypeName,
			&immatureTestInstructionAttribute.TestInstructionAttributeValueAsString,
			&immatureTestInstructionAttribute.TestInstructionAttributeValueUuid,
			&immatureTestInstructionAttribute.TestInstructionAttributeVisible,
			&immatureTestInstructionAttribute.TestInstructionAttributeEnable,
			&immatureTestInstructionAttribute.TestInstructionAttributeMandatory,
			&immatureTestInstructionAttribute.TestInstructionAttributeVisibleInTestCaseArea,
			&immatureTestInstructionAttribute.TestInstructionAttributeIsDeprecated,
			&tempTestInstructionAttributeInputMask,
			&immatureTestInstructionAttribute.TestInstructionAttributeUIType,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "7cd322cb-2219-4c4d-a8c8-2770a42b0c23",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Add BondAttribute to BondsAttributes
		testInstructionAttributesMessage = append(testInstructionAttributesMessage, &immatureTestInstructionAttribute)

	}

	return testInstructionAttributesMessage, err
}
