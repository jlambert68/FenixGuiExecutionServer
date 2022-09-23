package main

import (
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListAllImmatureTestInstructionAttributes - *********************************************************************
// The TestCase Builder asks for all TestInstructions Attributes that the user must add values to in TestCase
func (s *fenixExecutionServerGrpcServicesServer) ListTestCasesOnExecutionQueue(ctx context.Context, listTestCasesInExecutionQueueRequest *fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueRequest) (*fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse, error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a88c93c3-cc86-4b4b-86ca-a11dd606b242",
	}).Debug("Incoming 'gRPC - ListTestCasesOnExecutionQueue'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "1ce24c71-11ae-4f76-a473-ce794e4610e6",
	}).Debug("Outgoing 'gRPC - ListTestCasesOnExecutionQueue'")

	// Current user
	userID := listTestCasesInExecutionQueueRequest.UserIdentification.UserId

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiExecutionServerObject.isClientUsingCorrectTestDataProtoFileVersion(userID, fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(listTestCasesInExecutionQueueRequest.UserIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse{
			AckNackResponse:           returnMessage,
			TestCasesInExecutionQueue: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var testCaseExecutionBasicInformation []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage

	// Get users ImmatureTestInstruction-data from CloudDB
	testCaseExecutionBasicInformation, err := fenixGuiExecutionServerObject.listTestCasesOnExecutionQueueLoadFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving TestCaseExecutionBasicInformationMessage from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiExecutionServerObject.getHighestFenixTestDataProtoFileVersion()),
			},
			TestCasesInExecutionQueue: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiExecutionServerObject.getHighestFenixTestDataProtoFileVersion()),
		},
		TestCasesInExecutionQueue: testCaseExecutionBasicInformation,
	}

	return responseMessage, nil
}
