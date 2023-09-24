package testerGuiOwnerEngine

import "FenixGuiExecutionServer/common_config"

// initiateTesterGuiOwnerEngine
// Initiate the channel reader which is used handling which GuiExecutionServer that is responsible for which TesterGui, regarding status-sending
func InitiateTesterGuiOwnerEngine() {

	go common_config.TestInstructionExecutionTesterGuiOwnerEngineEngineObject.startTesterGuiOwnerEngineChannelReader()

}
