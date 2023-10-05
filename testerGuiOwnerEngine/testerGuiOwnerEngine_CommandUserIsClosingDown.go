package testerGuiOwnerEngine

import "FenixGuiExecutionServer/common_config"

// Process the actual command 'ChannelCommand_UserIsClosingDown'
func commandUserIsClosingDown(
	testerGuiIsClosingDown *common_config.TesterGuiIsClosingDownStruct) {

	// Broadcast message to other GuiExecutionServer
	broadcastSenderForChannelMessage_TesterGuiIsClosingDownMessage(*testerGuiIsClosingDown)

}
