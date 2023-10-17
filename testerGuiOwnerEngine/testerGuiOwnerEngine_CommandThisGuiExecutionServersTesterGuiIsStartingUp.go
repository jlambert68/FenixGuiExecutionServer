package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"github.com/jlambert68/FenixSyncShared/pubSubHelpers"
)

// Process the actual command 'ChannelCommand_ThisGuiExecutionServersTesterGuiIsStartingUp'
func commandThisGuiExecutionServersTesterGuiIsStartingUp(
	testerGuiIsStartingUp *common_config.TesterGuiIsStartingUpStruct) {

	// Create PubSub-Topic
	var pubSubTopicToLookFor string
	pubSubTopicToLookFor = GeneratePubSubTopicForExecutionStatusUpdates(testerGuiIsStartingUp.UserId)

	// Secure that PubSub Topic, DeadLetteringTopic and their Subscriptions exist
	var err error
	err = pubSubHelpers.CreateTopicDeadLettingAndSubscriptionIfNotExists(
		pubSubTopicToLookFor, common_config.TestExecutionStatusPubSubTopicSchema)
	if err != nil {
		return
	}

	// Broadcast message to other GuiExecutionServer
	broadcastSenderForChannelMessage_ThisGuiExecutionServersTesterGuiIsStartingUp(*testerGuiIsStartingUp)

}
