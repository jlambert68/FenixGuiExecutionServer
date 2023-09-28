package testerGuiOwnerEngine

// InitiateTesterGuiOwnerEngine
// Initiate the channel reader which is used handling which GuiExecutionServer that is responsible for which TesterGui,
// regarding status-sending
// Initiate BroadcastListeners for Channel 1 and Channel 2
func InitiateTesterGuiOwnerEngine() {

	go InitiateAndStartBroadcastChannelListenerEngine()
	go startTesterGuiOwnerEngineChannelReader()

}
