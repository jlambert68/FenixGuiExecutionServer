package testerGuiOwnerEngine

import "FenixGuiExecutionServer/common_config"

// Process the actual command 'ChannelCommand_ThisGuiExecutionServerIsClosingDown'
func commandThisGuiExecutionServerIsClosingDown(
	guiExecutionServerIsClosingDown *common_config.GuiExecutionServerIsClosingDownStruct) {

	// Extract the responsibilities for this GuiExecutionServer
	var guiExecutionServerResponsibilities []common_config.GuiExecutionServerResponsibilityStruct
	guiExecutionServerResponsibilities = listAllTestCaseExecutionsSubscriptionsFromMap()

	// Add this GuiExecutionServer's responsibilities to the message to be broadcast
	guiExecutionServerIsClosingDown.GuiExecutionServerResponsibilities = guiExecutionServerResponsibilities

	// Broadcast message to other GuiExecutionServer
	broadcastSenderForChannelMessage_ThisGuiExecutionServerIsClosingDownMessage(*guiExecutionServerIsClosingDown)

}
