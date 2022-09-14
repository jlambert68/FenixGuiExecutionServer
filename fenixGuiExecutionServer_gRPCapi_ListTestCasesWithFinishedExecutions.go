package main

import (
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListAllImmatureTestInstructionAttributes - *********************************************************************
// The TestCase Builder asks for all TestInstructions Attributes that the user must add values to in TestCase
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) ListTestCasesWithFinishedExecutions(ctx context.Context, listTestCasesWithFinishedExecutionsRequest *fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsRequest) (*fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse, error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - ListAllSingleTestCaseExecutions'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - ListAllSingleTestCaseExecutions'")

	// Current user
	userID := listTestCasesWithFinishedExecutionsRequest.UserIdentification.UserId

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiExecutionServerObject.isClientUsingCorrectTestDataProtoFileVersion(userID, fenixExecutionServerGuiGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(listTestCasesWithFinishedExecutionsRequest.UserIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
			SingleTestCaseExecutionSummary: nil,
			AckNackResponse:                returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var singleTestCaseExecutionSummary []*fenixExecutionServerGuiGrpcApi.SingleTestCaseExecutionSummaryMessage

	// Get users ImmatureTestInstruction-data from CloudDB
	singleTestCaseExecutionSummary, err := fenixGuiExecutionServerObject.loadSingleTestCaseExecutionSummaryFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
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
	responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
		SingleTestCaseExecutionSummary: singleTestCaseExecutionSummary,
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiExecutionServerObject.getHighestFenixTestDataProtoFileVersion()),
		},
	}

	return responseMessage, nil
}
