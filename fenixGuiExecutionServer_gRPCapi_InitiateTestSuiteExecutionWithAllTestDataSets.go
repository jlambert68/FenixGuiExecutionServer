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

// InitiateTestSuiteExecutionWithAllTestDataSets - *********************************************************************
// Initiate a TestExecution from a TestSuite with all its TestDataSets
func (s *fenixGuiExecutionServerGrpcServicesServer) InitiateTestSuiteExecutionWithAllTestDataSets(ctx context.Context,
	initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage *fenixExecutionServerGuiGrpcApi.InitiateTestSuiteExecutionWithAllTestDataSetsRequestMessage) (
	*fenixExecutionServerGuiGrpcApi.InitiateSingleTestSuiteExecutionResponseMessage, error) {

	var err error

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "fb95f94a-40bb-4a59-8cbc-351a1647e096",
		"initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage": initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage,
	}).Debug("Incoming 'gRPC - InitiateTestSuiteExecutionWithAllTestDataSets'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "18b9bb93-ae89-48c4-9917-1e6a037f7cdd",
	}).Debug("Outgoing 'gRPC - InitiateTestSuiteExecutionWithAllTestDataSets'")

	// Check if Client is using correct proto files version
	ackNackRespons := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage.UserAndApplicationRunTimeIdentification.GetUserIdOnComputer(),
		initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient)
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

	// Initiate response variable
	var initiateSingleSuiteCaseExecutionResponseMessage *fenixExecutionServerGuiGrpcApi.
		InitiateSingleTestSuiteExecutionResponseMessage

	// Load all TestDataSets to be used
	//var testDataFromSimpleTestDataAreaFileMessages []*fenixGuiExecutionServerObject.TestDataFromOneSimpleTestDataAreaFileMessage
	var testDataForTestCaseExecutionMessages []*fenixExecutionServerGuiGrpcApi.TestDataForTestCaseExecutionMessage
	testDataForTestCaseExecutionMessages, err = fenixGuiExecutionServerObject.
		initiateLoadTestSuitesAllTestDataSetsFromCloudDB(
			initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage.GetTestSuiteUuid())

	if err != nil {
		return initiateSingleSuiteCaseExecutionResponseMessage, nil
	}

	// Save TestCaseExecutions in Cloud DB
	initiateSingleSuiteCaseExecutionResponseMessage = fenixGuiExecutionServerObject.
		prepareInitiateTestSuiteExecutionSaveToCloudDB(
			initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage.UserAndApplicationRunTimeIdentification,
			initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage.GetTestSuiteUuid(),
			fenixExecutionServerGuiGrpcApi.ExecutionPriorityEnum_HIGH_SINGLE_TESTSUITE,
			initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage.GetExecutionStatusReportLevel(),
			testDataForTestCaseExecutionMessages)

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
					initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage.ExecutionStatusReportLevel),
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
				"initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage": initiateTestSuiteExecutionWithAllTestDataSetsRequestMessage,
			}).Error("Problem when doing gRPC-call to FenixExecutionServer")

			return //initiateSingleSuiteCaseExecutionResponseMessage, nil
		}

	}()

	return initiateSingleSuiteCaseExecutionResponseMessage, nil

}
