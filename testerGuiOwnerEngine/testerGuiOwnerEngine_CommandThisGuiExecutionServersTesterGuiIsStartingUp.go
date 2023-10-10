package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/outgoingPubSubMessages"
)

// Process the actual command 'ChannelCommand_ThisGuiExecutionServersTesterGuiIsStartingUp'
func commandThisGuiExecutionServersTesterGuiIsStartingUp(
	testerGuiIsStartingUp *common_config.TesterGuiIsStartingUpStruct) {

	// Create PubSub-Topic
	var pubSubTopicToLookFor string
	pubSubTopicToLookFor = generatePubSubTopicForExecutionStatusUpdates(testerGuiIsStartingUp.UserId)

	// Secure that PubSub Topic, DeadLetteringTopic and their Subscriptions exist
	var err error
	err = outgoingPubSubMessages.CreateTopicDeadLettingAndSubscriptionIfNotExists(pubSubTopicToLookFor)
	if err != nil {
		return
	}

	// Broadcast message to other GuiExecutionServer
	broadcastSenderForChannelMessage_ThisGuiExecutionServersTesterGuiIsStartingUp(*testerGuiIsStartingUp)

}
