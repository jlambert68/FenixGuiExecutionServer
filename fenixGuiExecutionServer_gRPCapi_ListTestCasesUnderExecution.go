package main

import (
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListAllImmatureTestInstructionAttributes - *********************************************************************
// The TestCase Builder asks for all TestInstructions Attributes that the user must add values to in TestCase
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) ListTestCasesUnderExecution(ctx context.Context, listTestCasesUnderExecutionRequest *fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionRequest) (*fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse, error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "2fa84f86-4ccc-479c-b73d-c329897c1873",
	}).Debug("Incoming 'gRPC - ListTestCasesUnderExecution'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "567a589f-d832-4e22-80a7-b37179eb6f94",
	}).Debug("Outgoing 'gRPC - ListTestCasesUnderExecution'")

	// Current user
	userID := listTestCasesUnderExecutionRequest.UserIdentification.UserId

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiExecutionServerObject.isClientUsingCorrectTestDataProtoFileVersion(userID, fenixExecutionServerGuiGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(listTestCasesUnderExecutionRequest.UserIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse{
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
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse{
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
	responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse{
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
