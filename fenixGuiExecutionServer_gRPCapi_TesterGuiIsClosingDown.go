package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// TesterGuiIsClosingDown
// TesterGui informs that it is closing down
func (s *fenixGuiExecutionServerGrpcServicesServer) TesterGuiIsClosingDown(
	ctx context.Context,
	userAndApplicationRunTimeIdentificationMessage *fenixExecutionServerGuiGrpcApi.UserAndApplicationRunTimeIdentificationMessage) (
	ackNackResponse *fenixExecutionServerGuiGrpcApi.AckNackResponse,
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a39cbc7a-7e49-460b-bed4-c2a2d6e0649c",
	}).Debug("Incoming 'gRPC - TesterGuiIsClosingDown'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "8c97ee72-e40e-4c7b-b17b-f59e254623b7",
	}).Debug("Outgoing 'gRPC - TesterGuiIsClosingDown'")

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
		TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_ThisGuiExecutionServersTesterGuiIsClosingDown,
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
