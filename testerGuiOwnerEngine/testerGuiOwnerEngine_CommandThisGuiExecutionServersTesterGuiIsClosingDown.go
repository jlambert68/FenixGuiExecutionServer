package testerGuiOwnerEngine

import "FenixGuiExecutionServer/common_config"

// Process the actual command 'ChannelCommand_ThisGuiExecutionServersTesterGuiIsClosingDown'
func commandThisGuiExecutionServersTesterGuiIsClosingDown(
	testerGuiIsClosingDown *common_config.TesterGuiIsClosingDownStruct) {

	// Delete All Subscription, for specific TesterGui, from the Subscriptions-Map
	deleteTesterGuiFromTestCaseExecutionsSubscriptionFromMap(testerGuiIsClosingDown.TesterGuiApplicationId)

	// Broadcast message to other GuiExecutionServer
	broadcastSenderForChannelMessage_ThisGuiExecutionServersTesterGuiIsClosingDownMessage(*testerGuiIsClosingDown)

}
