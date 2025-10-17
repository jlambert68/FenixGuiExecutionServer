package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// GetSingleTestSuiteExecution - *********************************************************************
// Get all information for a single TestSuiteExecution
func (s *fenixGuiExecutionServerGrpcServicesServer) GetSingleTestSuiteExecution(
	ctx context.Context,
	getSingleTestSuiteExecutionRequest *fenixExecutionServerGuiGrpcApi.GetSingleTestSuiteExecutionRequest) (
	getSingleTestSuiteExecutionResponse *fenixExecutionServerGuiGrpcApi.GetSingleTestSuiteExecutionResponse,
	err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "641e2fd2-c796-4bc7-a39a-68fd92869730",
	}).Debug("Incoming 'gRPC - GetSingleTestSuiteExecution'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "c59bf65a-6603-4433-87f2-a1cc075186c3",
	}).Debug("Outgoing 'gRPC - GetSingleTestSuiteExecution'")

	// Current user
	userIdOnComputer := getSingleTestSuiteExecutionRequest.UserAndApplicationRunTimeIdentification.UserIdOnComputer

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdOnComputer,
		fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
			getSingleTestSuiteExecutionRequest.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		getSingleTestSuiteExecutionResponse = &fenixExecutionServerGuiGrpcApi.GetSingleTestSuiteExecutionResponse{
			AckNackResponse:            returnMessage,
			TestSuiteExecutionResponse: nil,
		}

		// Exiting
		return getSingleTestSuiteExecutionResponse, nil
	}

	// Define variables to store data from DB in
	var testSuiteExecutionResponseMessages []*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionResponseMessage

	// Extract TestSuiteExecution from Database
	testSuiteExecutionResponseMessages, err = fenixGuiExecutionServerObject.loadFullTestSuitesExecutionInformation(
		[]*fenixExecutionServerGuiGrpcApi.TestSuiteExecutionKeyMessage{getSingleTestSuiteExecutionRequest.TestSuiteExecutionKey},
		getSingleTestSuiteExecutionRequest.UserAndApplicationRunTimeIdentification.GetGCPAuthenticatedUser())

	if err != nil {
		// Something went wrong so return an error to caller
		getSingleTestSuiteExecutionResponse = &fenixExecutionServerGuiGrpcApi.GetSingleTestSuiteExecutionResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving TestSuiteExecutionInformation from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestSuiteExecutionResponse: nil,
		}

		// Exiting
		return getSingleTestSuiteExecutionResponse, nil
	}

	// Exact one TestSuiteExecution must be found
	if len(testSuiteExecutionResponseMessages) != 1 {
		getSingleTestSuiteExecutionResponse = &fenixExecutionServerGuiGrpcApi.GetSingleTestSuiteExecutionResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Did not find exact one TestSuiteExecution",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestSuiteExecutionResponse: nil,
		}

		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":                                 "bb4f9fba-8b69-49e7-968f-15d8fcfc483e",
			"getSingleTestSuiteExecutionRequest": getSingleTestSuiteExecutionRequest,
		}).Error("Did not find exact one TestSuiteExecution")

		// Exiting
		return getSingleTestSuiteExecutionResponse, nil
	}

	// Create the response to caller
	getSingleTestSuiteExecutionResponse = &fenixExecutionServerGuiGrpcApi.GetSingleTestSuiteExecutionResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
		TestSuiteExecutionResponse: testSuiteExecutionResponseMessages[0],
	}

	return getSingleTestSuiteExecutionResponse, nil
}
