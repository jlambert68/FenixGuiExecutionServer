package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"strconv"
)

// Process the actual command 'ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination'
func commandUserUnsubscribesToUserAndTestCaseExecutionCombination(
	userUnsubscribesToUserAndTestCaseExecutionCombination *common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct) {

	// When sender is this GuiExecutionServer then remove the subscription from the map
	if userUnsubscribesToUserAndTestCaseExecutionCombination.GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

		// Create Key used for 'testCaseExecutionsSubscriptionsMap'
		var testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType
		testCaseExecutionsSubscriptionsMapKey = testCaseExecutionsSubscriptionsMapKeyType(
			userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid +
				strconv.Itoa(int(userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion)))

		// Remove this responsibility subscription
		deleteTestCaseExecutionsSubscriptionFromMap(testCaseExecutionsSubscriptionsMapKey)
	}

}
