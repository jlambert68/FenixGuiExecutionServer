package main

import (
	"FenixGuiExecutionServer/broadcastEngine"
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/gcp"
	"FenixGuiExecutionServer/messagesToExecutionServer"
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

		//log.Println("Close DB_session: %v", DB_session)
		//DB_session.Close()
	}
}

func fenixGuiExecutionServerMain() {

	// Create Unique Uuid for run time instance used as identification when communication with GuiExecutionServer
	common_config.ApplicationRunTimeUuid = uuidGenerator.New().String()
	fmt.Println("common_config.ApplicationRunTimeUuid: " + common_config.ApplicationRunTimeUuid)

	// Connect to CloudDB
	fenixSyncShared.ConnectToDB()

	// Set up BackendObject
	fenixGuiExecutionServerObject = &fenixGuiExecutionServerObjectStruct{}

	// Init logger
	fenixGuiExecutionServerObject.InitLogger("")
	common_config.Logger = fenixGuiExecutionServerObject.logger

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// Initiate 'MessagesToExecutionServerObject' tfor messages to be sent to ExecutionServer
	messagesToExecutionServer.MessagesToExecutionServerObject = messagesToExecutionServer.MessagesToExecutionServerObjectStruct{
		Logger: fenixGuiExecutionServerObject.logger,
	}

	// Initiate the handler that 'knows' who is subscribing to which TestCaseExecutions, regarding status updates
	broadcastEngine.InitiateSubscriptionHandler()

	// Start listen for Broadcasts regarding change in status TestCaseExecutions and TestInstructionExecutions
	broadcastEngine.InitiateAndStartBroadcastNotifyEngine()

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
