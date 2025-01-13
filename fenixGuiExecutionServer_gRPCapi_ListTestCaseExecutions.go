package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListTestCaseExecutions - *********************************************************************
// Call from TesterGui to GuiExecution to get a list of all TestCaseExecutions with their current execution status
func (s *fenixGuiExecutionServerGrpcServicesServer) ListTestCaseExecutions(
	ctx context.Context,
	listTestCaseExecutionsRequest *fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsRequest) (
	*fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse,
	error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "22e9da2f-5e5b-4221-b1b3-2e50fb360687",
	}).Debug("Incoming 'gRPC - ListTestCaseExecutions'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "37cf1243-c248-4b6b-afea-5b40507504a0",
	}).Debug("Outgoing 'gRPC - ListTestCaseExecutions'")

	// Current user
	userIdOnComputer := listTestCaseExecutionsRequest.UserAndApplicationRunTimeIdentification.GetUserIdOnComputer()

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdOnComputer,
		listTestCaseExecutionsRequest.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse{
			AckNackResponse:                            returnMessage,
			TestCaseExecutionsList:                     nil,
			LatestUniqueTestCaseExecutionDatabaseRowId: 0,
			MoreRowsExists:                             false,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var listTestCaseExecutionsResponse *fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse
	var err error

	// Get TestCasesOnExecutionQueue from CloudDB
	listTestCaseExecutionsResponse, err = fenixGuiExecutionServerObject.listTestCaseExecutionsFromCloudDB(
		listTestCaseExecutionsRequest)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCaseExecutionsResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ListTestCaseExecutions from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCaseExecutionsList:                     nil,
			LatestUniqueTestCaseExecutionDatabaseRowId: 0,
			MoreRowsExists:                             false,
		}

		// Exiting
		return responseMessage, nil
	}

	return listTestCaseExecutionsResponse, nil
}
