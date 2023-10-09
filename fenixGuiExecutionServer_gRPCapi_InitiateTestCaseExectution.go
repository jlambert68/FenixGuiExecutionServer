package main

import (
	"FenixGuiExecutionServer/broadcastEngine_ExecutionStatusUpdate"
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/messagesToExecutionServer"
	"context"
	"fmt"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// InitiateTestCaseExecution - *********************************************************************
// Initiate a TestExecution from a TestCase and a TestDataSet
func (s *fenixGuiExecutionServerGrpcServicesServer) InitiateTestCaseExecution(ctx context.Context,
	initiateSingleTestCaseExecutionRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionRequestMessage) (
	*fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage, error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a93fb1bd-1a5b-4417-80c3-082d34267c06",
		"initiateSingleTestCaseExecutionRequestMessage": initiateSingleTestCaseExecutionRequestMessage,
	}).Debug("Incoming 'gRPC - InitiateTestCaseExecution'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "981ad10a-2bfb-4a39-9b4d-35cac0d7481a",
	}).Debug("Outgoing 'gRPC - InitiateTestCaseExecution'")

	// Check if Client is using correct proto files version
	ackNackRespons := common_config.IsClientUsingCorrectTestDataProtoFileVersion(initiateSingleTestCaseExecutionRequestMessage.UserAndApplicationRunTimeIdentification.UserId, initiateSingleTestCaseExecutionRequestMessage.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient)
	if ackNackRespons != nil {
		// Not correct proto-file version is used
		// Exiting
		initiateSingleTestCaseExecutionResponseMessage := fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse:           ackNackRespons,
		}

		return &initiateSingleTestCaseExecutionResponseMessage, nil
	}

	// Save TestCaseExecution in Cloud DB
	var initiateSingleTestCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.InitiateSingleTestCaseExecutionResponseMessage
	initiateSingleTestCaseExecutionResponseMessage = fenixGuiExecutionServerObject.prepareInitiateTestCaseExecutionSaveToCloudDB(initiateSingleTestCaseExecutionRequestMessage)

	// Exit due to error in saving TestCaseExecution in database
	if initiateSingleTestCaseExecutionResponseMessage.AckNackResponse.AckNack == false {
		return initiateSingleTestCaseExecutionResponseMessage, nil
	}

	// *********
	// Create a Subscription on this 'TestCaseExecution' for this 'TestGui'
	broadcastEngine_ExecutionStatusUpdate.AddSubscriptionForTestCaseExecutionToTesterGui(
		broadcastEngine_ExecutionStatusUpdate.ApplicationRunTimeUuidType(
			initiateSingleTestCaseExecutionRequestMessage.UserAndApplicationRunTimeIdentification.ApplicationRunTimeUuid),
		broadcastEngine_ExecutionStatusUpdate.TestCaseExecutionUuidType(
			initiateSingleTestCaseExecutionResponseMessage.TestCasesInExecutionQueue.TestCaseExecutionUuid),
		1)

	// ******
	// Add a Subscription, using new the new PubSub-system, on this 'TestCaseExecution' for this 'TesterGui'
	// Create message to be put on 'testGuiExecutionEngineChannel' to be processed
	var tempUserSubscribesToUserAndTestCaseExecutionCombination common_config.UserSubscribesToUserAndTestCaseExecutionCombinationStruct
	tempUserSubscribesToUserAndTestCaseExecutionCombination = common_config.UserSubscribesToUserAndTestCaseExecutionCombinationStruct{
		TesterGuiApplicationId: initiateSingleTestCaseExecutionRequestMessage.
			UserAndApplicationRunTimeIdentification.GetApplicationRunTimeUuid(),
		UserId: initiateSingleTestCaseExecutionRequestMessage.
			UserAndApplicationRunTimeIdentification.GetUserId(),
		GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
		TestCaseExecutionUuid: initiateSingleTestCaseExecutionResponseMessage.TestCasesInExecutionQueue.
			GetTestCaseExecutionUuid(),
		TestCaseExecutionVersion: int32(initiateSingleTestCaseExecutionResponseMessage.TestCasesInExecutionQueue.GetTestCaseExecutionVersion()),
		MessageTimeStamp:         time.Now(),
	}

	var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
	testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
		TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_ThisGuiExecutionServersUserSubscribesToUserAndTestCaseExecutionCombination,
		TesterGuiIsClosingDown:                                nil,
		GuiExecutionServerIsClosingDown:                       nil,
		UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
		GuiExecutionServerIsStartingUp:                        nil,
		GuiExecutionServerStartedUpTimeStampRefresher:         nil,
		UserSubscribesToUserAndTestCaseExecutionCombination:   &tempUserSubscribesToUserAndTestCaseExecutionCombination,
	}

	// Put on GuiOwnerEngineChannel
	common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

	// Send Execution to ExecutionServer
	go func() {
		// Prepare message to be sent to ExecutionServer
		var testCaseExecutionsToProcessMessage *fenixExecutionServerGrpcApi.TestCaseExecutionsToProcessMessage
		var testCaseExecutionToProcess *fenixExecutionServerGrpcApi.TestCaseExecutionToProcess
		var testCaseExecutionsToProcess []*fenixExecutionServerGrpcApi.TestCaseExecutionToProcess

		testCaseExecutionToProcess = &fenixExecutionServerGrpcApi.TestCaseExecutionToProcess{
			TestCaseExecutionsUuid: initiateSingleTestCaseExecutionResponseMessage.TestCasesInExecutionQueue.
				TestCaseExecutionUuid,
			TestCaseExecutionVersion: 1,
			ExecutionStatusReportLevel: fenixExecutionServerGrpcApi.ExecutionStatusReportLevelEnum(
				initiateSingleTestCaseExecutionRequestMessage.ExecutionStatusReportLevel),
		}
		testCaseExecutionsToProcess = append(testCaseExecutionsToProcess, testCaseExecutionToProcess)

		testCaseExecutionsToProcessMessage = &fenixExecutionServerGrpcApi.TestCaseExecutionsToProcessMessage{
			TestCaseExecutionsToProcess: testCaseExecutionsToProcess,
		}

		// Trigger ExecutionEngine to start process TestCase from TestCaseExecution-queue
		var sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse *fenixExecutionServerGrpcApi.
			AckNackResponse
		sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse = messagesToExecutionServer.
			MessagesToExecutionServerObject.SendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServer(
			testCaseExecutionsToProcessMessage)

		// If triggering ExecutionServer to read TestCaseExecutionQueue wasn't successful then change 'initiateSingleTestCaseExecutionResponseMessage'
		if sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse.AckNack == false {
			var ackNackResponseToRespond *fenixExecutionServerGuiGrpcApi.AckNackResponse
			ackNackResponseToRespond = &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:    initiateSingleTestCaseExecutionResponseMessage.AckNackResponse.AckNack,
				Comments:   fmt.Sprintf("Message from ExecutionServer is: '%s'", initiateSingleTestCaseExecutionResponseMessage.AckNackResponse.Comments),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.
					CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			}

			initiateSingleTestCaseExecutionResponseMessage.AckNackResponse = ackNackResponseToRespond

			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"id": "f847771c-7947-4bc8-8902-5aa7ac2c7f88",
				"initiateSingleTestCaseExecutionRequestMessage": initiateSingleTestCaseExecutionRequestMessage,
			}).Error("Problem when doing gRPC-call to FenixExecutionServer")

			return //initiateSingleTestCaseExecutionResponseMessage, nil
		}

	}()

	return initiateSingleTestCaseExecutionResponseMessage, nil

}
