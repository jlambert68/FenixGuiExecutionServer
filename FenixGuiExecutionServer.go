package main

import (
	"FenixGuiExecutionServer/broadcastEngine"
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/messagesToExecutionServer"
	"fmt"
	uuidGenerator "github.com/google/uuid"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
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

	// Start Backend gRPC-server
	fenixGuiExecutionServerObject.InitGrpcServer()

}
