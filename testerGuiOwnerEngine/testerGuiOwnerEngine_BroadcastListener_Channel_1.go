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

// InitiateAndStartBroadcastNotifyEngine
// Start listen for Broadcasts regarding Channel 1
// MessageForSomeoneIsClosingDown
func InitiateAndStartBroadcastChannel1ListenerEngine() {

	go func() {
		for {
			err := BroadcastListener()
			if err != nil {

				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "b66878ed-52ee-4ca8-afa1-d3b3ee6edf51",
					"err": err,
				}).Error("Error return from Broadcast listener for Channel 1. Will retry in 5 seconds")
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func BroadcastListener() error {

	var err error
	var broadcastMessageForSomeoneIsClosingDown BroadcastMessageForSomeoneIsClosingDownStruct

	if fenixSyncShared.DbPool == nil {
		return errors.New("empty pool reference")
	}

	conn, err := fenixSyncShared.DbPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "LISTEN testerGuiOwnerEngineChannel1")
	if err != nil {
		return err
	}

	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "c417ecd3-99bf-422f-9af0-f5c5f94d889f",
				"err": err,
			}).Error("Error waiting for notification from 'testerGuiOwnerEngineChannel1'")

			return err
		}

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                        "5d9190b1-d904-45c0-a93f-edf9e3947e52",
			"accepted message from pid": notification.PID,
			"channel":                   notification.Channel,
			"payload":                   notification.Payload,
		}).Debug("Got Broadcast message from Postgres Database")

		err = json.Unmarshal([]byte(notification.Payload), &broadcastMessageForSomeoneIsClosingDown)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "28232706-01e4-4b95-98fc-cf729c2c1930",
				"err": err,
			}).Error("Got some error when Unmarshal incoming json in, Channel 1, over Broadcast system")
		} else {

			// Who is closing down
			if broadcastMessageForSomeoneIsClosingDown.WhoISClosingDown == common_config.GuiExecutionServer {
				// A 'GuiExecutionServer' is closing Down

				// Which 'GuiExecutionServer' is closing Down
				if broadcastMessageForSomeoneIsClosingDown.ApplicationId == common_config.ApplicationRunTimeUuid {
					// This 'GuiExecutionServer' is closing Down

					// Do nothing

				} else {
					// Other 'GuiExecutionServer' is closing Down

					// Convert message into channel-version of message
					var tempSomeoneIsClosingDown common_config.SomeoneIsClosingDownStruct
					tempSomeoneIsClosingDown = common_config.SomeoneIsClosingDownStruct{
						WhoISClosingDown: broadcastMessageForSomeoneIsClosingDown.WhoISClosingDown,
						ApplicationId:    broadcastMessageForSomeoneIsClosingDown.ApplicationId,
						UserId:           broadcastMessageForSomeoneIsClosingDown.UserId,
						MessageTimeStamp: broadcastMessageForSomeoneIsClosingDown.MessageTimeStamp,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand: common_config.ChannelCommand_AnotherGuiExecutionServerIsClosingDown,
						SomeoneIsClosingDown:               &tempSomeoneIsClosingDown,
						ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination:              nil,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

				}

			} else {
				// A 'TesterGui' is closing down

				// Was call originated from this 'GuiExecutionServer'
				if broadcastMessageForSomeoneIsClosingDown.ApplicationId == common_config.ApplicationRunTimeUuid {
					// Call was originated from this 'GuiExecutionServer'

					// Do nothing

				} else {
					// Call was originated from other 'GuiExecutionServer'

					// Convert message into channel-version of message
					var tempSomeoneIsClosingDown common_config.SomeoneIsClosingDownStruct
					tempSomeoneIsClosingDown = common_config.SomeoneIsClosingDownStruct{
						WhoISClosingDown: broadcastMessageForSomeoneIsClosingDown.WhoISClosingDown,
						ApplicationId:    broadcastMessageForSomeoneIsClosingDown.ApplicationId,
						UserId:           broadcastMessageForSomeoneIsClosingDown.UserId,
						MessageTimeStamp: broadcastMessageForSomeoneIsClosingDown.MessageTimeStamp,
					}

					// Put message on 'testGuiExecutionEngineChannel' to be processed
					var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
					testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
						TesterGuiOwnerEngineChannelCommand: common_config.ChannelCommand_UserIsClosingDown,
						SomeoneIsClosingDown:               &tempSomeoneIsClosingDown,
						ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: nil,
						UserUnsubscribesToUserAndTestCaseExecutionCombination:              nil,
					}

					// Put on EngineChannel
					common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

				}

			}

		}
	}
}
