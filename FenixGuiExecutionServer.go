package main

import (
	"FenixGuiExecutionServer/broadcastingEngine"
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/messagesToExecutionServer"
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

	// Start listen for Broadcasts regarding change in status TestCaseExecutions and TestInstructionExecutions
	broadcastingEngine.InitiateAndStartBroadcastNotifyEngine()

	// Start Backend gRPC-server
	fenixGuiExecutionServerObject.InitGrpcServer()

}
