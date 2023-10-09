package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// UnSubscribeToMessages
// TesterGui unsubscribes to status updates for specific TestCaseExecutions
func (s *fenixGuiExecutionServerGrpcServicesServer) UnSubscribeToMessages(
	ctx context.Context,
	unSubscribeToMessagesRequest *fenixExecutionServerGuiGrpcApi.UnSubscribeToMessagesRequest) (
	ackNackResponse *fenixExecutionServerGuiGrpcApi.AckNackResponse,
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "4e609f22-ebfa-42ad-8c6e-528b56f3118a",
	}).Debug("Incoming 'gRPC - UnSubscribeToMessages'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "2a3f183c-58ba-4042-9ea4-4c206c47179f",
	}).Debug("Outgoing 'gRPC - UnSubscribeToMessages'")

	// Current user
	userID := unSubscribeToMessagesRequest.ApplicationRunTimeIdentification.UserId

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userID,
		unSubscribeToMessagesRequest.ApplicationRunTimeIdentification.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Exiting
		return returnMessage, nil
	}

	// Loop all TestCaseExecutions, to subscribe to, and put them 'testGuiExecutionEngineChannel' to be processed
	for _, tempTestCaseExecutionsStatusUnSubscriptions := range unSubscribeToMessagesRequest.TestCaseExecutionsStatusUnSubscriptions {

		var tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination common_config.
			UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
		tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination = common_config.
			UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{
			TesterGuiApplicationId:          unSubscribeToMessagesRequest.ApplicationRunTimeIdentification.ApplicationRunTimeUuid,
			UserId:                          userID,
			GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
			TestCaseExecutionUuid:           tempTestCaseExecutionsStatusUnSubscriptions.GetTestCaseExecutionUuid(),
			TestCaseExecutionVersion:        tempTestCaseExecutionsStatusUnSubscriptions.GetTestCaseExecutionVersion(),
			MessageTimeStamp:                time.Now(),
		}

		var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
		testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
			TesterGuiOwnerEngineChannelCommand:                                 common_config.ChannelCommand_ThisGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination,
			TesterGuiIsClosingDown:                                             nil,
			GuiExecutionServerIsClosingDown:                                    nil,
			UserUnsubscribesToUserAndTestCaseExecutionCombination:              &tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination,
			GuiExecutionServerIsStartingUp:                                     nil,
			GuiExecutionServerStartedUpTimeStampRefresher:                      nil,
			UserSubscribesToUserAndTestCaseExecutionCombination:                nil,
			AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination: nil}

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
