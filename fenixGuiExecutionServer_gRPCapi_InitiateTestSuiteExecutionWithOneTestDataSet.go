package main

import (
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/messagesToExecutionServer"
	"context"
	"fmt"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// InitiateTestSuiteExecutionWithOneTestDataSet - *********************************************************************
// Initiate a TestExecution from a TestSuite with one TestDataSet
func (s *fenixGuiExecutionServerGrpcServicesServer) InitiateTestSuiteExecutionWithOneTestDataSet(ctx context.Context,
	initiateTestSuiteExecutionWithOneTestDataSetRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateTestSuiteExecutionWithOneTestDataSetRequestMessage) (
	*fenixExecutionServerGuiGrpcApi.InitiateSingleTestSuiteExecutionResponseMessage, error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "be8da457-ea3e-472d-9437-60661edefc96",
		"initiateTestSuiteExecutionWithOneTestDataSetRequestMessage": initiateTestSuiteExecutionWithOneTestDataSetRequestMessage,
	}).Debug("Incoming 'gRPC - InitiateTestSuiteExecutionWithOneTestDataSet'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "87fd70e3-1815-4712-9d91-94eb7d680937",
	}).Debug("Outgoing 'gRPC - InitiateTestSuiteExecutionWithOneTestDataSet'")

	// Check if Client is using correct proto files version
	ackNackRespons := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		initiateTestSuiteExecutionWithOneTestDataSetRequestMessage.UserAndApplicationRunTimeIdentification.GetUserIdOnComputer(),
		initiateTestSuiteExecutionWithOneTestDataSetRequestMessage.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient)
	if ackNackRespons != nil {
		// Not correct proto-file version is used
		// Exiting
		initiateSingleSuiteCaseExecutionResponseMessage := fenixExecutionServerGuiGrpcApi.
			InitiateSingleTestSuiteExecutionResponseMessage{
			TestCasesInExecutionQueue: nil,
			AckNackResponse:           ackNackRespons,
		}

		return &initiateSingleSuiteCaseExecutionResponseMessage, nil
	}

	// Save TestCaseExecution in Cloud DB
	var initiateSingleSuiteCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.
		InitiateSingleTestSuiteExecutionResponseMessage
	initiateSingleSuiteCaseExecutionResponseMessage = fenixGuiExecutionServerObject.
		prepareInitiateTestSuiteExecutionSaveToCloudDB(
			initiateTestSuiteExecutionWithOneTestDataSetRequestMessage.UserAndApplicationRunTimeIdentification,
			initiateTestSuiteExecutionWithOneTestDataSetRequestMessage.GetTestSuiteUuid(),
			fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_HIGH_SINGLE_TESTSUITE,
			initiateTestSuiteExecutionWithOneTestDataSetRequestMessage.GetExecutionStatusReportLevel(),
			[]*fenixExecutionServerGuiGrpcApi.TestDataForTestCaseExecutionMessage{initiateTestSuiteExecutionWithOneTestDataSetRequestMessage.GetTestDataForTestCaseExecution()})

	// Exit due to error in saving TestCaseExecution in database
	if initiateSingleSuiteCaseExecutionResponseMessage.AckNackResponse.AckNack == false {
		return initiateSingleSuiteCaseExecutionResponseMessage, nil
	}

	// Send Execution to ExecutionServer
	go func() {
		// Prepare message to be sent to ExecutionServer
		var testCaseExecutionsToProcessMessage *fenixExecutionServerGrpcApi.TestCaseExecutionsToProcessMessage
		var testCaseExecutionsToProcess []*fenixExecutionServerGrpcApi.TestCaseExecutionToProcess

		// Loop a TestCaseExecutions to inform ExecutionServer to execute all of them
		for _, tempTestCaseInExecutionQueue := range initiateSingleSuiteCaseExecutionResponseMessage.GetTestCasesInExecutionQueue() {

			var testCaseExecutionToProcess *fenixExecutionServerGrpcApi.TestCaseExecutionToProcess
			testCaseExecutionToProcess = &fenixExecutionServerGrpcApi.TestCaseExecutionToProcess{
				TestCaseExecutionsUuid:   tempTestCaseInExecutionQueue.GetTestCaseExecutionUuid(),
				TestCaseExecutionVersion: 1,
				ExecutionStatusReportLevel: fenixExecutionServerGrpcApi.ExecutionStatusReportLevelEnum(
					initiateTestSuiteExecutionWithOneTestDataSetRequestMessage.ExecutionStatusReportLevel),
			}

			testCaseExecutionsToProcess = append(testCaseExecutionsToProcess, testCaseExecutionToProcess)
		}

		testCaseExecutionsToProcessMessage = &fenixExecutionServerGrpcApi.TestCaseExecutionsToProcessMessage{
			TestCaseExecutionsToProcess: testCaseExecutionsToProcess,
		}

		// Trigger ExecutionEngine to start process TestCase from TestCaseExecution-queue
		var sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse *fenixExecutionServerGrpcApi.
			AckNackResponse
		sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse = messagesToExecutionServer.
			MessagesToExecutionServerObject.SendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServer(
			testCaseExecutionsToProcessMessage)

		// If triggering ExecutionServer to read TestCaseExecutionQueue wasn't successful then change 'initiateSingleSuiteCaseExecutionResponseMessage'
		if sendInformThatThereAreNewTestCasesOnExecutionQueueToExecutionServerResponse.AckNack == false {
			var ackNackResponseToRespond *fenixExecutionServerGuiGrpcApi.AckNackResponse
			ackNackResponseToRespond = &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:    initiateSingleSuiteCaseExecutionResponseMessage.AckNackResponse.AckNack,
				Comments:   fmt.Sprintf("Message from ExecutionServer is: '%s'", initiateSingleSuiteCaseExecutionResponseMessage.AckNackResponse.Comments),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.
					CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			}

			initiateSingleSuiteCaseExecutionResponseMessage.AckNackResponse = ackNackResponseToRespond

			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"id": "e5dc557d-d902-421e-8b1f-0c73e71b54c1",
				"initiateTestSuiteExecutionWithOneTestDataSetRequestMessage": initiateTestSuiteExecutionWithOneTestDataSetRequestMessage,
			}).Error("Problem when doing gRPC-call to FenixExecutionServer")

			return //initiateSingleSuiteCaseExecutionResponseMessage, nil
		}

	}()

	return initiateSingleSuiteCaseExecutionResponseMessage, nil

}
