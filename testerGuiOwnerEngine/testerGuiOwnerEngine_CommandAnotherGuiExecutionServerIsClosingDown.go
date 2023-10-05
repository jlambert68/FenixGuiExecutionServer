package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"strconv"
)

// Process the actual command 'ChannelCommand_AnotherGuiExecutionServerIsClosingDown'
func commandAnotherGuiExecutionServerIsClosingDown(
	tempGuiExecutionServerIsClosingDown *common_config.GuiExecutionServerIsClosingDownStruct) {

	// Verify that it is not this GuiExecutionServer in the message, if so then just exit
	if tempGuiExecutionServerIsClosingDown.GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

		return
	}

	// Delete GuiExecutionServer-information from slice with GuiExecutionServers
	removeGuiExecutionServerFromSlice(tempGuiExecutionServerIsClosingDown.GuiExecutionServerApplicationId)

	// Should this GuiExecutionServer take  over responsibilities from GuiExecutionServer that is closing down
	// Take over responsibilities when this GuiExecutionServer is the only item in the slice 'guiExecutionServerStartUpOrder'
	if len(guiExecutionServerStartUpOrder) == 1 {

		// Loop over the responsibilities
		for _, tempGuiExecutionServerResponsibility := range tempGuiExecutionServerIsClosingDown.GuiExecutionServerResponsibilities {

			var guiExecutionServerResponsibility *common_config.GuiExecutionServerResponsibilityStruct
			guiExecutionServerResponsibility = &common_config.GuiExecutionServerResponsibilityStruct{
				TesterGuiApplicationId:   tempGuiExecutionServerResponsibility.TesterGuiApplicationId,
				UserId:                   tempGuiExecutionServerResponsibility.UserId,
				TestCaseExecutionUuid:    tempGuiExecutionServerResponsibility.TestCaseExecutionUuid,
				TestCaseExecutionVersion: tempGuiExecutionServerResponsibility.TestCaseExecutionVersion,
			}

			// Create Key used for 'testCaseExecutionsSubscriptionsMap'
			var testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType
			testCaseExecutionsSubscriptionsMapKey = testCaseExecutionsSubscriptionsMapKeyType(
				tempGuiExecutionServerResponsibility.TestCaseExecutionUuid +
					strconv.Itoa(int(tempGuiExecutionServerResponsibility.TestCaseExecutionVersion)))

			// Save this responsibility
			saveToTestCaseExecutionsSubscriptionToMap(
				testCaseExecutionsSubscriptionsMapKey, guiExecutionServerResponsibility)
		}
	}

}
