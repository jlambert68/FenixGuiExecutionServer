package main

import (
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"


	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// ListAllImmatureTestInstructionAttributes - *********************************************************************
// The TestCase Builder asks for all TestInstructions Attributes that the user must add values to in TestCase
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) ListAllSingleTestCaseExecutions(ctx context.Context, listAllSingleTestCaseExecutionsRequest *fenixExecutionServerGuiGrpcApi.ListAllSingleTestCaseExecutionsRequest) (*fenixExecutionServerGuiGrpcApi.ListAllSingleTestCaseExecutionsResponse, error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListAllSingleTestCaseExecutionsResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - ListAllSingleTestCaseExecutions'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - ListAllSingleTestCaseExecutions'")

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiExecutionServerObject.isClientUsingCorrectTestDataProtoFileVersion(listAllSingleTestCaseExecutionsRequest.UserIdentification.UserId, fenixExecutionServerGuiGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(listAllSingleTestCaseExecutionsRequest.UserIdentification.ProtoFileVersionUsedByClient)))
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListAllSingleTestCaseExecutionsResponse{
			SingleTestCaseExecutionSummary: nil,
			A:                returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Current user
	userID := userIdentificationMessage.UserId

	// Define variables to store data from DB in
	var testInstructionAttributesList []*fenixExecutionServerGuiGrpcApi.SingleTestCaseExecutionSummaryMessage

	// Get users ImmatureTestInstruction-data from CloudDB
	testInstructionAttributesList, err := fenixGuiExecutionServerObject.loadClientsImmatureTestInstructionAttributesFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListAllSingleTestCaseExecutionsResponse{
			SingleTestCaseExecutionSummary: nil,
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ImmatureTestInstructionAttributes from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiExecutionServerObject.getHighestFenixTestDataProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixExecutionServerGuiGrpcApi.ImmatureTestInstructionAttributesMessage{
		TestInstructionAttributesList: testInstructionAttributesList,
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiExecutionServerObject.getHighestFenixTestDataProtoFileVersion()),
		},
	}

	return responseMessage, nil
}

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadClientsImmatureTestInstructionAttributesFromCloudDB(userID string) (testInstructionAttributesMessage []*fenixExecutionServerGuiGrpcApi.ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage, err error) {

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
		immatureTestInstructionAttribute := fenixExecutionServerGuiGrpcApi.ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage{}

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
