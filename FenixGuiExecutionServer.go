package main

import (
	"FenixGuiExecutionServer/broadcastEngine_ExecutionStatusUpdate"
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/gcp"
	"FenixGuiExecutionServer/messagesToExecutionServer"
	pubsub "FenixGuiExecutionServer/outgoingPubSubMessages"
	"FenixGuiExecutionServer/testerGuiOwnerEngine"
	"context"
	"fmt"
	uuidGenerator "github.com/google/uuid"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// Used for only process cleanup once
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// Cleanup before close down application
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{}).Info("Clean up and shut down servers")

		// Stop Backend gRPC Server
		fenixGuiExecutionServerObject.StopGrpcServer()

		// Close Database Connection
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id": "0bc657a0-db05-4781-9083-02ebd44567ff",
		}).Info("Closing Database connection")

		fenixSyncShared.DbPool.Close()
	}
}

func fenixGuiExecutionServerMain() {

	// Create Unique Uuid for run time instance used as identification when communication with GuiExecutionServer
	common_config.ApplicationRunTimeUuid = uuidGenerator.New().String()
	fmt.Println("common_config.ApplicationRunTimeUuid: " + common_config.ApplicationRunTimeUuid)

	// Set start up time for this instance
	common_config.ApplicationRunTimeStartUpTime = time.Now()

	// Connect to CloudDB
	fenixSyncShared.ConnectToDB()

	// Set up BackendObject
	fenixGuiExecutionServerObject = &fenixGuiExecutionServerObjectStruct{}

	// Init logger
	fenixGuiExecutionServerObject.InitLogger("")
	common_config.Logger = fenixGuiExecutionServerObject.logger

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	pubsub.MyTestPubSubFunctions()

	// Start TesterGuiOwnerEngine
	testerGuiOwnerEngine.InitiateTesterGuiOwnerEngine()

	// Send that this 'GuiExecutionServer' is closing down, over Broadcast system
	defer func() {

		// Create response channel
		var responseChannel chan bool
		responseChannel = make(chan bool)

		// Put message on 'testGuiExecutionEngineChannel' to be processed
		var tempGuiExecutionServerIsClosingDown common_config.GuiExecutionServerIsClosingDownStruct
		tempGuiExecutionServerIsClosingDown = common_config.GuiExecutionServerIsClosingDownStruct{
			GuiExecutionServerApplicationId:                     common_config.ApplicationRunTimeUuid,
			MessageTimeStamp:                                    time.Now(),
			CurrentGuiExecutionServerIsClosingDownReturnChannel: &responseChannel,
			GuiExecutionServerResponsibilities:                  nil,
		}

		var testerGuiOwnerEngineChannelCommand common_config.TesterGuiOwnerEngineChannelCommandStruct
		testerGuiOwnerEngineChannelCommand = common_config.TesterGuiOwnerEngineChannelCommandStruct{
			TesterGuiOwnerEngineChannelCommand:                                 common_config.ChannelCommand_ThisGuiExecutionServerIsClosingDown,
			TesterGuiIsClosingDown:                                             nil,
			GuiExecutionServerIsClosingDown:                                    &tempGuiExecutionServerIsClosingDown,
			ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination: nil,
			UserUnsubscribesToUserAndTestCaseExecutionCombination:              nil,
		}

		// Put on GuiOwnerEngineChannel
		common_config.TesterGuiOwnerEngineChannelEngineCommandChannel <- &testerGuiOwnerEngineChannelCommand

		// Wait until command has been processed by 'GuiOwnerEngine'
		<-responseChannel

	}()

	// Initiate 'MessagesToExecutionServerObject' for messages to be sent to ExecutionServer
	messagesToExecutionServer.MessagesToExecutionServerObject = messagesToExecutionServer.MessagesToExecutionServerObjectStruct{
		Logger: fenixGuiExecutionServerObject.logger,
	}

	// Initiate the handler that 'knows' who is subscribing to which TestCaseExecutions, regarding status updates
	broadcastEngine_ExecutionStatusUpdate.InitiateSubscriptionHandler()

	// Start listen for Broadcasts regarding change in status TestCaseExecutions and TestInstructionExecutions
	broadcastEngine_ExecutionStatusUpdate.InitiateAndStartBroadcastNotifyEngine()

	// When ExecutionServer runs on GCP, then set up access
	if common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP { //&&
		//common_config. GCPAuthentication == true &&
		//common_config.TurnOffCallToWorker == false {

		gcp.Gcp = gcp.GcpObjectStruct{}

		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

		// Generate first time Access token
		_, returnMessageAckNack, returnMessageString := gcp.Gcp.GenerateGCPAccessToken(ctx)
		if returnMessageAckNack == false {

			// If there was any problem then exit program
			common_config.Logger.WithFields(logrus.Fields{
				"id": "20c90d94-eef7-4819-ba8c-b7a56a39f995",
			}).Fatalf("Couldn't generate access token for GCP, return message: '%s'", returnMessageString)

		}
	}

	// Start Backend gRPC-server
	fenixGuiExecutionServerObject.InitGrpcServer()

}
