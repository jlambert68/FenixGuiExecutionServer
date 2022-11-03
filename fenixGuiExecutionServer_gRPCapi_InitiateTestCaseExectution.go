package main

import (
	"FenixGuiExecutionServer/messagesToExecutionServer"
	"context"
	"fmt"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// InitiateTestCaseExecution - *********************************************************************
// Initiate a TestExecution from a TestCase and a TestDataSet
func (s *fenixGuiExecutionServerGrpcServicesServer) InitiateTestCaseExecution(ctx context.Context, initiateSingleTestCaseExecutionRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionRequestMessage) (*fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage, error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a93fb1bd-1a5b-4417-80c3-082d34267c06",
	}).Debug("Incoming 'gRPC - InitiateTestCaseExecution'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "981ad10a-2bfb-4a39-9b4d-35cac0d7481a",
	}).Debug("Outgoing 'gRPC - InitiateTestCaseExecution'")

	// Check if Client is using correct proto files version
	ackNackRespons := fenixGuiExecutionServerObject.isClientUsingCorrectTestDataProtoFileVersion(initiateSingleTestCaseExecutionRequestMessage.UserIdentification.UserId, initiateSingleTestCaseExecutionRequestMessage.UserIdentification.ProtoFileVersionUsedByClient)
	if ackNackRespons != nil {
		// Not correct proto-file version is used
		// Exiting
		initiateSingleTestCaseExecutionResponseMessage := fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCaseExecutionUuid: "",
			AckNackResponse:       ackNackRespons,
		}

		return &initiateSingleTestCaseExecutionResponseMessage, nil
	}

	// Save TestCaseExecution in Cloud DB
	initiateSingleTestCaseExecutionResponseMessage := fenixGuiExecutionServerObject.prepareInitiateTestCaseExecutionSaveToCloudDB(initiateSingleTestCaseExecutionRequestMessage)

	// Exit due to error in saving TestCaseExecution in database
	if initiateSingleTestCaseExecutionResponseMessage.AckNackResponse.AckNack == false {
		return initiateSingleTestCaseExecutionResponseMessage, nil
	}

	// Trigger ExecutionEngine to start process TestCase from TestCaseExecution-queue
	var sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse *fenixExecutionServerGrpcApi.AckNackResponse
	sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse = messagesToExecutionServer.MessagesToExecutionServerObject.SendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServer()

	// If triggering ExecutionServer to read TestCaseExecutionQueue wasn't successful then change 'initiateSingleTestCaseExecutionResponseMessage'
	if sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse.AckNack == false {
		var ackNackResponseToRespond *fenixExecutionServerGuiGrpcApi.AckNackResponse
		ackNackResponseToRespond = &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      ackNackRespons.AckNack,
			Comments:                     fmt.Sprintf("Message from ExecutionServer is: '%s'", ackNackRespons.Comments),
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiExecutionServerObject.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		}

		initiateSingleTestCaseExecutionResponseMessage.AckNackResponse = ackNackResponseToRespond
	}

	return initiateSingleTestCaseExecutionResponseMessage, nil

}
