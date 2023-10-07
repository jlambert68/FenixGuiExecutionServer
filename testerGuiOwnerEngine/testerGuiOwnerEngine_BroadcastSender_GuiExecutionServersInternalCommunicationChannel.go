package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"encoding/json"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

//TesterGuiIsClosingDownMessage

// Broadcast message to all other GuiExecutionServers that 'TesterGuiIsClosingDownMessage'
func broadcastSenderForChannelMessage_TesterGuiIsClosingDownMessage(
	tempGuiExecutionServerIsClosingDown common_config.TesterGuiIsClosingDownStruct) (
	err error) {

	// Convert into Broadcast message type
	var broadcastMessageForGuiExecutionServersInternalCommunicationChannel BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct
	broadcastMessageForGuiExecutionServersInternalCommunicationChannel = BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct{
		GuiExecutionServersInternalCommunicationChannelType:                TesterGuiIsClosingDownMessage,
		TesterGuiIsClosingDown:                                             tempGuiExecutionServerIsClosingDown,
		GuiExecutionServerIsClosingDown:                                    common_config.GuiExecutionServerIsClosingDownStruct{},
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{},
		GuiExecutionServerSendStartedUpTimeStamp:                           common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct{},
		AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination: common_config.AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombinationStruct{},
	}

	// Broadcast message via Channel on Broadcast system
	err = broadcastSenderForGuiExecutionServersInternalCommunicationChannel(&broadcastMessageForGuiExecutionServersInternalCommunicationChannel)

	return err
}

// Broadcast message to all other GuiExecutionServers that 'GuiExecutionServerIsClosingDownMessage'
func broadcastSenderForChannelMessage_GuiExecutionServerIsClosingDownMessage(
	tempGuiExecutionServerIsClosingDown common_config.GuiExecutionServerIsClosingDownStruct) (
	err error) {

	// Convert into Broadcast message type
	var broadcastMessageForGuiExecutionServersInternalCommunicationChannel BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct
	broadcastMessageForGuiExecutionServersInternalCommunicationChannel = BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct{
		GuiExecutionServersInternalCommunicationChannelType:                GuiExecutionServerIsClosingDownMessage,
		TesterGuiIsClosingDown:                                             common_config.TesterGuiIsClosingDownStruct{},
		GuiExecutionServerIsClosingDown:                                    tempGuiExecutionServerIsClosingDown,
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{},
		GuiExecutionServerSendStartedUpTimeStamp:                           common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct{},
		AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination: common_config.AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombinationStruct{},
	}

	// Broadcast message via Channel on Broadcast system
	err = broadcastSenderForGuiExecutionServersInternalCommunicationChannel(&broadcastMessageForGuiExecutionServersInternalCommunicationChannel)

	return err
}

// Broadcast message to all other GuiExecutionServers that 'ThisGuiExecutionServerIsStartingUp'
func broadcastSenderForChannelMessage_ThisGuiExecutionServerIsStartingUp(
	tempGuiExecutionServerIsStartingUp common_config.GuiExecutionServerIsStartingUpStruct) (
	err error) {

	// Convert into Broadcast message type
	var broadcastMessageForGuiExecutionServersInternalCommunicationChannel BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct
	broadcastMessageForGuiExecutionServersInternalCommunicationChannel = BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct{
		GuiExecutionServersInternalCommunicationChannelType:                GuiExecutionServerIsStartingUpMessage,
		TesterGuiIsClosingDown:                                             common_config.TesterGuiIsClosingDownStruct{},
		GuiExecutionServerIsClosingDown:                                    common_config.GuiExecutionServerIsClosingDownStruct{},
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{},
		GuiExecutionServerIsStartingUp:                                     tempGuiExecutionServerIsStartingUp,
		GuiExecutionServerSendStartedUpTimeStamp:                           common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct{},
		AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination: common_config.AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombinationStruct{},
	}

	// Broadcast message via Channel  on Broadcast system
	err = broadcastSenderForGuiExecutionServersInternalCommunicationChannel(&broadcastMessageForGuiExecutionServersInternalCommunicationChannel)

	return err
}

