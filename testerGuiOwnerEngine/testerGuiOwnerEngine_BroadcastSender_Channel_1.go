package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"encoding/json"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

func broadcastSenderForChannel1(tempSomeoneIsClosingDown *common_config.SomeoneIsClosingDownStruct) (err error) {

	// Convert into Broadcast message type
	var broadcastMessageForSomeoneIsClosingDown BroadcastMessageForSomeoneIsClosingDownStruct
	broadcastMessageForSomeoneIsClosingDown = BroadcastMessageForSomeoneIsClosingDownStruct{
		WhoISClosingDown: tempSomeoneIsClosingDown.WhoISClosingDown,
		ApplicationId:    tempSomeoneIsClosingDown.ApplicationId,
		UserId:           tempSomeoneIsClosingDown.UserId,
		MessageTimeStamp: tempSomeoneIsClosingDown.MessageTimeStamp,
	}

	// Create json as string
	var broadcastMessageForSomeoneIsClosingDownAsByteSlice []byte
	var broadcastMessageForSomeoneIsClosingDownAsByteSliceAsString string
	broadcastMessageForSomeoneIsClosingDownAsByteSlice, err = json.Marshal(broadcastMessageForSomeoneIsClosingDown)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":  "0254a586-51de-4276-a8bb-82310715c7c3",
			"err": err,
		}).Error("Couldn't convert into byte slice ")

		return err
	}

	// Convert byte slice into string
	broadcastMessageForSomeoneIsClosingDownAsByteSliceAsString = string(broadcastMessageForSomeoneIsClosingDownAsByteSlice)

	conn, err := fenixSyncShared.DbPool.Acquire(context.Background())
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":  "f29a19a8-d77b-4f9d-9185-9c0052b76de8",
			"err": err.Error(),
		}).Error("Error when acquiring sql-connection for Channel 1")

		return err
	}
	defer conn.Release()

	_, err = fenixSyncShared.DbPool.Exec(context.Background(),
		"SELECT pg_notify('testerGuiOwnerEngineChannel1', $1)", broadcastMessageForSomeoneIsClosingDownAsByteSlice)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id": "60e44727-5479-4664-8ca8-9aaa1220e06c",
			"broadcastMessageForSomeoneIsClosingDownAsByteSliceAsString": broadcastMessageForSomeoneIsClosingDownAsByteSliceAsString,
			"err": err.Error(),
		}).Error("Error sending 'broadcastMessageForSomeoneIsClosingDownAsByteSliceAsString' on Channel 1")

		return err
	}

	common_config.Logger.WithFields(logrus.Fields{
		"id": "e6d0b3d6-871a-4f16-aca1-bcbfde9ea92e",
		"broadcastMessageForSomeoneIsClosingDownAsByteSliceAsString": broadcastMessageForSomeoneIsClosingDownAsByteSliceAsString,
	}).Debug("Message sent over Broadcast system, on Channel 1")

	return err
}
