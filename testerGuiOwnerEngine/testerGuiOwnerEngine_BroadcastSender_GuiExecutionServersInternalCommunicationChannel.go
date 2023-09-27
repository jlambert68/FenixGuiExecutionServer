package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"encoding/json"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// Broadcast message to all other GuiExecutionServers that 'ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination'
func broadcastSenderForChannel2_ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination(
	tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct common_config.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct) (
	err error) {

	// Convert into Broadcast message type
	var tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct
	tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination = ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct{
		TesterGuiApplicationId:          tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct.TesterGuiApplicationId,
		UserId:                          tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct.UserId,
		GuiExecutionServerApplicationId: tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct.GuiExecutionServerApplicationId,
		TestCaseExecutionUuid:           tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct.TestCaseExecutionUuid,
		TestCaseExecutionVersion:        tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct.TestCaseExecutionVersion,
		MessageTimeStamp:                tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct.MessageTimeStamp,
	}
	var broadcastMessageForSomeoneIsClosingDown Channel2TakeOverUserAndTestCaseExecutionCombinationOrTesterGuiUnsubscribesStruct
	broadcastMessageForSomeoneIsClosingDown = Channel2TakeOverUserAndTestCaseExecutionCombinationOrTesterGuiUnsubscribesStruct{
		PostgresChannel2MessageMessageType:                                 ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationMessage,
		ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: tempThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination,
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{},
	}

	broadcastMessageForSomeoneIsClosingDown = Channel2TakeOverUserAndTestCaseExecutionCombinationOrTesterGuiUnsubscribesStruct{
		PostgresChannel2MessageMessageType:                                 0,
		ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct{},
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{},
	}

	// Broadcast message via Channel 2 on Broadcast system
	err = broadcastSenderForChannel2(&broadcastMessageForSomeoneIsClosingDown)

	return err
}

// Broadcast message to all other GuiExecutionServers that 'UserUnsubscribesToUserAndTestCaseExecutionCombinationMessage'
func broadcastSenderForChannel2_UserUnsubscribesToUserAndTestCaseExecutionCombinationMessage(
	userUnsubscribesToUserAndTestCaseExecutionCombination common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct) (
	err error) {

	// Convert into Broadcast message type
	var tempUserUnsubscribesToUserAndTestCaseExecutionCombination UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
	tempUserUnsubscribesToUserAndTestCaseExecutionCombination = UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{
		TesterGuiApplicationId:          userUnsubscribesToUserAndTestCaseExecutionCombination.TesterGuiApplicationId,
		UserId:                          userUnsubscribesToUserAndTestCaseExecutionCombination.UserId,
		GuiExecutionServerApplicationId: userUnsubscribesToUserAndTestCaseExecutionCombination.GuiExecutionServerApplicationId,
		TestCaseExecutionUuid:           userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionUuid,
		TestCaseExecutionVersion:        userUnsubscribesToUserAndTestCaseExecutionCombination.TestCaseExecutionVersion,
		MessageTimeStamp:                userUnsubscribesToUserAndTestCaseExecutionCombination.MessageTimeStamp,
	}

	var broadcastMessageForUnsubscribesToUserAndTestCaseExecutionCombination Channel2TakeOverUserAndTestCaseExecutionCombinationOrTesterGuiUnsubscribesStruct
	broadcastMessageForUnsubscribesToUserAndTestCaseExecutionCombination = Channel2TakeOverUserAndTestCaseExecutionCombinationOrTesterGuiUnsubscribesStruct{
		PostgresChannel2MessageMessageType:                                 UserUnsubscribesToUserAndTestCaseExecutionCombinationMessage,
		ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct{},
		UserUnsubscribesToUserAndTestCaseExecutionCombination:              tempUserUnsubscribesToUserAndTestCaseExecutionCombination,
	}

	// Broadcast message via Channel 2 on Broadcast system
	err = broadcastSenderForChannel2(&broadcastMessageForUnsubscribesToUserAndTestCaseExecutionCombination)

	return err
}

// Broadcast message on Channel 2
func broadcastSenderForChannel2(
	broadcastMesageForPostgresChannel2Message *Channel2TakeOverUserAndTestCaseExecutionCombinationOrTesterGuiUnsubscribesStruct) (
	err error) {

	// Create json as string
	var broadcastMessageForPostgresChannel2AsByteSlice []byte
	var broadcastMessageForPostgresChannel2AsByteSliceAsString string
	broadcastMessageForPostgresChannel2AsByteSlice, err = json.Marshal(broadcastMesageForPostgresChannel2Message)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":  "92f9e3d3-ed73-482a-a6ce-5c03f08a00ff",
			"err": err,
		}).Error("Couldn't convert into byte slice ")

		return err
	}

	// Convert byte slice into string
	broadcastMessageForPostgresChannel2AsByteSliceAsString = string(broadcastMessageForPostgresChannel2AsByteSlice)

	conn, err := fenixSyncShared.DbPool.Acquire(context.Background())
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":  "cf404e66-7b97-424b-b377-776b28adbf7f",
			"err": err.Error(),
		}).Error("Error when acquiring sql-connection for Channel 2")

		return err
	}
	defer conn.Release()

	_, err = fenixSyncShared.DbPool.Exec(context.Background(),
		"SELECT pg_notify('testerGuiOwnerEngineChannel2', $1)", broadcastMessageForPostgresChannel2AsByteSlice)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id": "d7c5685c-f19b-4885-83b8-bba6b0408ec9",
			"broadcastMessageForPostgresChannel2AsByteSliceAsString": broadcastMessageForPostgresChannel2AsByteSliceAsString,
			"err": err.Error(),
		}).Error("Error sending 'broadcastMesageForPostgresChannel2Message' on Channel 2")

		return err
	}

	common_config.Logger.WithFields(logrus.Fields{
		"id": "b87e98c9-a9a3-4d2d-8c6a-e6067b4f31c5",
		"broadcastMessageForPostgresChannel2AsByteSliceAsString": broadcastMessageForPostgresChannel2AsByteSliceAsString,
	}).Debug("Message sent over Broadcast system, on Channel 2")

	return err
}
