package broadcastEngine

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"encoding/json"
	"errors"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"time"
)

type BroadcastingMessageForExecutionsStruct struct {
	BroadcastTimeStamp        string                           `json:"timestamp"`
	TestCaseExecutions        []TestCaseExecutionStruct        `json:"testcaseexecutions"`
	TestInstructionExecutions []TestInstructionExecutionStruct `json:"testinstructionexecutions"`
}

type TestCaseExecutionStruct struct {
	TestCaseExecutionUuid    string `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion string `json:"testcaseexecutionversion"`
	TestCaseExecutionStatus  string `json:"testcaseexecutionstatus"`
}

type TestInstructionExecutionStruct struct {
	TestCaseExecutionUuid           string `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion        string `json:"testcaseexecutionversion"`
	TestInstructionExecutionUuid    string `json:"testinstructionexecutionuuid"`
	TestInstructionExecutionVersion string `json:"testinstructionexecutionversion"`
	TestInstructionExecutionStatus  string `json:"testinstructionexecutionstatus"`
}

// Start listen for Broadcasts regarding change in status TestCaseExecutions and TestInstructionExecutions
func InitiateAndStartBroadcastNotifyEngine() {

	go func() {
		for {
			err := BroadcastListener()
			if err != nil {
				log.Println("unable start listener:", err)

				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "c46d3d7c-3a13-4fe2-8633-d339a5f594db",
					"err": err,
				}).Error("Unable to start Broadcast listener. Will retry in 5 seconds")
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func BroadcastListener() error {

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
		} else {
			// Create message to be sent ove MessageChannel to be forwarded to TestGui
			var testCaseExecutionsStatusForChannelMessage []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage

			// Create ChannelMessage for TestCaseExecutions
			for _, testCaseExecutionFromBroadcasrMessage := range broadcastingMessageForExecutions.TestCaseExecutions {

				testCaseExecutionVersionAsInteger, err := strconv.Atoi(testCaseExecutionFromBroadcasrMessage.TestCaseExecutionVersion)
				if err != nil {
					common_config.Logger.WithFields(logrus.Fields{
						"Id":  "b91af162-4fc7-416b-8681-ea101cb5ebd5",
						"err": err,
						"testCaseExecutionFromBroadcasrMessage.TestCaseExecutionVersion": testCaseExecutionFromBroadcasrMessage.TestCaseExecutionVersion,
					}).Error("Couldn't convert 'TestCaseExecutionVersion' from Broadcast-message into an integer")

				} else {

					var testCaseExecutionStatusForChannelMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage
					testCaseExecutionStatusForChannelMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage{
						TestCaseExecutionUuid:    testCaseExecutionFromBroadcasrMessage.TestCaseExecutionUuid,
						TestCaseExecutionVersion: int32(testCaseExecutionVersionAsInteger),
						TestCaseExecutionStatus:  fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum(fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_value[testCaseExecutionFromBroadcasrMessage.TestCaseExecutionStatus]),
					}

					testCaseExecutionsStatusForChannelMessage = append(testCaseExecutionsStatusForChannelMessage, testCaseExecutionStatusForChannelMessage)
				}
			}

			// Create ChannelMessage for TestInstructionExecutions
			var testInstructionExecutionsStatusForChannelMessage []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage

			for _, testInstructionExecutionFromBroadcastMessage := range broadcastingMessageForExecutions.TestInstructionExecutions {

				testCaseExecutionVersionAsInteger, testCaseExecutionVersionError := strconv.Atoi(testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion)
				if testCaseExecutionVersionError != nil {
					common_config.Logger.WithFields(logrus.Fields{
						"Id":                            "da61719b-1444-4b35-ad55-22dd1d83f491",
						"testCaseExecutionVersionError": testCaseExecutionVersionError,
						"testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion": testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion,
					}).Error("Couldn't convert 'TestCaseExecutionVersion' from Broadcast-message into an integer")

				}

				testInstructionExecutionVersionAsInteger, testInstructionExecutionVersionError := strconv.Atoi(testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionVersion)
				if testInstructionExecutionVersionError != nil {
					common_config.Logger.WithFields(logrus.Fields{
						"Id":                                   "0d345833-2b64-4b1d-8433-6d9a7f2d88f6",
						"testInstructionExecutionVersionError": testInstructionExecutionVersionError,
						"testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion": testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionVersion,
					}).Error("Couldn't convert 'TestInstructionExecutionVersion' from Broadcast-message into an integer")

				}

				// Both integer conversions needs to be OK
				if testCaseExecutionVersionError == nil &&
					testInstructionExecutionVersionError == nil {

					var testInstructionExecutionStatusForChannelMessage *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage
					testInstructionExecutionStatusForChannelMessage = &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage{
						TestCaseExecutionUuid:           testInstructionExecutionFromBroadcastMessage.TestCaseExecutionUuid,
						TestCaseExecutionVersion:        int32(testCaseExecutionVersionAsInteger),
						TestInstructionExecutionUuid:    testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionUuid,
						TestInstructionExecutionVersion: int32(testInstructionExecutionVersionAsInteger),
						TestInstructionExecutionStatus:  fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum(fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum_value[testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionStatus]),
					}

					testInstructionExecutionsStatusForChannelMessage = append(testInstructionExecutionsStatusForChannelMessage, testInstructionExecutionStatusForChannelMessage)
				}
			}

			// Message holding TestCaseExecutions and TestInstructionExecutions, that will be added to 'subscribeToMessagesStreamResponse'-message
			var testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage
			testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage{
				ProtoFileVersionUsedByClient:    fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				TestCaseExecutionsStatus:        testCaseExecutionsStatusForChannelMessage,
				TestInstructionExecutionsStatus: testInstructionExecutionsStatusForChannelMessage,
			}

			// The 'subscribeToMessagesStreamResponse' that will be added into Channel message
			var subscribeToMessagesStreamResponse *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
			subscribeToMessagesStreamResponse = &fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse{
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				IsKeepAliveMessage:           false,
				ExecutionsStatus:             testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage,
			}

			// Channel Message to be sent over channel, and later sent to TesterGui
			var messageToTestGuiForwardChannelStruct MessageToTestGuiForwardChannelStruct
			messageToTestGuiForwardChannelStruct = MessageToTestGuiForwardChannelStruct{
				SubscribeToMessagesStreamResponse: subscribeToMessagesStreamResponse,
				IsKeepAliveMessage:                false,
			}

			// Send Message over 'MessageChannel'
			MessageToTesterGuiForwardChannel <- messageToTestGuiForwardChannelStruct

		}
	}
}
