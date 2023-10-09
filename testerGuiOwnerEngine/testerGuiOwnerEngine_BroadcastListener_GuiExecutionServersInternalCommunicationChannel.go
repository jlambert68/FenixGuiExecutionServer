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

// InitiateAndStartBroadcastChannelListenerEngine
// Start listen for Broadcasts regarding Channel 1
// 'MessageForSomeoneIsClosingDown'
func InitiateAndStartBroadcastChannelListenerEngine() {

	go func() {
		for {
			err := BroadcastListener_GuiExecutionServersInternalCommunicationChannel()
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

func BroadcastListener_GuiExecutionServersInternalCommunicationChannel() error {

	var err error
	var broadcastMessageForGuiExecutionServersInternalCommunicationChannel BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct

	if fenixSyncShared.DbPool == nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "09901ce4-2d84-4282-983d-1f99ffe5bf91",
			"err": err,
		}).Error("empty pool reference")

		return errors.New("empty pool reference")
	}

	conn, err := fenixSyncShared.DbPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "LISTEN guiExecutionServersInternalCommunicationChannel")
	if err != nil {
		return err
	}

	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "fed10f28-01de-4025-92ca-9a7947eaaff0",
				"err": err,
			}).Error("Error waiting for notification from 'guiExecutionServersInternalCommunicationChannel'")

			return err
		}

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                        "a190a5b2-983e-4217-b32e-107538a026f0",
			"accepted message from pid": notification.PID,
			"channel":                   notification.Channel,
			"payload":                   notification.Payload,
		}).Debug("Got Broadcast message from Postgres Database, on 'guiExecutionServersInternalCommunicationChannel'")

		err = json.Unmarshal([]byte(notification.Payload), &broadcastMessageForGuiExecutionServersInternalCommunicationChannel)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "e30c85ba-e41f-4946-a289-f2f77e68a5de",
				"err": err,
			}).Error("Got some error when Unmarshal incoming json in, 'guiExecutionServersInternalCommunicationChannel', over Broadcast system")
		} else {

			// Which message type was sent
			switch broadcastMessageForGuiExecutionServersInternalCommunicationChannel.GuiExecutionServersInternalCommunicationChannelType {
			case TesterGuiIsClosingDownMessage:
				// A 'TesterGui' is closing down

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
					TesterGuiIsClosingDown.TesterGuiApplicationId == common_config.ApplicationRunTimeUuid {

					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message and Put message on
					// 'testGuiExecutionEngineChannel' to be processed

					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_AnotherGuiExecutionServersTesterGuiIsClosingDown,
						TesterGuiIsClosingDown:                                &broadcastMessageForGuiExecutionServersInternalCommunicationChannel.TesterGuiIsClosingDown,
						GuiExecutionServerIsClosingDown:                       nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
						GuiExecutionServerIsStartingUp:                        nil,
						GuiExecutionServerStartedUpTimeStampRefresher:         nil,
						UserSubscribesToUserAndTestCaseExecutionCombination:   nil,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

				}

			case GuiExecutionServerIsClosingDownMessage:
				// A 'GuiExecutionServer' is closing Down

				// Which 'GuiExecutionServer' is closing Down
				if broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
					GuiExecutionServerIsClosingDown.GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {
					// This 'GuiExecutionServer' is closing Down

					// Do nothing

				} else {
					// Other 'GuiExecutionServer' is closing Down

					// Convert message into channel-version of message
					var tempGuiExecutionServerResponsibilities []common_config.GuiExecutionServerResponsibilityStruct
					for _, tempGuiExecutionServerResponsibility := range broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
						GuiExecutionServerIsClosingDown.GuiExecutionServerResponsibilities {

						var guiExecutionServerResponsibility common_config.GuiExecutionServerResponsibilityStruct
						guiExecutionServerResponsibility = common_config.GuiExecutionServerResponsibilityStruct{
							TesterGuiApplicationId:   tempGuiExecutionServerResponsibility.TesterGuiApplicationId,
							UserId:                   tempGuiExecutionServerResponsibility.UserId,
							TestCaseExecutionUuid:    tempGuiExecutionServerResponsibility.TestCaseExecutionUuid,
							TestCaseExecutionVersion: tempGuiExecutionServerResponsibility.TestCaseExecutionVersion,
						}

						tempGuiExecutionServerResponsibilities = append(
							tempGuiExecutionServerResponsibilities, guiExecutionServerResponsibility)
					}

					var tempGuiExecutionServerIsClosingDown common_config.GuiExecutionServerIsClosingDownStruct
					tempGuiExecutionServerIsClosingDown = common_config.GuiExecutionServerIsClosingDownStruct{
						GuiExecutionServerApplicationId: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
							GuiExecutionServerIsClosingDown.GuiExecutionServerApplicationId,
						MessageTimeStamp: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
							GuiExecutionServerIsClosingDown.MessageTimeStamp,
						GuiExecutionServerResponsibilities: tempGuiExecutionServerResponsibilities,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_AnotherGuiExecutionServerIsClosingDown,
						TesterGuiIsClosingDown:                                nil,
						GuiExecutionServerIsClosingDown:                       &tempGuiExecutionServerIsClosingDown,
						UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
						GuiExecutionServerIsStartingUp:                        nil,
						GuiExecutionServerStartedUpTimeStampRefresher:         nil,
						UserSubscribesToUserAndTestCaseExecutionCombination:   nil,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand
				}

			case UserUnsubscribesToUserAndTestCaseExecutionCombinationMessage:
				// A call from TesterGui saying that it unSubscribes to ExecutionStatusUpdates

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
					UserUnsubscribesToUserAndTestCaseExecutionCombination.
					GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message
					var tempUserUnsubscribesToUserAndTestCaseExecutionCombinationStruct common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
					tempUserUnsubscribesToUserAndTestCaseExecutionCombinationStruct = common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct{
						TesterGuiApplicationId: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.UserUnsubscribesToUserAndTestCaseExecutionCombination.
							TesterGuiApplicationId,
						UserId: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.UserUnsubscribesToUserAndTestCaseExecutionCombination.
							UserId,
						GuiExecutionServerApplicationId: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.UserUnsubscribesToUserAndTestCaseExecutionCombination.
							GuiExecutionServerApplicationId,
						TestCaseExecutionUuid: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.UserUnsubscribesToUserAndTestCaseExecutionCombination.
							TestCaseExecutionUuid,
						TestCaseExecutionVersion: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.UserUnsubscribesToUserAndTestCaseExecutionCombination.
							TestCaseExecutionVersion,
						MessageTimeStamp: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.UserUnsubscribesToUserAndTestCaseExecutionCombination.
							MessageTimeStamp,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_AnotherGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination,
						TesterGuiIsClosingDown:                                nil,
						GuiExecutionServerIsClosingDown:                       nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination: &tempUserUnsubscribesToUserAndTestCaseExecutionCombinationStruct,
						GuiExecutionServerIsStartingUp:                        nil,
						GuiExecutionServerStartedUpTimeStampRefresher:         nil,
						UserSubscribesToUserAndTestCaseExecutionCombination:   nil,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand
				}

			case GuiExecutionServerIsStartingUpMessage:
				// A call from TesterGui saying that it is starting up

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
					GuiExecutionServerIsStartingUp.
					GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message
					var tempGuiExecutionServerIsStartingUp common_config.GuiExecutionServerIsStartingUpStruct
					tempGuiExecutionServerIsStartingUp = common_config.GuiExecutionServerIsStartingUpStruct{
						GuiExecutionServerApplicationId: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
							GuiExecutionServerIsStartingUp.GuiExecutionServerApplicationId,
						MessageTimeStamp: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
							GuiExecutionServerIsStartingUp.MessageTimeStamp,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_AnotherGuiExecutionServerIsStartingUp,
						TesterGuiIsClosingDown:                                nil,
						GuiExecutionServerIsClosingDown:                       nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
						GuiExecutionServerIsStartingUp:                        &tempGuiExecutionServerIsStartingUp,
						GuiExecutionServerStartedUpTimeStampRefresher:         nil,
						UserSubscribesToUserAndTestCaseExecutionCombination:   nil,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand
				}

			case GuiExecutionServerSendsStartedUpTimeStampMessage:
				// A periodic call from GuiExecutionServer informing other GuiExecutionServers of its StartUpTimeStamp
				// Use to secure that all GuiExecutionServers are in sync

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
					GuiExecutionServerSendStartedUpTimeStamp.
					GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message
					var tempGuiExecutionServerStartedUpTimeStampRefresher common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct
					tempGuiExecutionServerStartedUpTimeStampRefresher = common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct{
						GuiExecutionServerApplicationId: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
							GuiExecutionServerIsStartingUp.GuiExecutionServerApplicationId,
						MessageTimeStamp: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
							GuiExecutionServerIsStartingUp.MessageTimeStamp,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_AnotherGuiExecutionServerSendsStartedUpTimeStamp,
						TesterGuiIsClosingDown:                                nil,
						GuiExecutionServerIsClosingDown:                       nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
						GuiExecutionServerIsStartingUp:                        nil,
						GuiExecutionServerStartedUpTimeStampRefresher:         &tempGuiExecutionServerStartedUpTimeStampRefresher,
						UserSubscribesToUserAndTestCaseExecutionCombination:   nil,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand
				}

			case UserSubscribesToUserAndTestCaseExecutionCombinationMessage:
				// Other GuiExecutionServer has taken this combination of TestCaseExecution and TesterGui

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
					GuiExecutionServerSendStartedUpTimeStamp.
					GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message
					var tempUserSubscribesToUserAndTestCaseExecutionCombinationStruct common_config.
						UserSubscribesToUserAndTestCaseExecutionCombinationStruct
					tempUserSubscribesToUserAndTestCaseExecutionCombinationStruct = common_config.
						UserSubscribesToUserAndTestCaseExecutionCombinationStruct{
						GuiExecutionServerApplicationId: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
							GuiExecutionServerIsStartingUp.GuiExecutionServerApplicationId,
						MessageTimeStamp: broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
							GuiExecutionServerIsStartingUp.MessageTimeStamp,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand: common_config.
							ChannelCommand_AnotherGuiExecutionServerSendsStartedUpTimeStamp,
						TesterGuiIsClosingDown:                                nil,
						GuiExecutionServerIsClosingDown:                       nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
						GuiExecutionServerIsStartingUp:                        nil,
						GuiExecutionServerStartedUpTimeStampRefresher:         nil,
						UserSubscribesToUserAndTestCaseExecutionCombination:   &tempUserSubscribesToUserAndTestCaseExecutionCombinationStruct,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand
				}

			case TesterGuiIsStartingUpMessage:
				// A 'TesterGui' is starting up

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMessageForGuiExecutionServersInternalCommunicationChannel.
					TesterGuiIsStartingUp.TesterGuiApplicationId == common_config.ApplicationRunTimeUuid {

					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message and Put message on
					// 'testGuiExecutionEngineChannel' to be processed

					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand:                    common_config.ChannelCommand_AnotherGuiExecutionServersTesterGuiIsStartingUp,
						TesterGuiIsClosingDown:                                nil,
						GuiExecutionServerIsClosingDown:                       nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination: nil,
						GuiExecutionServerIsStartingUp:                        nil,
						GuiExecutionServerStartedUpTimeStampRefresher:         nil,
						UserSubscribesToUserAndTestCaseExecutionCombination:   nil,
						TesterGuiIsStartingUp:                                 &broadcastMessageForGuiExecutionServersInternalCommunicationChannel.TesterGuiIsStartingUp,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand
				}

			default:
				common_config.Logger.WithFields(logrus.Fields{
					"Id": "ddfcf5d8-6e59-4ab9-a03e-eb12c3f54106",
					"broadcastMessageForGuiExecutionServersInternalCommunicationChannel.GuiExecutionServersInternalCommunicationChannelType": broadcastMessageForGuiExecutionServersInternalCommunicationChannel.GuiExecutionServersInternalCommunicationChannelType,
				}).Fatal("Unhandled 'GuiExecutionServersInternalCommunicationChannelType'")

			}
		}
	}
}
