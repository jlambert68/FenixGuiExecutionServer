package broadcastEngine_ExecutionStatusUpdate

import (
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/testerGuiOwnerEngine"
	"context"
	"encoding/json"
	"errors"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixSyncShared/pubSubHelpers"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"strconv"
	"time"
)

type BroadcastingMessageForExecutionsStruct struct {
	OriginalMessageCreationTimeStamp string                                           `json:"originalmessagecreationtimestamp"`
	TestCaseExecutions               []TestCaseExecutionBroadcastMessageStruct        `json:"testcaseexecutions"`
	TestInstructionExecutions        []TestInstructionExecutionBroadcastMessageStruct `json:"testinstructionexecutions"`
}

type TestCaseExecutionBroadcastMessageStruct struct {
	TestCaseExecutionUuid          string `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion       string `json:"testcaseexecutionversion"`
	TestCaseExecutionStatus        string `json:"testcaseexecutionstatus"`
	ExecutionStartTimeStamp        string `json:"executionstarttimeStamp"`        // The timestamp when the execution was put for execution, not on queue for execution
	ExecutionStopTimeStamp         string `json:"executionstoptimestamp"`         // The timestamp when the execution was ended, in anyway
	ExecutionHasFinished           string `json:"executionhasfinished"`           // A simple status telling if the execution has ended or not
	ExecutionStatusUpdateTimeStamp string `json:"executionstatusupdatetimestamp"` // The timestamp when the status was last updated
	BroadcastTimeStamp             string `json:"broadcasttimestamp"`
	PreviousBroadcastTimeStamp     string `json:"previousbroadcasttimestamp"`
}

type TestInstructionExecutionBroadcastMessageStruct struct {
	TestCaseExecutionUuid                string `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion             string `json:"testcaseexecutionversion"`
	TestInstructionExecutionUuid         string `json:"testinstructionexecutionuuid"`
	TestInstructionExecutionVersion      string `json:"testinstructionexecutionversion"`
	SentTimeStamp                        string `json:"senttimestamp"`
	ExpectedExecutionEndTimeStamp        string `json:"expectedexecutionendtimestamp"`
	TestInstructionExecutionStatusName   string `json:"testinstructionexecutionstatusname"`
	TestInstructionExecutionStatusValue  string `json:"testinstructionexecutionstatusvalue"`
	TestInstructionExecutionEndTimeStamp string `json:"testinstructionexecutionendtimestamp"`
	TestInstructionExecutionHasFinished  string `json:"testinstructionexecutionhasfinished"`
	UniqueDatabaseRowCounter             string `json:"uniquedatabaserowcounter"`
	TestInstructionCanBeReExecuted       string `json:"testinstructioncanbereexecuted"`
	ExecutionStatusUpdateTimeStamp       string `json:"executionstatusupdatetimestamp"`
	BroadcastTimeStamp                   string `json:"broadcasttimestamp"`
	PreviousBroadcastTimeStamp           string `json:"previousbroadcasttimestamp"`
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

			// Restart broadcast engine when error occurs. Most probably because nothing is coming
			defer func() {
				_ = BroadcastListener()
			}()
			return err
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

	var originalMessageCreationTimeStamp time.Time
	var err error
	var timeStampLayoutForParser string //:= "2006-01-02 15:04:05.999999999 -0700 MST"

	// Convert Original Message creation Timestamp into time-variable
	timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(broadcastingMessageForExecutions.OriginalMessageCreationTimeStamp)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "d2cb561b-9976-407a-a263-a588529019f1",
			"err": err,
			"broadcastingMessageForExecutions.OriginalMessageCreationTimeStamp": broadcastingMessageForExecutions.OriginalMessageCreationTimeStamp,
		}).Error("Couldn't generate parser layout from TimeStamp")

		return
	}

	originalMessageCreationTimeStamp, err = time.Parse(timeStampLayoutForParser, broadcastingMessageForExecutions.OriginalMessageCreationTimeStamp)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                               "422159b0-de42-4b5d-a707-34dfabbf5082",
			"err":                              err,
			"broadcastingMessageForExecutions": broadcastingMessageForExecutions,
		}).Error("Couldn't parse TimeStamp in Broadcast-message")

		return
	}

	// Convert Original message creation Timestamp into gRPC-version
	var originalMessageCreationTimeStampForGrpc *timestamppb.Timestamp
	originalMessageCreationTimeStampForGrpc = timestamppb.New(originalMessageCreationTimeStamp)

	// Create map for messages, grouped by Subscription-parameter-key('TestCaseExecutionUuid'+'TestCaseExecutionVersion') to be sent over MessageChannel to be forwarded to TestGui
	var testCaseExecutionsStatusForChannelMessageMap map[string][]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage
	testCaseExecutionsStatusForChannelMessageMap = make(map[string][]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage)
	var mapKey string
	var testCaseExecutionVersionAsInteger int

	var tempExecutionStartTimeStamp time.Time
	var tempExecutionStopTimeStamp time.Time
	var tempExecutionStatusUpdateTimeStamp time.Time
	var tempExecutionHasFinished bool

	var broadcastTimeStamp time.Time
	var previousBroadcastTimeStamp time.Time

	// Create ChannelMessages for TestCaseExecutions
	var testCaseExecutionFromBroadcastMessage TestCaseExecutionBroadcastMessageStruct
	for _, testCaseExecutionFromBroadcastMessage = range broadcastingMessageForExecutions.TestCaseExecutions {

		// Convert Broadcast Timestamp into time-variable
		timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testCaseExecutionFromBroadcastMessage.BroadcastTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "dcc1f424-f375-4700-9b4c-129932676b98",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.BroadcastTimeStamp": testCaseExecutionFromBroadcastMessage.BroadcastTimeStamp,
			}).Error("Couldn't generate parser layout from TimeStamp")

			return
		}

		broadcastTimeStamp, err = time.Parse(timeStampLayoutForParser, testCaseExecutionFromBroadcastMessage.BroadcastTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                               "60b77ba4-3c99-4e39-b35a-33bed1c7155b",
				"err":                              err,
				"broadcastingMessageForExecutions": broadcastingMessageForExecutions,
			}).Error("Couldn't parse TimeStamp in Broadcast-message")

			return
		}

		// Convert Previous Broadcast Timestamp into time-variable
		timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testCaseExecutionFromBroadcastMessage.PreviousBroadcastTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "8b929496-cc27-4858-95f9-49e9012a28f6",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.PreviousBroadcastTimeStamp": testCaseExecutionFromBroadcastMessage.PreviousBroadcastTimeStamp,
			}).Error("Couldn't generate parser layout from TimeStamp")

			return
		}

		previousBroadcastTimeStamp, err = time.Parse(timeStampLayoutForParser, testCaseExecutionFromBroadcastMessage.PreviousBroadcastTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                               "acfac042-c396-4bd8-a12e-df35076e7ad8",
				"err":                              err,
				"broadcastingMessageForExecutions": broadcastingMessageForExecutions,
			}).Error("Couldn't parse Previous TimeStamp in Broadcast-message")

			return
		}

		// Convert string-versions from BroadcastMessage
		testCaseExecutionVersionAsInteger, err = strconv.Atoi(testCaseExecutionFromBroadcastMessage.TestCaseExecutionVersion)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "b91af162-4fc7-416b-8681-ea101cb5ebd5",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.TestCaseExecutionVersion": testCaseExecutionFromBroadcastMessage.TestCaseExecutionVersion,
			}).Error("Couldn't convert 'TestCaseExecutionVersion' from Broadcast-message into an integer")

			return

		}

		// Use fewer decimals for seconds in 'Layout' For TimeStamp-Parser
		//timeStampLayoutForParser = "2006-01-02 15:04:05.99999 -0700 MST"
		timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testCaseExecutionFromBroadcastMessage.ExecutionStartTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "301a5e6c-b63e-4678-8d98-206570241fee",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.ExecutionStartTimeStamp": testCaseExecutionFromBroadcastMessage.ExecutionStartTimeStamp,
			}).Error("Couldn't generate parser layout from TimeStamp")

			return
		}

		// Convert 'ExecutionStartTimeStamp'
		tempExecutionStartTimeStamp, err = time.Parse(timeStampLayoutForParser, testCaseExecutionFromBroadcastMessage.ExecutionStartTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "99c61a6a-8caf-4557-a5eb-01d354f69e90",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.ExecutionStartTimeStamp": testCaseExecutionFromBroadcastMessage.ExecutionStartTimeStamp,
			}).Error("Couldn't parse TimeStamp in Broadcast-message")

			return
		}

		// Convert 'ExecutionHasFinished'
		tempExecutionHasFinished, err = strconv.ParseBool(testCaseExecutionFromBroadcastMessage.ExecutionHasFinished)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "47dd1a78-316b-48be-a255-50a5818b8761",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.ExecutionHasFinished": testCaseExecutionFromBroadcastMessage.ExecutionHasFinished,
			}).Error("Couldn't parse TimeStamp in Broadcast-message")

			return
		}

		// Convert 'ExecutionStopTimeStamp' if 'ExecutionHasFinished'
		if tempExecutionHasFinished == true {
			timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testCaseExecutionFromBroadcastMessage.ExecutionStopTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "3d916f23-5e25-46cc-9329-32c019713db9",
					"err": err,
					"testCaseExecutionFromBroadcastMessage.ExecutionStopTimeStamp": testCaseExecutionFromBroadcastMessage.ExecutionStopTimeStamp,
				}).Error("Couldn't generate parser layout from TimeStamp")

				return
			}

			tempExecutionStopTimeStamp, err = time.Parse(timeStampLayoutForParser, testCaseExecutionFromBroadcastMessage.ExecutionStopTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "820005de-6f98-4aa8-a202-95759fcc07e2",
					"err": err,
					"testCaseExecutionFromBroadcastMessage.ExecutionStopTimeStamp": testCaseExecutionFromBroadcastMessage.ExecutionStopTimeStamp,
				}).Error("Couldn't parse TimeStamp in Broadcast-message")

				return
			}
		}

		// Convert 'ExecutionStatusUpdateTimeStamp'
		timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testCaseExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "e891bef9-1d83-4336-92d3-120d9b7594db",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp": testCaseExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp,
			}).Error("Couldn't generate parser layout from TimeStamp")

			return
		}

		tempExecutionStatusUpdateTimeStamp, err = time.Parse(timeStampLayoutForParser, testCaseExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "8ffa06b0-7396-4859-9d7a-461d9da153ce",
				"err": err,
				"testCaseExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp": testCaseExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp,
			}).Error("Couldn't parse TimeStamp in Broadcast-message")

			return
		}

		// Convert Broadcast Timestamps into gRPC-version
		var broadcastTimeStampForGrpc *timestamppb.Timestamp
		broadcastTimeStampForGrpc = timestamppb.New(broadcastTimeStamp)

		// Convert Previous Broadcast Timestamp into gRPC-version
		var previousBroadcastTimeStampForGrpc *timestamppb.Timestamp
		previousBroadcastTimeStampForGrpc = timestamppb.New(previousBroadcastTimeStamp)

		var testCaseExecutionStatusForChannelMessage *fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage
		testCaseExecutionStatusForChannelMessage = &fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusMessage{
			TestCaseExecutionUuid:      testCaseExecutionFromBroadcastMessage.TestCaseExecutionUuid,
			TestCaseExecutionVersion:   int32(testCaseExecutionVersionAsInteger),
			BroadcastTimeStamp:         broadcastTimeStampForGrpc,
			PreviousBroadcastTimeStamp: previousBroadcastTimeStampForGrpc,
			TestCaseExecutionDetails: &fenixExecutionServerGuiGrpcApi.TestCaseExecutionDetailsMessage{
				ExecutionStartTimeStamp: timestamppb.New(tempExecutionStartTimeStamp),

				TestCaseExecutionStatus:        fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum(fenixExecutionServerGuiGrpcApi.TestCaseExecutionStatusEnum_value[testCaseExecutionFromBroadcastMessage.TestCaseExecutionStatus]),
				ExecutionHasFinished:           tempExecutionHasFinished,
				ExecutionStatusUpdateTimeStamp: timestamppb.New(tempExecutionStatusUpdateTimeStamp),
			},
		}

		// Only add Stop time when TestCaseExecution is finished
		if tempExecutionHasFinished == true {
			testCaseExecutionStatusForChannelMessage.TestCaseExecutionDetails.ExecutionStopTimeStamp = timestamppb.New(tempExecutionStopTimeStamp)
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

	// Create map for messages, grouped by Subscription-parameter-key('TestInstructionExecutionUuid'+'TestInstructionExecutionVersion') to be sent over MessageChannel to be forwarded to TestGui
	var testInstructionExecutionsStatusForChannelMessageMap map[string][]*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage
	testInstructionExecutionsStatusForChannelMessageMap = make(map[string][]*fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage)
	var testCaseExecutionVersionError error
	var testInstructionExecutionFromBroadcastMessage TestInstructionExecutionBroadcastMessageStruct
	var testInstructionExecutionVersionError error
	var testInstructionExecutionVersionAsInteger int
	var sentTimeStampAsTime time.Time
	var expectedExecutionEndTimeStampAsTime time.Time
	var testInstructionExecutionStatusAsInteger int
	var testInstructionExecutionEndTimeStampAsTime time.Time
	var testInstructionExecutionHasFinishedAsBool bool
	var uniqueDatabaseRowCounterAsInteger int
	var testInstructionCanBeReExecutedAsBool bool
	var executionStatusUpdateTimeStampAsTime time.Time

	// Create ChannelMessages for TestInstructionExecutions
	for _, testInstructionExecutionFromBroadcastMessage = range broadcastingMessageForExecutions.TestInstructionExecutions {

		// Convert Broadcast Timestamp into time-variable
		timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testInstructionExecutionFromBroadcastMessage.BroadcastTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "c72a7a03-54fb-4ae6-bd8c-79aca9d60390",
				"err": err,
				"testInstructionExecutionFromBroadcastMessage.BroadcastTimeStamp": testInstructionExecutionFromBroadcastMessage.BroadcastTimeStamp,
			}).Error("Couldn't generate parser layout from TimeStamp")

			return
		}

		broadcastTimeStamp, err = time.Parse(timeStampLayoutForParser, testInstructionExecutionFromBroadcastMessage.BroadcastTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                               "556f2893-eced-436c-b2df-f517a08398c0",
				"err":                              err,
				"broadcastingMessageForExecutions": broadcastingMessageForExecutions,
			}).Error("Couldn't parse TimeStamp in Broadcast-message")

			return
		}

		previousBroadcastTimeStamp, err = time.Parse(timeStampLayoutForParser, testInstructionExecutionFromBroadcastMessage.PreviousBroadcastTimeStamp)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                               "a7b5932c-ebc6-4f34-a846-a2c5224a7e3c",
				"err":                              err,
				"broadcastingMessageForExecutions": broadcastingMessageForExecutions,
			}).Error("Couldn't parse Previous TimeStamp in Broadcast-message")

			return
		}

		// Convert string-versions from BroadcastMessage
		testCaseExecutionVersionAsInteger, err = strconv.Atoi(testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "d1500b59-b38e-4e0b-b01a-a07d2ade1de6",
				"err": err,
				"testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion": testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion,
			}).Error("Couldn't convert 'TestCaseExecutionVersion' from Broadcast-message into an integer")

			return

		}

		// Parse 'TestCaseExecutionVersion' from Broadcast-message
		testCaseExecutionVersionAsInteger, testCaseExecutionVersionError = strconv.Atoi(testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion)
		if testCaseExecutionVersionError != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                            "da61719b-1444-4b35-ad55-22dd1d83f491",
				"testCaseExecutionVersionError": testCaseExecutionVersionError,
				"testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion": testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion,
			}).Error("Couldn't convert 'TestCaseExecutionVersion' from Broadcast-message into an integer")

			return
		}

		// Parse 'TestInstructionExecutionVersion' from Broadcast-message
		testInstructionExecutionVersionAsInteger, testInstructionExecutionVersionError = strconv.Atoi(testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionVersion)
		if testInstructionExecutionVersionError != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                                   "0d345833-2b64-4b1d-8433-6d9a7f2d88f6",
				"testInstructionExecutionVersionError": testInstructionExecutionVersionError,
				"testInstructionExecutionFromBroadcastMessage.TestCaseExecutionVersion": testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionVersion,
			}).Error("Couldn't convert 'TestInstructionExecutionVersion' from Broadcast-message into an integer")

			return
		}

		// Parse 'SentTimeStamp' from Broadcast-message
		if testInstructionExecutionFromBroadcastMessage.SentTimeStamp != "" {
			timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testInstructionExecutionFromBroadcastMessage.SentTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "e3c2d0e5-ba54-4409-b4a8-384d9c252407",
					"err": err,
					"testInstructionExecutionFromBroadcastMessage.SentTimeStamp": testInstructionExecutionFromBroadcastMessage.SentTimeStamp,
				}).Error("Couldn't generate parser layout from TimeStamp")

				return
			}

			sentTimeStampAsTime, err = time.Parse(timeStampLayoutForParser, testInstructionExecutionFromBroadcastMessage.SentTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "93f3bedb-dc41-45ad-b523-bb0214f823ec",
					"err": err,
					"testInstructionExecutionFromBroadcastMessage.SentTimeStamp": testInstructionExecutionFromBroadcastMessage.SentTimeStamp,
				}).Error("Couldn't parse TimeStamp in Broadcast-message")

				return
			}
		}

		// Parse 'ExpectedExecutionEndTimeStamp' from Broadcast-message
		if testInstructionExecutionFromBroadcastMessage.ExpectedExecutionEndTimeStamp != "" {
			timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testInstructionExecutionFromBroadcastMessage.ExpectedExecutionEndTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "0309db1d-30b8-41d3-8a42-4ce4ef6c56fc",
					"err": err,
					"testInstructionExecutionFromBroadcastMessage.ExpectedExecutionEndTimeStamp": testInstructionExecutionFromBroadcastMessage.ExpectedExecutionEndTimeStamp,
				}).Error("Couldn't generate parser layout from TimeStamp")

				return
			}

			expectedExecutionEndTimeStampAsTime, err = time.Parse(timeStampLayoutForParser, testInstructionExecutionFromBroadcastMessage.ExpectedExecutionEndTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "a67badc1-14a8-4c91-9c15-b0636f9ff374",
					"err": err,
					"testInstructionExecutionFromBroadcastMessage.ExpectedExecutionEndTimeStamp": testInstructionExecutionFromBroadcastMessage.ExpectedExecutionEndTimeStamp,
				}).Error("Couldn't parse TimeStamp in Broadcast-message")

				return
			}
		}

		// Parse 'TestInstructionExecutionStatusValue' from Broadcast-message
		testInstructionExecutionStatusAsInteger, testCaseExecutionVersionError = strconv.Atoi(testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionStatusValue)
		if testCaseExecutionVersionError != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                            "db59ef02-8113-42f6-94e7-a4119eaa3e52",
				"testCaseExecutionVersionError": testCaseExecutionVersionError,
				"testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionStatusValue": testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionStatusValue,
			}).Error("Couldn't convert 'TestInstructionExecutionStatusValue' from Broadcast-message into an integer")

			return

		}

		// Parse 'TestInstructionExecutionEndTimeStamp' from Broadcast-message
		if testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionEndTimeStamp != "" {
			timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionEndTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "60cd129f-6148-4702-89d4-1c32cc884e6e",
					"err": err,
					"testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionEndTimeStamp": testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionEndTimeStamp,
				}).Error("Couldn't generate parser layout from TimeStamp")

				return
			}

			testInstructionExecutionEndTimeStampAsTime, err = time.Parse(timeStampLayoutForParser, testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionEndTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "a6673e3d-9dc2-4d36-bf7e-a604c0a86a4c",
					"err": err,
					"testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionEndTimeStamp": testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionEndTimeStamp,
				}).Error("Couldn't parse TimeStamp in Broadcast-message")

				return
			}
		}

		// Parse 'TestInstructionExecutionHasFinished' from Broadcast-message
		testInstructionExecutionHasFinishedAsBool, err = strconv.ParseBool(testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionHasFinished)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "70878467-53d7-4ea1-b27b-703ab50c80d0",
				"err": err,
				"testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionHasFinished": testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionHasFinished,
			}).Error("Couldn't parse Boolean in Broadcast-message")

			return
		}

		// Parse 'UniqueDatabaseRowCounter' from Broadcast-message
		uniqueDatabaseRowCounterAsInteger, testCaseExecutionVersionError = strconv.Atoi(testInstructionExecutionFromBroadcastMessage.UniqueDatabaseRowCounter)
		if testCaseExecutionVersionError != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                            "e4cd8f5d-2acc-430a-8034-35e1fee1a1dc",
				"testCaseExecutionVersionError": testCaseExecutionVersionError,
				"testInstructionExecutionFromBroadcastMessage.UniqueDatabaseRowCounter": testInstructionExecutionFromBroadcastMessage.UniqueDatabaseRowCounter,
			}).Error("Couldn't convert 'UniqueDatabaseRowCounter' from Broadcast-message into an integer")

			return
		}

		// Parse 'TestInstructionCanBeReExecuted' from Broadcast-message
		testInstructionCanBeReExecutedAsBool, err = strconv.ParseBool(testInstructionExecutionFromBroadcastMessage.TestInstructionCanBeReExecuted)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":  "f2d35949-acca-49af-80e7-66be68ae42fb",
				"err": err,
				"testInstructionExecutionFromBroadcastMessage.TestInstructionCanBeReExecuted": testInstructionExecutionFromBroadcastMessage.TestInstructionCanBeReExecuted,
			}).Error("Couldn't parse Boolean in Broadcast-message")

			return
		}

		// Parse 'ExecutionStatusUpdateTimeStamp' from Broadcast-message
		if testInstructionExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp != "" {
			timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(testInstructionExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "cc142492-a4be-4074-b70b-e2e8142bfdb0",
					"err": err,
					"testInstructionExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp": testInstructionExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp,
				}).Error("Couldn't generate parser layout from TimeStamp")

				return
			}

			executionStatusUpdateTimeStampAsTime, err = time.Parse(timeStampLayoutForParser, testInstructionExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "7f398611-d998-4733-8e57-d203b10437d9",
					"err": err,
					"testInstructionExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp": testInstructionExecutionFromBroadcastMessage.ExecutionStatusUpdateTimeStamp,
				}).Error("Couldn't parse TimeStamp in Broadcast-message")

				return
			}
		}

		// Convert Broadcast Timestamps into gRPC-version
		var broadcastTimeStampForGrpc *timestamppb.Timestamp
		broadcastTimeStampForGrpc = timestamppb.New(broadcastTimeStamp)

		// Convert Previous Broadcast Timestamp into gRPC-version
		var previousBroadcastTimeStampForGrpc *timestamppb.Timestamp
		previousBroadcastTimeStampForGrpc = timestamppb.New(previousBroadcastTimeStamp)

		// Build TestInstructionExecution-part of status update message
		var testInstructionExecutionStatusForChannelMessage *fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage
		testInstructionExecutionStatusForChannelMessage = &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusMessage{
			TestCaseExecutionUuid:           testInstructionExecutionFromBroadcastMessage.TestCaseExecutionUuid,
			TestCaseExecutionVersion:        int32(testCaseExecutionVersionAsInteger),
			TestInstructionExecutionUuid:    testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionUuid,
			TestInstructionExecutionVersion: int32(testInstructionExecutionVersionAsInteger),
			TestInstructionExecutionStatus: fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum(
				fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum_value[testInstructionExecutionFromBroadcastMessage.TestInstructionExecutionStatusValue]),
			BroadcastTimeStamp:         broadcastTimeStampForGrpc,
			PreviousBroadcastTimeStamp: previousBroadcastTimeStampForGrpc,
			TestInstructionExecutionsStatusInformation: &fenixExecutionServerGuiGrpcApi.TestInstructionExecutionsInformationMessage{
				SentTimeStamp:                        timestamppb.New(sentTimeStampAsTime),
				ExpectedExecutionEndTimeStamp:        timestamppb.New(expectedExecutionEndTimeStampAsTime),
				TestInstructionExecutionStatus:       fenixExecutionServerGuiGrpcApi.TestInstructionExecutionStatusEnum(testInstructionExecutionStatusAsInteger),
				TestInstructionExecutionEndTimeStamp: timestamppb.New(testInstructionExecutionEndTimeStampAsTime),
				TestInstructionExecutionHasFinished:  testInstructionExecutionHasFinishedAsBool,
				UniqueDatabaseRowCounter:             uint64(uniqueDatabaseRowCounterAsInteger),
				TestInstructionCanBeReExecuted:       testInstructionCanBeReExecutedAsBool,
				ExecutionStatusUpdateTimeStamp:       timestamppb.New(executionStatusUpdateTimeStampAsTime),
			},
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
	for tempKey, _ = range mapKeysMap {
		testCaseExecutionUuidAndTestCaseExecutionVersionKeySlice = append(testCaseExecutionUuidAndTestCaseExecutionVersionKeySlice, tempKey)
	}

	// Loop slice of combinations of ('TestCaseExecutionUuid' + 'TestCaseExecutionVersion')
	var tempTestCaseExecutionUuidTestCaseExecutionVersion string
	var executionType string
	for _, tempTestCaseExecutionUuidTestCaseExecutionVersion = range testCaseExecutionUuidAndTestCaseExecutionVersionKeySlice {

		/*
			// Extract which TesterGuis that are subscribing to this 'TestCaseExecution(Version)'
			var messageToTesterGuiForwardChannels []*MessageToTesterGuiForwardChannelType
			messageToTesterGuiForwardChannels = whoIsSubscribingToTestCaseExecution(tempTestCaseExecutionUuidTestCaseExecutionVersion)

			// If there aren't any subscribers then continue to next 'TestCaseExecutionUuid+TestCaseExecutionVersion'
			if len(messageToTesterGuiForwardChannels) == 0 {
				continue
			}
		*/

		// Extract which TesterGuis that are subscribing to this 'TestCaseExecution(Version)'
		var guiExecutionServerResponsibilities []common_config.GuiExecutionServerResponsibilityStruct
		guiExecutionServerResponsibilities = testerGuiOwnerEngine.
			ListAllTestCaseExecutionsSubscriptionsForExecutionsSubscriptionsMapKeyFromMap(tempTestCaseExecutionUuidTestCaseExecutionVersion)

		// extract info about if there are TestCaseExecutions and/or TestInstructionExecutions
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

		// Convert into PubSub-version 'TestCaseExecutionsStatus-slice'
		var pubSubTestCaseExecutionsStatusSlice []*fenixExecutionServerGuiGrpcApi.
			ExecutionStatusMessagesPubSubSchema_TestCaseExecutionStatusMessage
		for _, tempPubSubTestCaseExecutionsStatus := range testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage.
			GetTestCaseExecutionsStatus() {

			var pubSubTestCaseExecutionsStatus *fenixExecutionServerGuiGrpcApi.
				ExecutionStatusMessagesPubSubSchema_TestCaseExecutionStatusMessage
			pubSubTestCaseExecutionsStatus = &fenixExecutionServerGuiGrpcApi.
				ExecutionStatusMessagesPubSubSchema_TestCaseExecutionStatusMessage{
				TestCaseExecutionUuid:    tempPubSubTestCaseExecutionsStatus.GetTestCaseExecutionUuid(),
				TestCaseExecutionVersion: tempPubSubTestCaseExecutionsStatus.GetTestCaseExecutionVersion(),
				BroadcastTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
					Seconds: tempPubSubTestCaseExecutionsStatus.GetBroadcastTimeStamp().GetSeconds(),
					Nanos:   tempPubSubTestCaseExecutionsStatus.GetBroadcastTimeStamp().GetNanos(),
				},
				PreviousBroadcastTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
					Seconds: tempPubSubTestCaseExecutionsStatus.GetPreviousBroadcastTimeStamp().GetSeconds(),
					Nanos:   tempPubSubTestCaseExecutionsStatus.GetPreviousBroadcastTimeStamp().GetNanos(),
				},
				TestCaseExecutionDetails: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_TestCaseExecutionDetailsMessage{
					ExecutionStartTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
						Seconds: tempPubSubTestCaseExecutionsStatus.TestCaseExecutionDetails.GetExecutionStartTimeStamp().
							GetSeconds(),
						Nanos: tempPubSubTestCaseExecutionsStatus.TestCaseExecutionDetails.GetExecutionStartTimeStamp().
							GetNanos(),
					},
					ExecutionStopTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
						Seconds: tempPubSubTestCaseExecutionsStatus.TestCaseExecutionDetails.GetExecutionStopTimeStamp().
							GetSeconds(),
						Nanos: tempPubSubTestCaseExecutionsStatus.TestCaseExecutionDetails.GetExecutionStopTimeStamp().
							GetNanos(),
					},
					TestCaseExecutionStatus: fenixExecutionServerGuiGrpcApi.
						ExecutionStatusMessagesPubSubSchema_TestCaseExecutionStatusEnum(tempPubSubTestCaseExecutionsStatus.
							TestCaseExecutionDetails.GetTestCaseExecutionStatus()),
					ExecutionHasFinished: tempPubSubTestCaseExecutionsStatus.TestCaseExecutionDetails.
						GetExecutionHasFinished(),
					ExecutionStatusUpdateTimeStamp: &fenixExecutionServerGuiGrpcApi.
						ExecutionStatusMessagesPubSubSchema_Timestamp{
						Seconds: tempPubSubTestCaseExecutionsStatus.TestCaseExecutionDetails.
							GetExecutionStatusUpdateTimeStamp().GetSeconds(),
						Nanos: tempPubSubTestCaseExecutionsStatus.TestCaseExecutionDetails.
							GetExecutionStatusUpdateTimeStamp().GetNanos(),
					},
					UniqueDatabaseRowCounter: tempPubSubTestCaseExecutionsStatus.TestCaseExecutionDetails.
						GetUniqueDatabaseRowCounter(),
				},
			}

			pubSubTestCaseExecutionsStatusSlice = append(pubSubTestCaseExecutionsStatusSlice, pubSubTestCaseExecutionsStatus)
		}

		// Convert into PubSub-version 'TestInstructionExecutionsStatus-slice'
		var pubSubTestInstructionExecutionsStatusSlice []*fenixExecutionServerGuiGrpcApi.
			ExecutionStatusMessagesPubSubSchema_TestInstructionExecutionStatusMessage
		for _, tempPubSubTestInstructionExecutionsStatus := range testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage.
			GetTestInstructionExecutionsStatus() {

			var pubSubTestInstructionExecutionsStatus *fenixExecutionServerGuiGrpcApi.
				ExecutionStatusMessagesPubSubSchema_TestInstructionExecutionStatusMessage
			pubSubTestInstructionExecutionsStatus = &fenixExecutionServerGuiGrpcApi.
				ExecutionStatusMessagesPubSubSchema_TestInstructionExecutionStatusMessage{
				TestCaseExecutionUuid:           tempPubSubTestInstructionExecutionsStatus.GetTestCaseExecutionUuid(),
				TestCaseExecutionVersion:        tempPubSubTestInstructionExecutionsStatus.GetTestCaseExecutionVersion(),
				TestInstructionExecutionUuid:    tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionUuid(),
				TestInstructionExecutionVersion: tempPubSubTestInstructionExecutionsStatus.GetTestCaseExecutionVersion(),
				TestInstructionExecutionStatus: fenixExecutionServerGuiGrpcApi.
					ExecutionStatusMessagesPubSubSchema_TestInstructionExecutionStatusEnum(
						tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionStatus()),
				BroadcastTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
					Seconds: tempPubSubTestInstructionExecutionsStatus.GetBroadcastTimeStamp().GetSeconds(),
					Nanos:   tempPubSubTestInstructionExecutionsStatus.GetBroadcastTimeStamp().GetNanos(),
				},
				PreviousBroadcastTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
					Seconds: tempPubSubTestInstructionExecutionsStatus.GetPreviousBroadcastTimeStamp().GetSeconds(),
					Nanos:   tempPubSubTestInstructionExecutionsStatus.GetPreviousBroadcastTimeStamp().GetNanos(),
				},
				TestInstructionExecutionsStatusInformation: &fenixExecutionServerGuiGrpcApi.
					ExecutionStatusMessagesPubSubSchema_TestInstructionExecutionsInformationMessage{
					SentTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
						Seconds: tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionsStatusInformation().
							GetSentTimeStamp().GetSeconds(),
						Nanos: tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionsStatusInformation().
							GetSentTimeStamp().GetNanos(),
					},
					ExpectedExecutionEndTimeStamp: &fenixExecutionServerGuiGrpcApi.
						ExecutionStatusMessagesPubSubSchema_Timestamp{
						Seconds: tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionsStatusInformation().
							GetExpectedExecutionEndTimeStamp().GetSeconds(),
						Nanos: tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionsStatusInformation().
							GetExpectedExecutionEndTimeStamp().GetNanos(),
					},
					TestInstructionExecutionStatus: fenixExecutionServerGuiGrpcApi.
						ExecutionStatusMessagesPubSubSchema_TestInstructionExecutionStatusEnum(
							tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionsStatusInformation().
								GetTestInstructionExecutionStatus()),
					TestInstructionExecutionEndTimeStamp: &fenixExecutionServerGuiGrpcApi.
						ExecutionStatusMessagesPubSubSchema_Timestamp{
						Seconds: tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionsStatusInformation().
							GetTestInstructionExecutionEndTimeStamp().GetSeconds(),
						Nanos: tempPubSubTestInstructionExecutionsStatus.GetTestInstructionExecutionsStatusInformation().
							GetTestInstructionExecutionEndTimeStamp().GetNanos(),
					},
					TestInstructionExecutionHasFinished: tempPubSubTestInstructionExecutionsStatus.
						GetTestInstructionExecutionsStatusInformation().GetTestInstructionExecutionHasFinished(),
					UniqueDatabaseRowCounter: tempPubSubTestInstructionExecutionsStatus.
						GetTestInstructionExecutionsStatusInformation().GetUniqueDatabaseRowCounter(),
					TestInstructionCanBeReExecuted: tempPubSubTestInstructionExecutionsStatus.
						GetTestInstructionExecutionsStatusInformation().GetTestInstructionCanBeReExecuted(),
					ExecutionStatusUpdateTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
						Seconds: tempPubSubTestInstructionExecutionsStatus.
							GetTestInstructionExecutionsStatusInformation().GetExecutionStatusUpdateTimeStamp().GetSeconds(),
						Nanos: tempPubSubTestInstructionExecutionsStatus.
							GetTestInstructionExecutionsStatusInformation().GetExecutionStatusUpdateTimeStamp().GetNanos(),
					},
				},
			}

			pubSubTestInstructionExecutionsStatusSlice = append(pubSubTestInstructionExecutionsStatusSlice, pubSubTestInstructionExecutionsStatus)
		}

		// Build the full PubSub-ExecutionStatus-message
		var executionStatusMessagesPubSubMessage fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema
		executionStatusMessagesPubSubMessage = fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema{
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			OriginalMessageCreationTimeStamp: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_Timestamp{
				Seconds: originalMessageCreationTimeStampForGrpc.GetSeconds(),
				Nanos:   originalMessageCreationTimeStampForGrpc.GetNanos(),
			},
			IsKeepAliveMessage: false,
			ExecutionsStatus: &fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_TestCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage{
				ProtoFileVersionUsedByClient:    fenixExecutionServerGuiGrpcApi.ExecutionStatusMessagesPubSubSchema_CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				TestCaseExecutionsStatus:        pubSubTestCaseExecutionsStatusSlice,
				TestInstructionExecutionsStatus: pubSubTestInstructionExecutionsStatusSlice,
			},
		}
		// Convert gPubSub-ExecutionStatus-message into json-string
		var processTestInstructionExecutionRequestAsJsonString string
		processTestInstructionExecutionRequestAsJsonString = protojson.Format(&executionStatusMessagesPubSubMessage)

		// Loop PubSub-targets and send message over PubSub
		var topicId string
		var returnMessageAckNack bool
		var returnMessageString string
		var err error
		for _, guiExecutionServerResponsibility := range guiExecutionServerResponsibilities {
			topicId = testerGuiOwnerEngine.GeneratePubSubTopicForExecutionStatusUpdates(guiExecutionServerResponsibility.UserId)
			returnMessageAckNack, returnMessageString, err = pubSubHelpers.PublishExecutionStatusOnPubSub(
				topicId, processTestInstructionExecutionRequestAsJsonString)

			// Something went wrong
			if returnMessageAckNack == false || err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"ID":      "2c62eb43-12ff-4cdd-aebf-b312b8201a12",
					"topicId": topicId,
					"processTestInstructionExecutionRequestAsJsonString": processTestInstructionExecutionRequestAsJsonString,
					"err":                 err,
					"returnMessageString": returnMessageString,
				}).Error("Something went wrong when sending pubsub-ExecutionStatus-message to TesterGui")

				continue
			}
		}

		/*
			// The 'subscribeToMessagesStreamResponse' that will be added into Channel message
			var subscribeToMessagesStreamResponse *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
			subscribeToMessagesStreamResponse = &fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse{
				ProtoFileVersionUsedByClient:     fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
				IsKeepAliveMessage:               false,
				ExecutionsStatus:                 testCaseExecutionsStatusAndTestInstructionExecutionsStatusMessage,
				OriginalMessageCreationTimeStamp: originalMessageCreationTimeStampForGrpc,
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

				common_config.Logger.WithFields(logrus.Fields{
					"Id":                             "b248f8b4-e610-4986-8f7d-2688eaf282cf",
					"messageToTestGuiForwardChannel": messageToTestGuiForwardChannel,
				}).Debug("ExecutionStatusMessage was put on channel")

			}
		*/

	}
}
