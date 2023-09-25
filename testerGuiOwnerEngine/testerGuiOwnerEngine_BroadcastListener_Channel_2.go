package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"encoding/json"
	"errors"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// InitiateAndStartBroadcastChannel2ListenerEngine
// Start listen for Broadcasts regarding Channel 1
// 'MessageForSomeoneIsClosingDown'
func InitiateAndStartBroadcastChannel2ListenerEngine() {

	go func() {
		for {
			err := BroadcastListenerChannel2()
			if err != nil {

				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "7010f54a-61ee-4274-aabb-4ad017a725d9",
					"err": err,
				}).Error("Error return from Broadcast listener for Channel 1. Will retry in 5 seconds")
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func BroadcastListenerChannel2() error {

	var err error
	var broadcastMesageForPostgresChannel2Message BroadcastMesageForPostgresChannel2MessageStruct

	if fenixSyncShared.DbPool == nil {
		return errors.New("empty pool reference")
	}

	conn, err := fenixSyncShared.DbPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "LISTEN testerGuiOwnerEngineChannel2")
	if err != nil {
		return err
	}

	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "fed10f28-01de-4025-92ca-9a7947eaaff0",
				"err": err,
			}).Error("Error waiting for notification from 'testerGuiOwnerEngineChannel2'")

			return err
		}

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                        "a190a5b2-983e-4217-b32e-107538a026f0",
			"accepted message from pid": notification.PID,
			"channel":                   notification.Channel,
			"payload":                   notification.Payload,
		}).Debug("Got Broadcast message from Postgres Database, on Channel 2")

		err = json.Unmarshal([]byte(notification.Payload), &broadcastMesageForPostgresChannel2Message)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "e30c85ba-e41f-4946-a289-f2f77e68a5de",
				"err": err,
			}).Error("Got some error when Unmarshal incoming json in, Channel 2, over Broadcast system")
		} else {

			// Which message type was sent
			// 'ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination' or
			// 'UserUnsubscribesToUserAndTestCaseExecutionCombination'
			if broadcastMesageForPostgresChannel2Message.PostgresChannel2MessageMessageType ==
				ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationMessage {
				// A call from TesterGui to a GuiExecutionServer was done for TestGui to receive ExecutionStatusUpdates

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
					GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {
					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message
					var tempGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination common_config.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct
					tempGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination = common_config.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct{
						TesterGuiApplicationId: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							TesterGuiApplicationId,
						UserId: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							UserId,
						GuiExecutionServerApplicationId: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							GuiExecutionServerApplicationId,
						TestCaseExecutionUuid: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							TestCaseExecutionUuid,
						TestCaseExecutionVersion: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							TestCaseExecutionVersion,
						MessageTimeStamp: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							MessageTimeStamp,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand: common_config.ChannelCommand_AnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination,
						SomeoneIsClosingDown:               nil,
						ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: &tempGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination,
						UserUnsubscribesToUserAndTestCaseExecutionCombination:              nil,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

				}

			} else {
				// A call from TesterGui saying that it unSubscribes to ExecutionStatusUpdates

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
					GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message
					var tempUserUnsubscribesToUserAndTestCaseExecutionCombinationStruct common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
					tempUserUnsubscribesToUserAndTestCaseExecutionCombinationStruct = common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{
						TesterGuiApplicationId: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							TesterGuiApplicationId,
						UserId: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							UserId,
						GuiExecutionServerApplicationId: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							GuiExecutionServerApplicationId,
						TestCaseExecutionUuid: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							TestCaseExecutionUuid,
						TestCaseExecutionVersion: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							TestCaseExecutionVersion,
						MessageTimeStamp: broadcastMesageForPostgresChannel2Message.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination.
							MessageTimeStamp,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand: common_config.ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination,
						SomeoneIsClosingDown:               nil,
						ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination:              &tempUserUnsubscribesToUserAndTestCaseExecutionCombinationStruct,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand
				}

			}

		}
	}
}
