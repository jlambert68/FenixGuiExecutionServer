package broadcastingEngine

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

type BroadcastingMessageForExecutionsStruct struct {
	BroadcastTimeStamp        string                           `json:"timestamp"`
	TestCaseExecutions        []TestCaseExecutionStruct        `json:"testcaseexecutions"`
	TestInstructionExecutions []TestInstructionExecutionStruct `json:"testinstructionexecutions"`
}

type TestCaseExecutionStruct struct {
	TestCaseExecutionUuid   string `json:"testcaseexecutionuuid"`
	TestCaseExecutionStatus string `json:"testcaseexecutionstatus"`
}

type TestInstructionExecutionStruct struct {
	TestInstructionExecutionUuid   string `json:"testinstructionuuid"`
	TestInstructionExecutionStatus string `json:"testinstructionstatus"`
}

// Start listen for Broadcasts regarding change in status TestCaseExecutions and TestInstructionExecutions
func InitiateAndStartBroadcastNotifyEngine() {

	go func() {
		err := Listen()
		if err != nil {
			log.Println("unable start listener:", err)
			os.Exit(1)
		}
	}()
}

func Listen() error {

	var err error
	var broadcastingMessageForExecutions BroadcastingMessageForExecutionsStruct

	if fenixSyncShared.DbPool == nil {
		return errors.New("empty pool reference")
	}

	conn, err := fenixSyncShared.DbPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "LISTEN notes")
	if err != nil {
		return err
	}

	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "78d8f31c-5323-4c73-8a6a-f6cfef66f649",
				"err": err,
			}).Error("Error waiting for notification")
		}

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                        "12874bd6-0868-4efd-b232-45624d29c3e5",
			"accepted message from pid": notification.PID,
			"channel":                   notification.Channel,
			"payload":                   notification.Payload,
		}).Debug("Got Broadcast message from Postgres Databas")

		err = json.Unmarshal([]byte(notification.Payload), &broadcastingMessageForExecutions)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "d359ea3c-d46e-4bd6-9a2c-df73a9509cd7",
				"err": err,
			}).Error("Got some error when Unmarshal incoming json over Broadcast system")
		}

		fmt.Println(broadcastingMessageForExecutions)
	}
}
