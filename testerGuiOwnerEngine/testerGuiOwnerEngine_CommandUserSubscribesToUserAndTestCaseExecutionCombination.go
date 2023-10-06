package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"strconv"
)

// Process the actual command 'ChannelCommand_UserSubscribesToUserAndTestCaseExecutionCombination'
func commandUserSubscribesToUserAndTestCaseExecutionCombination(
	userSubscribesToUserAndTestCaseExecutionCombination *common_config.UserSubscribesToUserAndTestCaseExecutionCombinationStruct) {

	// When sender is this GuiExecutionServer then add the subscription to the map
	if userSubscribesToUserAndTestCaseExecutionCombination.GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

		var guiExecutionServerResponsibility *common_config.GuiExecutionServerResponsibilityStruct
		guiExecutionServerResponsibility = &common_config.GuiExecutionServerResponsibilityStruct{
			TesterGuiApplicationId:   userSubscribesToUserAndTestCaseExecutionCombination.TesterGuiApplicationId,
			UserId:                   userSubscribesToUserAndTestCaseExecutionCombination.UserId,
			TestCaseExecutionUuid:    userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid,
			TestCaseExecutionVersion: userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion,
		}

		// Create Key used for 'testCaseExecutionsSubscriptionsMap'
		var testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType
		testCaseExecutionsSubscriptionsMapKey = testCaseExecutionsSubscriptionsMapKeyType(
			userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid +
				strconv.Itoa(int(userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion)))

		// Save this responsibility
		saveToTestCaseExecutionsSubscriptionToMap(
			testCaseExecutionsSubscriptionsMapKey, guiExecutionServerResponsibility)

		// Inform Other GuiExecutionServers to remove this Key from their maps
		// Create channel message
		var tempGuiExecutionServerStartedUpTimeStampRefresher common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
		tempGuiExecutionServerStartedUpTimeStampRefresher = common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct{
			GuiExecutionServerApplicationId: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
				GuiExecutionServerIsStartingUp.GuiExecutionServerApplicationId,
			MessageTimeStamp: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
				GuiExecutionServerIsStartingUp.MessageTimeStamp,
		}

		// Put message on 'testGuiExecutionEngineChannel' to be processed
		var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
		testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
			TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_ThisGuiExecutionServerSendsStartedUpTimeStamp,
			TesterGuiIsClosingDown:                                nil,
			GuiExecutionServerIsClosingDown:                       nil,
			UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
			GuiExecutionServerIsStartingUp:                        nil,
			GuiExecutionServerStartedUpTimeStampRefresher:         &tempGuiExecutionServerStartedUpTimeStampRefresher,
		}

		// Put on EngineChannel
		common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

	}

}
