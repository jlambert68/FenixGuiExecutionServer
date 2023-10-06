package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"strconv"
)

// Process the actual command 'ChannelCommand_AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination'
func commandAnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination(
	anotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination *common_config.AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombinationStruct) {

	// When sender is other GuiExecutionServer then remove the subscription from the map
	if anotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination.GuiExecutionServerApplicationId != common_config.ApplicationRunTimeUuid {

		// Create Key used for 'testCaseExecutionsSubscriptionsMap'
		var testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType
		testCaseExecutionsSubscriptionsMapKey = testCaseExecutionsSubscriptionsMapKeyType(
			anotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination.TestCaseExecutionUuid +
				strconv.Itoa(int(anotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination.TestCaseExecutionVersion)))

		// Remove this responsibility subscription
		deleteTestCaseExecutionsSubscriptionFromMap(testCaseExecutionsSubscriptionsMapKey)
	}

}
