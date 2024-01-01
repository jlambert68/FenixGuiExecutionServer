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
func (s *fenixGuiExecutionServerGrpcServicesServer) ListTestCasesOnExecutionQueue(ctx context.Context, listTestCasesInExecutionQueueRequest *fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueRequest) (*fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse, error) {

	// Define the response message
	var responseMessage *fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "a88c93c3-cc86-4b4b-86ca-a11dd606b242",
	}).Debug("Incoming 'gRPC - ListTestCasesOnExecutionQueue'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "1ce24c71-11ae-4f76-a473-ce794e4610e6",
	}).Debug("Outgoing 'gRPC - ListTestCasesOnExecutionQueue'")

	// Current user
	userIdOnComputer := listTestCasesInExecutionQueueRequest.UserAndApplicationRunTimeIdentification.GetUserIdOnComputer()

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdOnComputer,
		fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(
			listTestCasesInExecutionQueueRequest.UserAndApplicationRunTimeIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse{
			AckNackResponse:           returnMessage,
			TestCasesInExecutionQueue: nil,
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
		}).Error("Problem to do 'DbPool.Begin' in 'ListTestCasesOnExecutionQueue'")

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		initiateSingleTestCaseExecutionResponseMessage := &fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem to do 'DbPool.Begin' in 'ListTestCasesOnExecutionQueue'",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCasesInExecutionQueue: nil,
		}

		return initiateSingleTestCaseExecutionResponseMessage, nil
	}

	// Close db-transaction when leaving this function
	defer txn.Commit(context.Background())

	// Define variables to store data from DB in
	var testCaseExecutionBasicInformation []*fenixExecutionServerGuiGrpcApi.TestCaseExecutionBasicInformationMessage

	// Get TestCasesOnExecutionQueue from CloudDB
	testCaseExecutionBasicInformation, err = fenixGuiExecutionServerObject.listTestCasesOnExecutionQueueLoadFromCloudDB(txn, userIdOnComputer, listTestCasesInExecutionQueueRequest.DomainUuids)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse{
			AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving TestCaseExecutionBasicInformationMessage from database",
				ErrorCodes:                   []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum{fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
			},
			TestCasesInExecutionQueue: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixExecutionServerGuiGrpcApi.ListTestCasesInExecutionQueueResponse{
		AckNackResponse: &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(common_config.GetHighestFenixGuiExecutionServerProtoFileVersion()),
		},
		TestCasesInExecutionQueue: testCaseExecutionBasicInformation,
	}

	return responseMessage, nil
}
