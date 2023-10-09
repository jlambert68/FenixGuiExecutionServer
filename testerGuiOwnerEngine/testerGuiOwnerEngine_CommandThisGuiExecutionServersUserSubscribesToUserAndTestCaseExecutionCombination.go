package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/outgoingPubSubMessages"
	"strconv"
)

// Process the actual command 'ChannelCommand_ThisGuiExecutionServersUserSubscribesToUserAndTestCaseExecutionCombination'
func commandThisGuiExecutionServersUserSubscribesToUserAndTestCaseExecutionCombination(
	userSubscribesToUserAndTestCaseExecutionCombination *common_config.UserSubscribesToUserAndTestCaseExecutionCombinationStruct) {

	// When sender is this GuiExecutionServer then add the subscription to the map
	if userSubscribesToUserAndTestCaseExecutionCombination.GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

		var guiExecutionServerResponsibility *common_config.GuiExecutionServerResponsibilityStruct
		guiExecutionServerResponsibility = &common_config.GuiExecutionServerResponsibilityStruct{
			TesterGuiApplicationId:   userSubscribesToUserAndTestCaseExecutionCombination.TesterGuiApplicationId,
			UserId:                   userSubscribesToUserAndTestCaseExecutionCombination.UserId,
			TestCaseExecutionUuid:    userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid,
			TestCaseExecutionVersion: userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion,
		}

		// Create Key used for 'testCaseExecutionsSubscriptionsMap'
		var testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType
		testCaseExecutionsSubscriptionsMapKey = testCaseExecutionsSubscriptionsMapKeyType(
			userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid +
				strconv.Itoa(int(userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion)))

		// Save this responsibility
		saveToTestCaseExecutionsSubscriptionToMap(
			testCaseExecutionsSubscriptionsMapKey, guiExecutionServerResponsibility)

		// Check if PubSub-Topic already exists
		var pubSubTopicToLookFor string
		pubSubTopicToLookFor = generatePubSubTopicForExecutionStatusUpdates(
			userSubscribesToUserAndTestCaseExecutionCombination.TesterGuiApplicationId)

		// Secure that PubSub exist, if not then creat both PubSubTopic and PubSubTopic-Subscription
		outgoingPubSubMessages.CreateTopicDeadLettingAndSubscriptionIfNotExists(pubSubTopicToLookFor)

		// Inform other GuiExecutionServers to remove this Key from their maps
		// Create message
		var tempUserSubscribesToUserAndTestCaseExecutionCombination common_config.
			UserSubscribesToUserAndTestCaseExecutionCombinationStruct
		tempUserSubscribesToUserAndTestCaseExecutionCombination = common_config.
			UserSubscribesToUserAndTestCaseExecutionCombinationStruct{
			TesterGuiApplicationId:          userSubscribesToUserAndTestCaseExecutionCombination.TesterGuiApplicationId,
			UserId:                          userSubscribesToUserAndTestCaseExecutionCombination.UserId,
			GuiExecutionServerApplicationId: common_config.ApplicationRunTimeUuid,
			TestCaseExecutionUuid:           userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid,
			TestCaseExecutionVersion:        userSubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion,
			MessageTimeStamp:                userSubscribesToUserAndTestCaseExecutionCombination.MessageTimeStamp,
		}

		// Send message to be broadcasted to other GuiExecutionServers
		broadcastSenderForChannelMessage_ThisGuiExecutionServersTesterGuiSubscribesToThisTestCaseExecutionCombination(
			tempUserSubscribesToUserAndTestCaseExecutionCombination)

	}
}
