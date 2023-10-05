package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
)

// Process the actual command 'ChannelCommand_AnotherGuiExecutionServerIsStartingUp'
func commandAnotherGuiExecutionServerIsStartingUp(
	tempGuiExecutionServerIsStartingUp *common_config.GuiExecutionServerIsStartingUpStruct) {

	// Verify that it is not this GuiExecutionServer in the message, if so then just exit
	if tempGuiExecutionServerIsStartingUp.GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

		return
	}

	// Try to insert GuiExecutionServer-information into slice with GuiExecutionServers
	// Logic if it can be inserted is handled by the function itself
	var guiExecutionServerToBeInsert *guiExecutionServerStartUpOrderStruct
	guiExecutionServerToBeInsert = &guiExecutionServerStartUpOrderStruct{
		applicationRunTimeUuid:        tempGuiExecutionServerIsStartingUp.GuiExecutionServerApplicationId,
		applicationRunTimeStartUpTime: tempGuiExecutionServerIsStartingUp.MessageTimeStamp,
	}

	insertGuiExecutionServerIntoTimeOrderedSlice(guiExecutionServerToBeInsert)

}
