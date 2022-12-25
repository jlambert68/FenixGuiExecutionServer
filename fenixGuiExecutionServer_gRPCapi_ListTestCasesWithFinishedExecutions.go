package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// ListTestCasesWithFinishedExecutions - *********************************************************************
// List all TestCaseExecutions that has finished their execution
func (s *fenixGuiExecutionServerGrpcServicesServer) ListTestCasesWithFinishedExecutions(ctx context.Context, listTestCasesWithFinishedExecutionsRequest *fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsRequest) (*fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse, error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "33451c5f-1230-4ea1-817e-08afa1c1192b",
	}).Debug("Incoming 'gRPC - ListTestCasesWithFinishedExecutions'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "1cabcbc9-86ab-4ffb-b41f-8562c3bc1d75",
	}).Debug("Outgoing 'gRPC - ListTestCasesWithFinishedExecutions'")

	// Current user
	userID := listTestCasesWithFinishedExecutionsRequest.UserAndApplicationRunTimeIdentification.UserId

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userID,
		fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
			listTestCasesWithFinishedExecutionsRequest.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
			AckNackResponse:               returnMessage,
			TestCaseWithFinishedExecution: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "9f4451d2-6d81-45fe-bd2a-5ac3a69f882e",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin' in 'ListTestCasesWithFinishedExecutions'")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem to do 'DbPool.Begin' in 'ListTestCasesWithFinishedExecutions'",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCaseWithFinishedExecution: nil,
		}

		return responseMessage, nil
	}

	// Close db-transaction when leaving this function
	defer txn.Commit(context.Background())

	// Define variables to store data from DB in
	var listTestCasesWithFinishedExecutionsResponse []*fenixExecutionServerGuiGrpcApi.TestCaseWithFinishedExecutionMessage

	// Get users ImmatureTestInstruction-data from CloudDB
	listTestCasesWithFinishedExecutionsResponse, err = fenixGuiExecutionServerObject.listTestCasesWithFinishedExecutionsLoadFromCloudDB(txn, userID, listTestCasesWithFinishedExecutionsRequest.DomainUuids)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ListTestCasesWithFinishedExecutions from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCaseWithFinishedExecution: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesWithFinishedExecutionsResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
		TestCaseWithFinishedExecution: listTestCasesWithFinishedExecutionsResponse,
	}

	return responseMessage, nil
}
