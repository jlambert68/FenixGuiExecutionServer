package testerGuiOwnerEngine

import "FenixGuiExecutionServer/common_config"

// Process the actual command 'ChannelCommand_ThisGuiExecutionServerIsStartingUp'
func commandThisGuiExecutionServerIsStartingUp(
	guiExecutionServerIsStartingUp *common_config.GuiExecutionServerIsStartingUpStruct) {

	// Broadcast message to other GuiExecutionServer
	broadcastSenderForChannelMessage_ThisGuiExecutionServerIsStartingUp(*guiExecutionServerIsStartingUp)

}
