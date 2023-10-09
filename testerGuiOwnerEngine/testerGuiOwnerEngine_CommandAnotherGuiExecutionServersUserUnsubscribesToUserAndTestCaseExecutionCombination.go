package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"strconv"
)

// Process the actual command 'ChannelCommand_AnotherGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination'
func commandAnotherGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination(
	userUnsubscribesToUserAndTestCaseExecutionCombination *common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct) {

	// Remove the subscription from the map
	// Create Key used for 'testCaseExecutionsSubscriptionsMap'
	var testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType
	testCaseExecutionsSubscriptionsMapKey = testCaseExecutionsSubscriptionsMapKeyType(
		userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid +
			strconv.Itoa(int(userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion)))

	// Remove this responsibility subscription
	deleteTestCaseExecutionsSubscriptionFromMap(testCaseExecutionsSubscriptionsMapKey)

}
