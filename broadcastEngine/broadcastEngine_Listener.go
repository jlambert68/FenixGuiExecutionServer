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

// InitiateAndStartBroadcastNotifyEngine
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

			// Break down 'broadcastingMessageForExecutions' and send correct content to correct sSubscribers.
			convertToChannelMessageAndPutOnChannels(broadcastingMessageForExecutions)

		}
	}
}

// Break down 'broadcastingMessageForExecutions' and send correct content to correct sSubscribers.
func convertToChannelMessageAndPutOnChannels(broadcastingMessageForExecutions BroadcastingMessageForExecutionsStruct) {

	// Create a map with all 'TestCaseExecutionUuid' + 'TestCaseExecutionVersion' combinations
	// Used to know which combinations that exists
	var mapKeysMap map[string][]string // map['TestCaseExecutionUuid' + 'TestCaseExecutionVersion'][]'TestCaseExecutionUuid' + 'TestCaseExecutionVersion' + indicator('TC' or 'TI')]'
	mapKeysMap = make(map[string][]string)
	var mapKeysMapKeyValue string
	var existInMap bool

	// Create map for messages, grouped by Subscription-parameter-key('TestCaseExecutionUuid'+'TestCaseExecutionVersion') to be sent over MessageChannel to be forwarded to TestGui
	var testCaseExecutionsStatusForChannelMessageMap map[string][]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage
	testCaseExecutionsStatusForChannelMessageMap = make(map[string][]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage)
	var mapKey string

	// Create ChannelMessages for TestCaseExecutions
	for _, testCaseExecutionFromBroadcastMessage := range broadcastingMessageForExecutions.TestCaseExecutions {

		testCaseExecutionVersionAsInteger, err := strconv.Atoi(testCaseExecutionFromBroadcastMessage.TestCaseExecutionVersion)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "b91af162-4fc7-416b-8681-ea101cb5ebd5",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.TestCaseExecutionVersion": testCaseExecutionFromBroadcastMessage.TestCaseExecutionVersion,
			}).Error("Couldn't convert 'TestCaseExecutionVersion' from Broadcast-message into an integer")

		} else {

			var testCaseExecutionStatusForChannelMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage
			testCaseExecutionStatusForChannelMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage{
				TestCaseExecutionUuid:    testCaseExecutionFromBroadcastMessage.TestCaseExecutionUuid,
				TestCaseExecutionVersion: int32(testCaseExecutionVersionAsInteger),
				TestCaseExecutionStatus:  fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum(fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_value[testCaseExecutionFromBroadcastMessage.TestCaseExecutionStatus]),
			}

			// Create mapKey consisting of 'TestCaseExecutionUuid' + 'TestCaseExecutionVersion'
			mapKey = testCaseExecutionFromBroadcastMessage.TestCaseExecutionUuid + testCaseExecutionFromBroadcastMessage.TestCaseExecutionVersion

			// Extract slice holding the status messages for TestCaseExecutions
			var tempTestCaseExecutionsStatusForChannelMessage []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage
			tempTestCaseExecutionsStatusForChannelMessage, existInMap = testCaseExecutionsStatusForChannelMessageMap[mapKey]

			if existInMap == false {
				// Add to 'mapKeyMap' that a new combination of 'TestCaseExecutionUuid' + 'TestCaseExecutionVersion'  for TestCaseExecutions was found for
				mapKeysMapKeyValue = mapKey + "TC"

				var mapKeysMapKeyValues []string
				mapKeysMapKeyValues, _ = mapKeysMap[mapKey]
				mapKeysMapKeyValues = append(mapKeysMapKeyValues, mapKeysMapKeyValue)

				mapKeysMap[mapKey] = mapKeysMapKeyValues

			}

			// Add new status message to slice and add slice back to map
			tempTestCaseExecutionsStatusForChannelMessage = append(tempTestCaseExecutionsStatusForChannelMessage, testCaseExecutionStatusForChannelMessage)
			testCaseExecutionsStatusForChannelMessageMap[mapKey] = tempTestCaseExecutionsStatusForChannelMessage
		}
	}

	// Create map for messages, grouped by Subscription-parameter-key('TestInstructionExecutionUuid'+'TestInstructionExecutionVersion') to be sent over MessageChannel to be forwarded to TestGui
	var testInstructionExecutionsStatusForChannelMessageMap map[string][]*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage
	testInstructionExecutionsStatusForChannelMessageMap = make(map[string][]*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage)

	// Create ChannelMessages for TestInstructionExecutions
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

		// Both integer conversions needs to be OK, 'TestCaseExecutionVersion' and 'TestInstructionExecutionVersion'
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

			// Create mapKey consisting of 'TestCaseExecutionUuid' + 'TestCaseExecutionVersion'
			mapKey = testInstructionExecutionFromBroadcastMessage.TestCaseExecutionUuid + testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion

			// Extract slice holding the status messages for TestInstructionExecutions
			var tempTestInstructionExecutionsStatusForChannelMessage []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage
			tempTestInstructionExecutionsStatusForChannelMessage, existInMap = testInstructionExecutionsStatusForChannelMessageMap[mapKey]

			if existInMap == false {
				// Add to 'mapKeyMap' that a new combination of 'TestCaseExecutionUuid' + 'TestCaseExecutionVersion'  for TestInstructionExecutions was found for
				mapKeysMapKeyValue = mapKey + "TI"

				var mapKeysMapKeyValues []string
				mapKeysMapKeyValues, _ = mapKeysMap[mapKey]
				mapKeysMapKeyValues = append(mapKeysMapKeyValues, mapKeysMapKeyValue)

				mapKeysMap[mapKey] = mapKeysMapKeyValues

			}

			// Add new status message to slice and add slice back to map
			tempTestInstructionExecutionsStatusForChannelMessage = append(tempTestInstructionExecutionsStatusForChannelMessage, testInstructionExecutionStatusForChannelMessage)
			testInstructionExecutionsStatusForChannelMessageMap[mapKey] = tempTestInstructionExecutionsStatusForChannelMessage
		}

		// Get all keys from 'mapKeysMap' to find all combinations of 'TestCaseExecutionUuid' + 'TestCaseExecutionVersion'
		var testCaseExecutionUuidAndTestCaseExecutionVersionKeySlice []string
		var tempKey string
		testCaseExecutionUuidAndTestCaseExecutionVersionKeySlice = make([]string, 0, len(mapKeysMap))
		for _, tempKey = range mapKeysMap {
			testCaseExecutionUuidAndTestCaseExecutionVersionKeySlice = append(testCaseExecutionUuidAndTestCaseExecutionVersionKeySlice, tempKey)
		}

		// Loop slice of combinations of ('TestCaseExecutionUuid' + 'TestCaseExecutionVersion')
		var tempTestCaseExecutionUuidTestCaseExecutionVersion string
		var executionType string
		for _, tempTestCaseExecutionUuidTestCaseExecutionVersion = range testCaseExecutionUuidAndTestCaseExecutionVersionKeySlice {

			// Extract which TesterGuis that are subscribing to this 'TestCaseExecution(Version)'
			var messageToTesterGuiForwardChannels []*MessageToTesterGuiForwardChannelType
			messageToTesterGuiForwardChannels = whoIsSubscribingToTestCaseExecution(tempTestCaseExecutionUuidTestCaseExecutionVersion)

			// If there aren't any subscribers then continue to next 'TestCaseExecutionUuid+TestCaseExecutionVersion'
			if len(messageToTesterGuiForwardChannels) == 0 {
				continue
			}

			// extract info about there are TestCaseExecutions and/or TestInstructionExecutions
			// map['TestCaseExecutionUuid' + 'TestCaseExecutionVersion'][]'TestCaseExecutionUuid' + 'TestCaseExecutionVersion' + indicator('TC' or 'TI')]'
			var tempTestCaseExecutionUuidTestCaseExecutionVersionSlice []string
			tempTestCaseExecutionUuidTestCaseExecutionVersionSlice, existInMap = mapKeysMap[tempTestCaseExecutionUuidTestCaseExecutionVersion]

			// Create object to be sent over channel
			var testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage
			testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage{
				ProtoFileVersionUsedByClient:    fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				TestCaseExecutionsStatus:        nil,
				TestInstructionExecutionsStatus: nil,
			}

			// Loop the slice to extract if there are TestCaseExecutionsStatuses and/or TestInstructionExecutionsStatuses
			var tempExecutionType string
			for _, tempExecutionType = range tempTestCaseExecutionUuidTestCaseExecutionVersionSlice {

				// Extract if the execution is TestCaseExecution(TC) or a TestInstructionExecution(TI), the last 2 characters
				executionType = tempExecutionType[len(tempExecutionType)-2:]

				// Based on ExecutionType add correct data
				switch executionType {

				case "TC":
					// There are TestCaseExecutionStatus-data then add that data to object, to be sent over channel
					var tempTestCaseExecutionsStatusForChannelMessage []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage
					tempTestCaseExecutionsStatusForChannelMessage, existInMap = testCaseExecutionsStatusForChannelMessageMap[tempTestCaseExecutionUuidTestCaseExecutionVersion]

					if existInMap == false {
						common_config.Logger.WithFields(logrus.Fields{
							"Id": "35818e92-5d81-4e7a-abfb-1b49cb87d97b",
							"tempTestCaseExecutionUuidTestCaseExecutionVersion": tempTestCaseExecutionUuidTestCaseExecutionVersion,
						}).Error("Couldn't find 'TestCaseExecutionUuid+TestCaseExecutionVersion' in 'testCaseExecutionsStatusForChannelMessageMap'. Key is missing")

						break
					}

					testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage.TestCaseExecutionsStatus = tempTestCaseExecutionsStatusForChannelMessage

				case "TI":
					// There are TestInstructionExecutionStatus-data then add that data to object, to be sent over channel
					var tempTestInstructionExecutionsStatusForChannelMessage []*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage
					tempTestInstructionExecutionsStatusForChannelMessage, existInMap = testInstructionExecutionsStatusForChannelMessageMap[tempTestCaseExecutionUuidTestCaseExecutionVersion]
					if existInMap == false {
						common_config.Logger.WithFields(logrus.Fields{
							"Id": "60fc7e82-35d1-448a-a04d-101a509e9183",
							"tempTestCaseExecutionUuidTestCaseExecutionVersion": tempTestCaseExecutionUuidTestCaseExecutionVersion,
						}).Error("Couldn't find 'TestCaseExecutionUuid+TestCaseExecutionVersion' in 'testInstructionExecutionsStatusForChannelMessageMap'. Key is missing")

						break
					}

					testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage.TestInstructionExecutionsStatus = tempTestInstructionExecutionsStatusForChannelMessage

				default:
					common_config.Logger.WithFields(logrus.Fields{
						"Id":                "739809e1-f927-4500-9f79-be92416f0a3a",
						"executionType":     executionType,
						"tempExecutionType": tempExecutionType,
					}).Error("Execution type isn't any of TestCaseExecution(TC) or TestInstructionExecution(TI)")

					break
				}

			}

			// The 'subscribeToMessagesStreamResponse' that will be added into Channel message
			var subscribeToMessagesStreamResponse *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
			subscribeToMessagesStreamResponse = &fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse{
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				IsKeepAliveMessage:           false,
				ExecutionsStatus:             testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage,
			}

			// Create channel Message to be sent over channel, and later sent to TesterGui
			var messageToTestGuiForwardChannel MessageToTestGuiForwardChannelStruct
			messageToTestGuiForwardChannel = MessageToTestGuiForwardChannelStruct{
				SubscribeToMessagesStreamResponse: subscribeToMessagesStreamResponse,
				IsKeepAliveMessage:                false,
			}

			// Loop subscribers channels and put message on channels
			var messageToTesterGuiForwardChannel *MessageToTesterGuiForwardChannelType
			for _, messageToTesterGuiForwardChannel = range messageToTesterGuiForwardChannels {
				// Send Message over 'MessageChannel'
				*messageToTesterGuiForwardChannel <- messageToTestGuiForwardChannel
			}

		}
	}
}
