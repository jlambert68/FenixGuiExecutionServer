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

	// Initiate variable holding Subscriptions handled by this GuiExecutionServer
	testCaseExecutionsSubscriptionsMap = make(map[testCaseExecutionsSubscriptionsMapKeyType]*common_config.GuiExecutionServerResponsibilityStruct)

	// Add this GuiExecutionServer to guiExecutionServerStartUpOrder slice
	var tempGuiExecutionServerStartUpOrder *guiExecutionServerStartUpOrderStruct
	tempGuiExecutionServerStartUpOrder = &guiExecutionServerStartUpOrderStruct{
		applicationRunTimeUuid:        common_config.ApplicationRunTimeUuid,
		applicationRunTimeStartUpTime: common_config.ApplicationRunTimeStartUpTime,
	}
	guiExecutionServerStartUpOrder = append(guiExecutionServerStartUpOrder, tempGuiExecutionServerStartUpOrder)

	// Start up broadcast Listener engine, used for receiving messages from other GuiExecutionServer
	go InitiateAndStartBroadcastChannelListenerEngine()

	// Start up GuiOwnerEngine
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
