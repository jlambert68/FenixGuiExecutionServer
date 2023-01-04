package messagesToExecutionServer

import (
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/gcp"
	"context"
	"fmt"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServer - Fenix Gui Execution Server inform ExecutionServer that there is/are new TestCase(s) on TestCaseExecutionQueue
func (messagesToExecutionServerObject *MessagesToExecutionServerObjectStruct) SendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServer(
	testCaseExecutionsToProcessMessage *fenixExecutionServerGrpcApi.TestCaseExecutionsToProcessMessage) (
	ackNackResponse *fenixExecutionServerGrpcApi.AckNackResponse) {

	messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
		"id":                                 "3d3de917-77fe-4768-a5a5-7e107173d74f",
		"testCaseExecutionsToProcessMessage": testCaseExecutionsToProcessMessage,
	}).Debug("Incoming 'SendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServer'")

	messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
		"id": "787a2437-7a81-4629-a8ef-ca676a9e18d3",
	}).Debug("Outgoing 'SendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServer'")

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string
	var err error

	ctx = context.Background()

	if FenixExecutionServerGrpcClient == nil {

		// Set up connection to Server
		ctx, err = messagesToExecutionServerObject.SetConnectionToExecutionServer(ctx)
		if err != nil {

			// Set Error codes to return message
			var errorCodes []fenixExecutionServerGrpcApi.ErrorCodesEnum
			var errorCode fenixExecutionServerGrpcApi.ErrorCodesEnum

			errorCode = fenixExecutionServerGrpcApi.ErrorCodesEnum_ERROR_UNSPECIFIED
			errorCodes = append(errorCodes, errorCode)

			// Create Return message
			ackNackResponse = &fenixExecutionServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     fmt.Sprintf("Couldn't set up connection to ExecutionServer"),
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(messagesToExecutionServerObject.GetHighestFenixExecutionServerProtoFileVersion()),
			}

			return ackNackResponse
		}
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
			"ID": "e4992093-6d22-40d6-a30c-f1e14e05253d",
		}).Debug("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when ExecutionServer is run on GCP
	if common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(ctx) //messagesToExecutionServerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {

			// Set Error codes to return message
			var errorCodes []fenixExecutionServerGrpcApi.ErrorCodesEnum
			var errorCode fenixExecutionServerGrpcApi.ErrorCodesEnum

			errorCode = fenixExecutionServerGrpcApi.ErrorCodesEnum_ERROR_UNSPECIFIED
			errorCodes = append(errorCodes, errorCode)

			// Create Return message
			ackNackResponse = &fenixExecutionServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("Couldn't generate GCPAccessToken for ExecutioNServer'. Return message: '%s'", returnMessageString),
				ErrorCodes: errorCodes,
			}

			return ackNackResponse

		}

	}

	// Finish the preparation of the message to ExecutionServer
	testCaseExecutionsToProcessMessage.ProtoFileVersionUsedByClient = fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(messagesToExecutionServerObject.GetHighestFenixExecutionServerProtoFileVersion())

	// slice with sleep time, in milliseconds, between each attempt to do gRPC-call to ExecutionServer
	var sleepTimeBetweenGrpcCallAttempts []int
	sleepTimeBetweenGrpcCallAttempts = []int{100, 200, 300, 300, 500, 500, 1000, 1000, 1000, 1000} // Total: 5.9 seconds

	// Do multiple attempts to do gRPC-call to ExecutionServer, when it fails
	var numberOfgRPCCallAttempts int
	var gRPCCallAttemptCounter int
	numberOfgRPCCallAttempts = len(sleepTimeBetweenGrpcCallAttempts)
	gRPCCallAttemptCounter = 0

	var informThatThereAreNewTestCasesOnExecutionQueueResponse *fenixExecutionServerGrpcApi.AckNackResponse

	for {

		// Do gRPC-call to ExecutionServer
		informThatThereAreNewTestCasesOnExecutionQueueResponse, err = FenixExecutionServerGrpcClient.InformThatThereAreNewTestCasesOnExecutionQueue(ctx, testCaseExecutionsToProcessMessage)

		// Exit when there was a success call
		if err == nil && informThatThereAreNewTestCasesOnExecutionQueueResponse.AckNack == true {
			return informThatThereAreNewTestCasesOnExecutionQueueResponse
		}

		// Add to counter for how many gRPC-call-attempts to Worker that have been done
		gRPCCallAttemptCounter = gRPCCallAttemptCounter + 1

		// Shouldn't happen
		if err != nil {
			messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
				"ID":    "e0e2175f-6ea0-4437-92dd-5f83359c8ea5",
				"error": err,
			}).Error("Problem to do gRPC-call to FenixExecutionServer for 'InformThatThereAreNewTestCasesOnExecutionQueue'")

			// Only return the error after last attempt
			if gRPCCallAttemptCounter >= numberOfgRPCCallAttempts {

				// Set Error codes to return message
				var errorCodes []fenixExecutionServerGrpcApi.ErrorCodesEnum
				var errorCode fenixExecutionServerGrpcApi.ErrorCodesEnum

				errorCode = fenixExecutionServerGrpcApi.ErrorCodesEnum_ERROR_UNSPECIFIED
				errorCodes = append(errorCodes, errorCode)

				// Create Return message
				ackNackResponse = &fenixExecutionServerGrpcApi.AckNackResponse{
					AckNack:                      false,
					Comments:                     fmt.Sprintf("Error when doing gRPC-call for ExecutionServe. Error message is: '%s'", err.Error()),
					ErrorCodes:                   errorCodes,
					ProtoFileVersionUsedByClient: fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(messagesToExecutionServerObject.GetHighestFenixExecutionServerProtoFileVersion()),
				}

				return ackNackResponse

			}

			// Sleep for some time before retrying to connect
			time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenGrpcCallAttempts[gRPCCallAttemptCounter-1]))

		} else if informThatThereAreNewTestCasesOnExecutionQueueResponse.AckNack == false {

			// ExecutionServer couldn't handle gPRC call
			messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
				"ID":                           "c104fc85-c6ca-4084-a756-409e53491bfe",
				"Message from ExecutionServer": informThatThereAreNewTestCasesOnExecutionQueueResponse.Comments,
			}).Error("Problem to do gRPC-call to FenixExecutionServer for 'InformThatThereAreNewTestCasesOnExecutionQueue'")

			// Set Error codes to return message
			var errorCodes []fenixExecutionServerGrpcApi.ErrorCodesEnum
			var errorCode fenixExecutionServerGrpcApi.ErrorCodesEnum

			errorCode = fenixExecutionServerGrpcApi.ErrorCodesEnum_ERROR_UNSPECIFIED
			errorCodes = append(errorCodes, errorCode)

			// Create Return message
			ackNackResponse = &fenixExecutionServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     fmt.Sprintf("AckNack=false when doing gRPC-call for ExecutionServer. Message is: '%s'", informThatThereAreNewTestCasesOnExecutionQueueResponse.Comments),
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: informThatThereAreNewTestCasesOnExecutionQueueResponse.ProtoFileVersionUsedByClient,
			}

			return ackNackResponse

		}

	}

	return informThatThereAreNewTestCasesOnExecutionQueueResponse

}
