package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"strconv"
)

// Process the actual command 'ChannelCommand_ThisGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination'
func commandThisGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination(
	userUnsubscribesToUserAndTestCaseExecutionCombination *common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct) {

	// Remove the subscription from the map
	// Create Key used for 'testCaseExecutionsSubscriptionsMap'
	var testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType
	testCaseExecutionsSubscriptionsMapKey = testCaseExecutionsSubscriptionsMapKeyType(
		userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid +
			strconv.Itoa(int(userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion)))

	// Remove this responsibility subscription
	deleteTestCaseExecutionsSubscriptionFromMap(testCaseExecutionsSubscriptionsMapKey)

	// Inform other GuiExecutionServers to remove this Key from their maps
	// Create message
	var tempUserUnsubscribesToUserAndTestCaseExecutionCombination common_config.
		UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
	tempUserUnsubscribesToUserAndTestCaseExecutionCombination = common_config.
		UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{
		TesterGuiApplicationId:          userUnsubscribesToUserAndTestCaseExecutionCombination.TesterGuiApplicationId,
		UserId:                          userUnsubscribesToUserAndTestCaseExecutionCombination.UserId,
		GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
		TestCaseExecutionUuid:           userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid,
		TestCaseExecutionVersion:        userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion,
		MessageTimeStamp:                userUnsubscribesToUserAndTestCaseExecutionCombination.MessageTimeStamp,
	}

	// Send message to be broadcasted to other GuiExecutionServers
	broadcastSenderForChannelMessage_ThisGuiExecutionServersTesterGuiUnsubscribesToThisTestCaseExecutionCombination(
		tempUserUnsubscribesToUserAndTestCaseExecutionCombination)

}
