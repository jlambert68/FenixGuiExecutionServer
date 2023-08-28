package main

import (
	"FenixGuiExecutionServer/broadcastEngine_ExecutionStatusUpdate"
	"FenixGuiExecutionServer/common_config"
	"errors"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SubscribeToMessageStream
// Used to send Messages from Fenix backend to TesterGui. TesterGui connects to GuiExecutionServer and then Responses are streamed back to TesterGui
func (s *fenixGuiExecutionServerGrpcServicesServer) SubscribeToMessageStream(userAndApplicationRunTimeIdentificationMessage *fenixExecutionServerGuiGrpcApi.UserAndApplicationRunTimeIdentificationMessage, streamServer fenixExecutionServerGuiGrpcApi.FenixExecutionServerGuiGrpcServicesForGuiClient_SubscribeToMessageStreamServer) (err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "d986194e-ec8c-4198-8160-bd7eb9838aca",
		"userAndApplicationRunTimeIdentificationMessage": userAndApplicationRunTimeIdentificationMessage,
	}).Debug("Incoming 'gRPCServer - SubscribeToMessageStream'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "1b9fb882-f3aa-4ffa-b575-910569aec6c4",
	}).Debug("Outgoing 'gRPCServer - SubscribeToMessageStream'")

	// Calling system
	userId := "TesterGui"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userId,
		userAndApplicationRunTimeIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		return errors.New(returnMessage.Comments)
	}

	// Check if TesterGui:s 'ApplicationRunTimeUuid' already exits
	var testCaseExecutionsSubscriptionChannelInformation *broadcastEngine_ExecutionStatusUpdate.TestCaseExecutionsSubscriptionChannelInformationStruct
	var existInMap bool

	testCaseExecutionsSubscriptionChannelInformation, existInMap =
		broadcastEngine_ExecutionStatusUpdate.TestCaseExecutionsSubscriptionChannelInformationMap[broadcastEngine_ExecutionStatusUpdate.ApplicationRunTimeUuidType(
			userAndApplicationRunTimeIdentificationMessage.ApplicationRunTimeUuid)]

	if existInMap == true {
		// Just recreate channel for incoming TestInstructionExecution from Execution Server for this TesterGui
		var tempMessageToTesterGuiForwardChannel broadcastEngine_ExecutionStatusUpdate.MessageToTesterGuiForwardChannelType
		tempMessageToTesterGuiForwardChannel = make(
			chan broadcastEngine_ExecutionStatusUpdate.MessageToTestGuiForwardChannelStruct,
			broadcastEngine_ExecutionStatusUpdate.MessageToTesterGuiForwardChannelMaxSize)
		testCaseExecutionsSubscriptionChannelInformation.MessageToTesterGuiForwardChannel = &tempMessageToTesterGuiForwardChannel

	} else {
		// Create the full ChannelObject from scratch
		var tempMessageToTesterGuiForwardChannel broadcastEngine_ExecutionStatusUpdate.MessageToTesterGuiForwardChannelType
		tempMessageToTesterGuiForwardChannel = make(
			chan broadcastEngine_ExecutionStatusUpdate.MessageToTestGuiForwardChannelStruct,
			broadcastEngine_ExecutionStatusUpdate.MessageToTesterGuiForwardChannelMaxSize)

		testCaseExecutionsSubscriptionChannelInformation = &broadcastEngine_ExecutionStatusUpdate.TestCaseExecutionsSubscriptionChannelInformationStruct{
			ApplicationRunTimeUuid: broadcastEngine_ExecutionStatusUpdate.ApplicationRunTimeUuidType(
				userAndApplicationRunTimeIdentificationMessage.ApplicationRunTimeUuid),
			LastConnectionFromTesterGui:      time.Now().UTC(),
			MessageToTesterGuiForwardChannel: &tempMessageToTesterGuiForwardChannel,
		}

		// Save ChannelObject in 'TestCaseExecutionsSubscriptionChannelInformationMap'
		broadcastEngine_ExecutionStatusUpdate.TestCaseExecutionsSubscriptionChannelInformationMap[broadcastEngine_ExecutionStatusUpdate.ApplicationRunTimeUuidType(
			userAndApplicationRunTimeIdentificationMessage.ApplicationRunTimeUuid)] = testCaseExecutionsSubscriptionChannelInformation

	}

	// Local channel to decide when Server stopped sending
	done := make(chan bool)

	go func() {

		// We have an active connection to TesterGui
		TesterGuiHasConnected = true

		for {
			// Wait for ExecutionStatus-message
			executionForwardChannelMessage := <-*testCaseExecutionsSubscriptionChannelInformation.MessageToTesterGuiForwardChannel

			common_config.Logger.WithFields(logrus.Fields{
				"Id":                             "1c876526-9b96-41a4-b694-7e176ef46656",
				"executionForwardChannelMessage": executionForwardChannelMessage,
			}).Debug("Receive ExecutionStatusMessage from channel")

			var executionsStatusMessage *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
			executionsStatusMessage = executionForwardChannelMessage.SubscribeToMessagesStreamResponse

			// If TesterGui stops responding then exit
			if TesterGuiHasConnected == false {
				done <- true //close(done)

				return
			}

			err = streamServer.Send(executionsStatusMessage)
			if err != nil {

				// We don't have an active connection to TesterGui
				TesterGuiHasConnected = false

				fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
					"id":                      "70ab1dcb-0be3-49b6-b49a-694bab529ed4",
					"err":                     err,
					"executionsStatusMessage": executionsStatusMessage,
				}).Error("Got some problem when doing reversed streaming of Messages to TesterGui. Stopping Reversed Streaming")

				// If message is not a keep alive message, then Put message back on Subscription channel
				if executionForwardChannelMessage.IsKeepAliveMessage == false {
					fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
						"id":                      "5210db0a-c91d-47da-954e-cb0a78667c76",
						"executionsStatusMessage": executionsStatusMessage,
					}).Debug("Put message back on channel")

					*testCaseExecutionsSubscriptionChannelInformation.MessageToTesterGuiForwardChannel <- executionForwardChannelMessage
				}
				// Have the gRPC-call be continued, end stream server
				done <- true //close(done)

				return

			}

			// Check if message only was a keep alive message to TesterGui
			if executionForwardChannelMessage.IsKeepAliveMessage == false {

				// Is a standard TestInstructionExecution that was sent to TesterGui
				fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
					"id":                      "6f5e6dc7-cef5-4008-a4ea-406be80ded4c",
					"executionsStatusMessage": executionsStatusMessage,
				}).Debug("Success in reversed streaming TestInstructionExecution to TesterGui")

			} else {

				// Is a keep alive message that was sent to TesterGui
				fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
					"id":                      "c1d5a756-b7fa-48ae-953e-59dedd0671f4",
					"executionsStatusMessage": executionsStatusMessage,
				}).Debug("Success in reversed streaming Keep-alive-message to TesterGui")
			}
		}

	}()

	// Feed 'MessageToTesterGuiForwardChannel' with messages every 15 seconds to check if TesterGui is alive
	go func() {

		// Create keep alive message
		var subscribeToMessagesStreamResponse *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
		subscribeToMessagesStreamResponse = &fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse{
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			IsKeepAliveMessage:           true,
			ExecutionsStatus:             nil,
		}
		var keepAliveMessageToTesterGui broadcastEngine_ExecutionStatusUpdate.MessageToTestGuiForwardChannelStruct
		keepAliveMessageToTesterGui = broadcastEngine_ExecutionStatusUpdate.MessageToTestGuiForwardChannelStruct{
			SubscribeToMessagesStreamResponse: subscribeToMessagesStreamResponse,
			IsKeepAliveMessage:                true,
		}

		var messageWasPickedFromExecutionForwardChannel bool

		for {

			// Sleep for 15 seconds before continue
			time.Sleep(time.Second * 15)

			// If we haven't got an answer from TesterGui in 30 seconds then it must be down.
			// We can get in this state if 'MessageToTesterGuiForwardChannel' is full and nobody picks the message from queue
			messageWasPickedFromExecutionForwardChannel = false

			go func() {
				time.Sleep(time.Second * 30)
				if messageWasPickedFromExecutionForwardChannel == false {
					// Stop in channel
					TesterGuiHasConnected = false
					fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
						"id": "ad24ded4-4218-4ddd-93bb-2b8ec1a1a046",
					}).Debug("No answer regarding Keep Alive-message, TesterGui is not responding")

					done <- true //close(done)

				}
			}()

			// Send Keep Alive message on channel to be sent to TesterGui
			*testCaseExecutionsSubscriptionChannelInformation.MessageToTesterGuiForwardChannel <- keepAliveMessageToTesterGui
			messageWasPickedFromExecutionForwardChannel = true

		}
	}()

	// Server stopped so wait for new connection
	<-done

	TesterGuiHasConnected = false

	return err

}
