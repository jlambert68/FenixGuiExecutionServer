package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SubscribeToMessages
// TesterGui subscribes to status updates for specific TestCaseExecutions
func (s *fenixGuiExecutionServerGrpcServicesServer) SubscribeToMessages(
	ctx context.Context,
	subscribeToMessagesRequest *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesRequest) (
	ackNackResponse *fenixExecutionServerGuiGrpcApi.AckNackResponse,
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "2fdb2a0a-67f2-430f-a37e-8457ded92657",
	}).Debug("Incoming 'gRPC - SubscribeToMessages'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "504e3f96-54c3-42e7-a22e-cb6b1c7bfcb5",
	}).Debug("Outgoing 'gRPC - SubscribeToMessages'")

	// Current user
	userIdOnComputer := subscribeToMessagesRequest.ApplicationRunTimeIdentification.GetUserIdOnComputer()

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdOnComputer,
		subscribeToMessagesRequest.ApplicationRunTimeIdentification.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Exiting
		return returnMessage, nil
	}

	// Loop all TestCaseExecutions, to subscribe to, and put them 'testGuiExecutionEngineChannel' to be processed
	for _, tempTestCaseExecutionsStatusSubscriptions := range subscribeToMessagesRequest.TestCaseExecutionsStatusSubscriptions {

		var tempUserSubscribesToUserAndTestCaseExecutionCombination common_config.
			UserSubscribesToUserAndTestCaseExecutionCombinationStruct
		tempUserSubscribesToUserAndTestCaseExecutionCombination = common_config.
			UserSubscribesToUserAndTestCaseExecutionCombinationStruct{
			TesterGuiApplicationId: subscribeToMessagesRequest.
				ApplicationRunTimeIdentification.ApplicationRunTimeUuid,
			UserId:                          userIdOnComputer,
			GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
			TestCaseExecutionUuid:           tempTestCaseExecutionsStatusSubscriptions.GetTestCaseExecutionUuid(),
			TestCaseExecutionVersion:        tempTestCaseExecutionsStatusSubscriptions.GetTestCaseExecutionVersion(),
			MessageTimeStamp:                time.Now(),
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

		// Put on EngineChannel
		common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

	}

	// Create Return message
	returnMessage = &fenixExecutionServerGuiGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
	}

	return returnMessage, nil
}
