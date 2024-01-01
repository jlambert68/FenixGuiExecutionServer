package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// ListAllImmatureTestInstructionAttributes - *********************************************************************
// The TestCase Builder asks for all TestInstructions Attributes that the user must add values to in TestCase
func (s *fenixGuiExecutionServerGrpcServicesServer) ListTestCasesUnderExecution(ctx context.Context, listTestCasesUnderExecutionRequest *fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionRequest) (*fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse, error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "30e50aed-e860-467f-abdf-673d37788616",
	}).Debug("Incoming 'gRPC - ListTestCasesUnderExecution'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "e316b0ac-edd9-48f7-a550-dcf1a30aeab8",
	}).Debug("Outgoing 'gRPC - ListTestCasesUnderExecution'")

	// Current user
	userIdOnComputer := listTestCasesUnderExecutionRequest.UserAndApplicationRunTimeIdentification.GetUserIdOnComputer()

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdOnComputer,
		fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
			listTestCasesUnderExecutionRequest.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse{
			AckNackResponse: returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"id":    "6f4ebe57-9606-4977-918d-f9f278127049",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin' in 'ListTestCasesUnderExecution'")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem to do 'DbPool.Begin' in 'ListTestCasesUnderExecution'",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCasesUnderExecution: nil,
		}

		return responseMessage, nil
	}

	// Close db-transaction when leaving this function
	defer txn.Commit(context.Background())

	// Define variables to store data from DB in
	var testCaseUnderExecutionMessage []*fenixExecutionServerGuiGrpcApi.TestCaseUnderExecutionMessage

	// Get users ImmatureTestInstruction-data from CloudDB
	testCaseUnderExecutionMessage, err = fenixGuiExecutionServerObject.listTestCasesUnderExecutionLoadFromCloudDB(txn, userIdOnComputer, listTestCasesUnderExecutionRequest.DomainUuids)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving TestCaseUnderExecutionMessage from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCasesUnderExecution: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesUnderExecutionResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
		TestCasesUnderExecution: testCaseUnderExecutionMessage,
	}

	return responseMessage, nil
}
