package testerGuiOwnerEngine

import "FenixGuiExecutionServer/common_config"

// Process the actual command 'ChannelCommand_AnotherGuiExecutionServersTesterGuiIsClosingDown'
func commandAnotherGuiExecutionServersTesterGuiIsClosingDown(
	testerGuiIsClosingDown *common_config.TesterGuiIsClosingDownStruct) {

	// Delete All Subscription, for specific TesterGui, from the Subscriptions-Map
	deleteTesterGuiFromTestCaseExecutionsSubscriptionFromMap(testerGuiIsClosingDown.TesterGuiApplicationId)

}
