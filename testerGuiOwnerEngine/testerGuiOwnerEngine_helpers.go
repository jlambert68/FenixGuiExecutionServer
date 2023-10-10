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

	// Add this GuiExecutionServer to slice with all other GuiExecutionServers that started after this one
	var tempGuiExecutionServerStartUpOrder *guiExecutionServerStartUpOrderStruct
	tempGuiExecutionServerStartUpOrder = &guiExecutionServerStartUpOrderStruct{
		applicationRunTimeUuid:        common_config.ApplicationRunTimeUuid,
		applicationRunTimeStartUpTime: common_config.ApplicationRunTimeStartUpTime,
	}
	insertGuiExecutionServerIntoTimeOrderedSlice(tempGuiExecutionServerStartUpOrder)

	// Start up broadcast Listener engine, used for receiving messages from other GuiExecutionServer
	go InitiateAndStartBroadcastChannelListenerEngine()

	// Start up GuiOwnerEngine
	go startTesterGuiOwnerEngineChannelReader()

	// Start up periodic broadcaster for StartUp-timestamp
	reInformOtherGuiExecutionServersAboutThatThisGuiExecutionServersStartingUpStartUpTimeStamp()

}

// Inform other running GuiExecutionServers that this server is starting up
func informOtherGuiExecutionServersThatThisGuiExecutionServerIsStartingUp() {

	// Put message on 'testGuiExecutionEngineChannel' to be processed
	var tempGuiExecutionServerIsStartingUp common_config.GuiExecutionServerIsStartingUpStruct
	tempGuiExecutionServerIsStartingUp = common_config.GuiExecutionServerIsStartingUpStruct{
		GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
		MessageTimeStamp:                common_config.ApplicationRunTimeStartUpTime,
	}

	var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
	testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
		TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_ThisGuiExecutionServerIsStartingUp,
		TesterGuiIsClosingDown:                                nil,
		GuiExecutionServerIsClosingDown:                       nil,
		UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
		GuiExecutionServerIsStartingUp:                        &tempGuiExecutionServerIsStartingUp,
		GuiExecutionServerStartedUpTimeStampRefresher:         nil,
	}

	// Put on GuiOwnerEngineChannel
	common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

}

// re-inform, in a periodic manner, other running GuiExecutionServers about this server's StartingUp-timestamp
func reInformOtherGuiExecutionServersAboutThatThisGuiExecutionServersStartingUpStartUpTimeStamp() {

	go func() {
		for {
			// Sleep for a while before broadcasting
			time.Sleep(timeStampBroadcastDuration)

			// Put message on 'testGuiExecutionEngineChannel' to be processed
			var tempGuiExecutionServerStartedUpTimeStampRefresher common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct
			tempGuiExecutionServerStartedUpTimeStampRefresher = common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct{
				GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
				MessageTimeStamp:                common_config.ApplicationRunTimeStartUpTime,
			}

			var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
			testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
				TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_ThisGuiExecutionServerSendsStartedUpTimeStamp,
				TesterGuiIsClosingDown:                                nil,
				GuiExecutionServerIsClosingDown:                       nil,
				UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
				GuiExecutionServerIsStartingUp:                        nil,
				GuiExecutionServerStartedUpTimeStampRefresher:         &tempGuiExecutionServerStartedUpTimeStampRefresher,
			}

			// Put on GuiOwnerEngineChannel
			common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

		}
	}()
}

// Create the PubSub-topic from TesterGui-ApplicationUuid
func generatePubSubTopicForExecutionStatusUpdates(testerGuiUserId string) (statusExecutionTopic string) {

	var pubSubTopicBase string
	pubSubTopicBase = common_config.TestExecutionStatusPubSubTopicBase

	// Get the first 8 characters from TesterGui-ApplicationUuid
	// var shortedAppUuid string
	// shortedAppUuid = testerGuiApplicationUuid[0:8]

	// Build PubSub-topic
	statusExecutionTopic = pubSubTopicBase + "-" + testerGuiUserId

	return statusExecutionTopic
}
