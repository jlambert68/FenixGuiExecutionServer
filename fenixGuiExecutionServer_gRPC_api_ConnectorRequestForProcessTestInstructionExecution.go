package main

import (
	"errors"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SubscribeToMessageStream
// Used to send Messages from Fenix backend to TesterGui. TesterGui connects to GuiExecutionServer and then Responses are streamed back to TesterGui
func (s *fenixGuiExecutionServerGrpcServicesServer) SubscribeToMessageStream(emptyParameter *fenixExecutionServerGuiGrpcApi.EmptyParameter, streamServer fenixExecutionServerGuiGrpcApi.FenixExecutionServerGuiGrpcServicesForGuiClient_SubscribeToMessageStreamServer) (err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "d986194e-ec8c-4198-8160-bd7eb9838aca",
	}).Debug("Incoming 'gRPCServer - SubscribeToMessageStream'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "1b9fb882-f3aa-4ffa-b575-910569aec6c4",
	}).Debug("Outgoing 'gRPCServer - SubscribeToMessageStream'")

	// Calling system
	userId := "TesterGui"

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiExecutionServerObject.isClientUsingCorrectTestDataProtoFileVersion(userId, emptyParameter.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		return errors.New(returnMessage.Comments)
	}

	// Recreate channel for incoming TestInstructionExecution from Execution Server
	messageToTesterGuiForwardChannel = make(chan messageToTestGuiForwardChannelStruct)

	// Local channel to decide when Server stopped sending
	done := make(chan bool)

	go func() {

		// We have an active connection to TesterGui
		TesterGuiHasConnected = true

		for {
			// Wait for incoming TestInstructionExecution from Execution Server
			executionForwardChannelMessage := <-messageToTesterGuiForwardChannel

			testInstructionExecution := executionForwardChannelMessage.subscribeToMessagesStreamResponse

			// If TesterGui stops responding then exit
			if TesterGuiHasConnected == false {
				done <- true //close(done)

				return
			}

			err = streamServer.Send(testInstructionExecution)
			if err != nil {

				// We don't have an active connection to TesterGui
				TesterGuiHasConnected = false

				fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
					"id":                       "70ab1dcb-0be3-49b6-b49a-694bab529ed4",
					"err":                      err,
					"testInstructionExecution": testInstructionExecution,
				}).Error("Got some problem when doing reversed streaming of Messages to TesterGui. Stopping Reversed Streaming")

				// Have the gRPC-call be continued, end stream server
				done <- true //close(done)

				return

			}

			// Check if message only was a keep alive message to TesterGui
			if executionForwardChannelMessage.isKeepAliveMessage == false {

				// Is a standard TestInstructionExecution that was sent to TesterGui
				fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
					"id":                       "6f5e6dc7-cef5-4008-a4ea-406be80ded4c",
					"testInstructionExecution": testInstructionExecution,
				}).Debug("Success in reversed streaming TestInstructionExecution to TesterGui")

			} else {

				// Is a keep alive message that was sent to TesterGui
				fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
					"id":                       "c1d5a756-b7fa-48ae-953e-59dedd0671f4",
					"testInstructionExecution": testInstructionExecution,
				}).Debug("Success in reversed streaming TestInstructionExecution to TesterGui")
			}
		}

	}()

	// Feed 'messageToTesterGuiForwardChannel' with messages every 15 seconds to check if TesterGui is alive
	go func() {

		// Create keep alive message
		var subscribeToMessagesStreamResponse *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
		subscribeToMessagesStreamResponse = &fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse{
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(fenixGuiExecutionServerObject.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			IsKeepAliveMessage:           true,
			ExecutionsStatus:             nil,
		}
		var keepAliveMessageToTesterGui messageToTestGuiForwardChannelStruct
		keepAliveMessageToTesterGui = messageToTestGuiForwardChannelStruct{
			subscribeToMessagesStreamResponse: subscribeToMessagesStreamResponse,
			isKeepAliveMessage:                true,
		}

		var messageWasPickedFromExecutionForwardChannel bool

		for {

			// Sleep for 15 seconds before continue
			time.Sleep(time.Second * 15)

			// If we haven't got an answer from TesterGui in 30 seconds then it must be down.
			// We can get in this state if 'messageToTesterGuiForwardChannel' is full and nobody picks the message from queue
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
			messageToTesterGuiForwardChannel <- keepAliveMessageToTesterGui
			messageWasPickedFromExecutionForwardChannel = true

		}
	}()

	// Server stopped so wait for new connection
	<-done

	TesterGuiHasConnected = false

	return err

}
