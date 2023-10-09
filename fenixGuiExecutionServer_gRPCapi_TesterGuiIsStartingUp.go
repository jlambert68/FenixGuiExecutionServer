package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// TesterGuiIsStartingUp
// TesterGui informs that it is closing down
func (s *fenixGuiExecutionServerGrpcServicesServer) TesterGuiIsStartingUp(
	ctx context.Context,
	userAndApplicationRunTimeIdentificationMessage *fenixExecutionServerGuiGrpcApi.UserAndApplicationRunTimeIdentificationMessage) (
	ackNackResponse *fenixExecutionServerGuiGrpcApi.AckNackResponse,
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "e89b9b23-f7be-4534-8224-4896688fedb7",
	}).Debug("Incoming 'gRPC - TesterGuiIsStartingUp'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "2d6e6ac8-a766-45c2-8df3-9d93d9084932",
	}).Debug("Outgoing 'gRPC - TesterGuiIsStartingUp'")

	// Current user
	userID := userAndApplicationRunTimeIdentificationMessage.UserId

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userID,
		userAndApplicationRunTimeIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Exiting
		return returnMessage, nil
	}

	// Put message on 'testGuiExecutionEngineChannel' to be processed
	var tempTesterGuiIsClosingDown common_config.TesterGuiIsClosingDownStruct
	tempTesterGuiIsClosingDown = common_config.TesterGuiIsClosingDownStruct{
		TesterGuiApplicationId:          userAndApplicationRunTimeIdentificationMessage.ApplicationRunTimeUuid,
		UserId:                          userID,
		GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
		MessageTimeStamp:                time.Now(),
	}

	var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
	testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
		TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_ThisGuiExecutionServersTesterGuiIsStartingUp,
		TesterGuiIsClosingDown:                                &tempTesterGuiIsClosingDown,
		GuiExecutionServerIsClosingDown:                       nil,
		UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
		GuiExecutionServerIsStartingUp:                        nil,
		GuiExecutionServerStartedUpTimeStampRefresher:         nil,
		UserSubscribesToUserAndTestCaseExecutionCombination:   nil,
	}

	// Put on EngineChannel
	common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

	// Create Return message
	returnMessage = &fenixExecutionServerGuiGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
	}

	return returnMessage, nil
}
