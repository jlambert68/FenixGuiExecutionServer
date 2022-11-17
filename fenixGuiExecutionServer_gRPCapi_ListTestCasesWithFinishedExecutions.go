package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListAllImmatureTestInstructionAttributes - *********************************************************************
// The TestCase Builder asks for all TestInstructions Attributes that the user must add values to in TestCase
func (s *fenixGuiExecutionServerGrpcServicesServer) ListTestCasesWithFinishedExecutions(ctx context.Context, listTestCasesWithFinishedExecutionsRequest *fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsRequest) (*fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse, error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "33451c5f-1230-4ea1-817e-08afa1c1192b",
	}).Debug("Incoming 'gRPC - ListTestCasesWithFinishedExecutions'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "1cabcbc9-86ab-4ffb-b41f-8562c3bc1d75",
	}).Debug("Outgoing 'gRPC - ListTestCasesWithFinishedExecutions'")

	// Current user
	userID := listTestCasesWithFinishedExecutionsRequest.UserAndApplicationRunTimeIdentification.UserId

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userID,
		fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
			listTestCasesWithFinishedExecutionsRequest.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
			AckNackResponse:               returnMessage,
			TestCaseWithFinishedExecution: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var listTestCasesWithFinishedExecutionsResponse []*fenixExecutionServerGuiGrpcApi.TestCaseWithFinishedExecutionMessage

	// Get users ImmatureTestInstruction-data from CloudDB
	listTestCasesWithFinishedExecutionsResponse, err := fenixGuiExecutionServerObject.listTestCasesWithFinishedExecutionsLoadFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ListTestCasesWithFinishedExecutions from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCaseWithFinishedExecution: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
		TestCaseWithFinishedExecution: listTestCasesWithFinishedExecutionsResponse,
	}

	return responseMessage, nil
}
