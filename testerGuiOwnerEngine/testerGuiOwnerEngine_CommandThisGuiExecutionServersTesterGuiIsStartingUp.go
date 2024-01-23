package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"github.com/jlambert68/FenixSyncShared/pubSubHelpers"
	"strings"
)

// Process the actual command 'ChannelCommand_ThisGuiExecutionServersTesterGuiIsStartingUp'
func commandThisGuiExecutionServersTesterGuiIsStartingUp(
	testerGuiIsStartingUp *common_config.TesterGuiIsStartingUpStruct) {

	// Create PubSub-Topic and remove characters that are not allowed
	var pubSubTopicToLookFor string
	pubSubTopicToLookFor = GeneratePubSubTopicForExecutionStatusUpdates(testerGuiIsStartingUp.UserId)
	pubSubTopicToLookFor = strings.ReplaceAll(pubSubTopicToLookFor, "@", "")

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
