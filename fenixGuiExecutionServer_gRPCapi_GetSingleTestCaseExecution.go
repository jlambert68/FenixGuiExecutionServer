package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// GetSingleTestCaseExecution - *********************************************************************
// Get all information for a single TestCaseExecution
func (s *fenixGuiExecutionServerGrpcServicesServer) GetSingleTestCaseExecution(
	ctx context.Context,
	getSingleTestCaseExecutionRequest *fenixExecutionServerGuiGrpcApi.GetSingleTestCaseExecutionRequest) (
	getSingleTestCaseExecutionResponse *fenixExecutionServerGuiGrpcApi.GetSingleTestCaseExecutionResponse,
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "90e7e541-b2dc-49e7-85bd-17226b0eecad",
	}).Debug("Incoming 'gRPC - GetSingleTestCaseExecution'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "6b97b5ec-fafb-4da1-9323-88154ade17d7",
	}).Debug("Outgoing 'gRPC - GetSingleTestCaseExecution'")

	// Current user
	userIdOnComputer := getSingleTestCaseExecutionRequest.UserAndApplicationRunTimeIdentification.UserIdOnComputer

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdOnComputer,
		fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
			getSingleTestCaseExecutionRequest.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		getSingleTestCaseExecutionResponse = &fenixExecutionServerGuiGrpcApi.GetSingleTestCaseExecutionResponse{
			AckNackResponse:           returnMessage,
			TestCaseExecutionResponse: nil,
		}

		// Exiting
		return getSingleTestCaseExecutionResponse, nil
	}

	// Define variables to store data from DB in
	var testCaseExecutionResponseMessages []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionResponseMessage

	// Extract TestCaseExecution from Database
	testCaseExecutionResponseMessages, err = fenixGuiExecutionServerObject.loadFullTestCasesExecutionInformation(
		[]*fenixExecutionServerGuiGrpcApi.TestCaseExecutionKeyMessage{getSingleTestCaseExecutionRequest.TestCaseExecutionKey})

	if err != nil {
		// Something went wrong so return an error to caller
		getSingleTestCaseExecutionResponse = &fenixExecutionServerGuiGrpcApi.GetSingleTestCaseExecutionResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving TestCaseExecutionInformation from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCaseExecutionResponse: nil,
		}

		// Exiting
		return getSingleTestCaseExecutionResponse, nil
	}

	// Exact one TestCaseExecution must be found
	if len(testCaseExecutionResponseMessages) != 1 {
		getSingleTestCaseExecutionResponse = &fenixExecutionServerGuiGrpcApi.GetSingleTestCaseExecutionResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Did not find exact one TestCaseExecution",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCaseExecutionResponse: nil,
		}

		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":                                "7ccf8b52-55de-4c39-8bd7-25615045b0cf",
			"getSingleTestCaseExecutionRequest": getSingleTestCaseExecutionRequest,
		}).Error("Did not find exact one TestCaseExecution")

		// Exiting
		return getSingleTestCaseExecutionResponse, nil
	}

	// Create the response to caller
	getSingleTestCaseExecutionResponse = &fenixExecutionServerGuiGrpcApi.GetSingleTestCaseExecutionResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
		TestCaseExecutionResponse: testCaseExecutionResponseMessages[0],
	}

	return getSingleTestCaseExecutionResponse, nil
}
