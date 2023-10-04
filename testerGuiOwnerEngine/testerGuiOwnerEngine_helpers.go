package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"time"
)

// InitiateTesterGuiOwnerEngine
// Initiate the channel reader which is used handling which GuiExecutionServer that is responsible for which TesterGui,
// regarding status-sending
// Initiate BroadcastListeners for Channel 1 and Channel 2
func InitiateTesterGuiOwnerEngine() {

	go InitiateAndStartBroadcastChannelListenerEngine()
	go startTesterGuiOwnerEngineChannelReader()

}

// Inform other running GuiExecutionServers that this server is starting up
func informOtherGuiExecutionServersThatThisGuiExecutionServerIsStartingUp() {

	// Put message on 'testGuiExecutionEngineChannel' to be processed
	var tempGuiExecutionServerIsStartingUp common_config.GuiExecutionServerIsStartingUpStruct
	tempGuiExecutionServerIsStartingUp = common_config.GuiExecutionServerIsStartingUpStruct{
		GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
		MessageTimeStamp:                time.Now(),
	}

	var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
	testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
		TesterGuiOwnerEngineChannelCommand:                                 common_config.ChannelCommand_ThisGuiExecutionServerIsClosingDown,
		TesterGuiIsClosingDown:                                             nil,
		GuiExecutionServerIsClosingDown:                                    nil,
		ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: nil,
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              nil,
		GuiExecutionServerIsStartingUp:                                     &tempGuiExecutionServerIsStartingUp,
	}

	// Put on GuiOwnerEngineChannel
	common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

}
