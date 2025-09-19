package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListTestSuiteExecutions - *********************************************************************
// Call from TesterGui to GuiExecution to get a list of all TestSuiteExecutions with their current execution status
func (s *fenixGuiExecutionServerGrpcServicesServer) ListTestSuiteExecutions(
	ctx context.Context,
	listTestSuiteExecutionsRequest *fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsRequest) (
	*fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse,
	error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "32daf1d0-bfec-4caf-9f17-9d5650130d7b",
	}).Debug("Incoming 'gRPC - ListTestSuiteExecutions'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "dff05164-73c2-4f54-b045-5c94efedd1c5",
	}).Debug("Outgoing 'gRPC - ListTestSuiteExecutions'")

	// Current user
	userIdOnComputer := listTestSuiteExecutionsRequest.UserAndApplicationRunTimeIdentification.GetUserIdOnComputer()

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdOnComputer,
		listTestSuiteExecutionsRequest.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse{
			AckNackResponse:                             returnMessage,
			TestSuiteExecutionsList:                     nil,
			LatestUniqueTestSuiteExecutionDatabaseRowId: 0,
			MoreRowsExists:                              false,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var listTestSuiteExecutionsResponse *fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse
	var err error

	// Get TestSuitesOnExecutionQueue from CloudDB
	listTestSuiteExecutionsResponse, err = fenixGuiExecutionServerObject.listTestSuiteExecutionsFromCloudDB(
		listTestSuiteExecutionsRequest)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestSuiteExecutionsResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ListTestSuiteExecutions from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestSuiteExecutionsList:                     nil,
			LatestUniqueTestSuiteExecutionDatabaseRowId: 0,
			MoreRowsExists:                              false,
		}

		// Exiting
		return responseMessage, nil
	}

	return listTestSuiteExecutionsResponse, nil
}