// Broadcast message to all other GuiExecutionServers about 'ThisGuiExecutionServerSendsStartedUpTimeStamp'
// This is done periodically to ensure that all other GuiExecutionServers have the same "world view"
func broadcastSenderForChannelMessage_ThisGuiExecutionServerSendsStartedUpTimeStamp(
	tempGuiExecutionServerStartedUpTimeStampRefresher common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct) (
	err error) {

	// Convert into Broadcast message type
	var broadcastMessageForGuiExecutionServersInternalCommunicationChannel BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct
	broadcastMessageForGuiExecutionServersInternalCommunicationChannel = BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct{
		GuiExecutionServersInternalCommunicationChannelType:                GuiExecutionServerSendsStartedUpTimeStampMessage,
		TesterGuiIsClosingDown:                                             common_config.TesterGuiIsClosingDownStruct{},
		GuiExecutionServerIsClosingDown:                                    common_config.GuiExecutionServerIsClosingDownStruct{},
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{},
		GuiExecutionServerIsStartingUp:                                     common_config.GuiExecutionServerIsStartingUpStruct{},
		GuiExecutionServerSendStartedUpTimeStamp:                           tempGuiExecutionServerStartedUpTimeStampRefresher,
		AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination: common_config.AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombinationStruct{},
	}

	// Broadcast message via Channel  on Broadcast system
	err = broadcastSenderForGuiExecutionServersInternalCommunicationChannel(&broadcastMessageForGuiExecutionServersInternalCommunicationChannel)

	return err
}

// Broadcast message to all other GuiExecutionServers that (this) 'AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination'
func broadcastSenderForChannelMessage_AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination(
	tempAnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination common_config.AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombinationStruct) (
	err error) {

	// Convert into Broadcast message type
	var broadcastMessageForGuiExecutionServersInternalCommunicationChannel BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct
	broadcastMessageForGuiExecutionServersInternalCommunicationChannel = BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct{
		GuiExecutionServersInternalCommunicationChannelType:                AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination,
		TesterGuiIsClosingDown:                                             common_config.TesterGuiIsClosingDownStruct{},
		GuiExecutionServerIsClosingDown:                                    common_config.GuiExecutionServerIsClosingDownStruct{},
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{},
		GuiExecutionServerSendStartedUpTimeStamp:                           common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct{},
		AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination: tempAnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination,
	}

	// Broadcast message via Channel  on Broadcast system
	err = broadcastSenderForGuiExecutionServersInternalCommunicationChannel(&broadcastMessageForGuiExecutionServersInternalCommunicationChannel)

	return err
}

// Broadcast message on 'GuiExecutionServersInternalCommunicationChannel'
func broadcastSenderForGuiExecutionServersInternalCommunicationChannel(
	broadcastMessageForGuiExecutionServersInternalCommunicationChannel *BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct) (
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"id": "a356d4dc-2672-4026-a6d8-539a47e1b564",
		"guiExecutionServersInternalCommunicationChannelTypeDescription": guiExecutionServersInternalCommunicationChannelTypeDescription[broadcastMessageForGuiExecutionServersInternalCommunicationChannel.GuiExecutionServersInternalCommunicationChannelType],
	}).Debug("Broadcasting this message type")

	// Create json as string
	var broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSlice []byte
	var broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSliceAsString string
	broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSlice, err = json.Marshal(broadcastMessageForGuiExecutionServersInternalCommunicationChannel)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":  "92f9e3d3-ed73-482a-a6ce-5c03f08a00ff",
			"err": err,
		}).Error("Couldn't convert into byte slice ")

		return err
	}

	// Convert byte slice into string
	broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSliceAsString =
		string(broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSlice)

	conn, err := fenixSyncShared.DbPool.Acquire(context.Background())
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":  "cf404e66-7b97-424b-b377-776b28adbf7f",
			"err": err.Error(),
		}).Error("Error when acquiring sql-connection for 'GuiExecutionServersInternalCommunicationChannel'")

		return err
	}
	defer conn.Release()

	_, err = fenixSyncShared.DbPool.Exec(context.Background(),
		"SELECT pg_notify('testerGuiOwnerEngineChannel2', $1)",
		broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSlice)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id": "d7c5685c-f19b-4885-83b8-bba6b0408ec9",
			"broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSliceAsString": broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSliceAsString,
			"err": err.Error(),
		}).Error("Error sending 'broadcastMessageForGuiExecutionServersInternalCommunicationChannel'")

		return err
	}

	common_config.Logger.WithFields(logrus.Fields{
		"id": "b87e98c9-a9a3-4d2d-8c6a-e6067b4f31c5",
		"broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSliceAsString": broadcastMessageForGuiExecutionServersInternalCommunicationChannelAsByteSliceAsString,
	}).Debug("Message sent over Broadcast system, on Channel ")

	return err
}
