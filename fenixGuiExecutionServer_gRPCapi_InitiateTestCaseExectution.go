package main

import (
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// InitiateTestCaseExecution - *********************************************************************
// Initiate a TestExecution from a TestCase and a TestDataSet
func (s *fenixExecutionServerGuiGrpcServicesServer) InitiateTestCaseExecution(ctx context.Context, initiateSingleTestCaseExecutionRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionRequestMessage) (*fenixExecutionServerGuiGrpcApi.AckNackResponse, error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a93fb1bd-1a5b-4417-80c3-082d34267c06",
	}).Debug("Incoming 'gRPC - InitiateTestCaseExecution'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "981ad10a-2bfb-4a39-9b4d-35cac0d7481a",
	}).Debug("Outgoing 'gRPC - InitiateTestCaseExecution'")

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiExecutionServerObject.isClientUsingCorrectTestDataProtoFileVersion(initiateSingleTestCaseExecutionRequestMessage.UserIdentification.UserId, initiateSingleTestCaseExecutionRequestMessage.UserIdentification.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		// Exiting
		return returnMessage, nil
	}

	// Save Pinned TestInstructions and pre-created TestInstructionContainers to Cloud DB
	returnMessage = fenixGuiExecutionServerObject.prepareInitiateTestCaseExecutionSaveToCloudDB(initiateSingleTestCaseExecutionRequestMessage)
	if returnMessage != nil {
		// Something went wrong when saving to database
		// Exiting
		return returnMessage, nil
	}

	return &fenixExecutionServerGuiGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiExecutionServerObject.getHighestFenixTestDataProtoFileVersion()),
	}, nil
}
