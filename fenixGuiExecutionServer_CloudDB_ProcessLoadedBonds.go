package main

import (
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) processVisibleBondAttributesInformation() (bondsAttributes []fenixTestCaseBuilderServerGrpcApi.BasicBondInformationMessage_VisibleBondAttributesMessage, err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	// **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT IB.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"ImmatureBonds\" IB "
	sqlToExecute = sqlToExecute + "ORDER BY IB.\"TCRuleDeletion\"; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "7161a37c-96ff-4df7-acfe-da41b80d2c29",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Variables to used when extract data from result set
	var tempTimeStamp time.Time
	var tempTestCaseModelElementTypeGrpcMappingId int
	var tempTestCaseModelElementTypeAsString string

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		bondAttributes := fenixTestCaseBuilderServerGrpcApi.BasicBondInformationMessage_VisibleBondAttributesMessage{}

		err := rows.Scan(
			&bondAttributes.BondUuid,
			&bondAttributes.BondName,
			&bondAttributes.BondDescription,
			&bondAttributes.BondMouseOverText,
			&bondAttributes.Deprecated,
			&bondAttributes.Enabled,
			&bondAttributes.Visible,
			&bondAttributes.BondColor,
			&bondAttributes.CanBeDeleted,
			&bondAttributes.CanBeSwappedOut,
			&bondAttributes.ShowBondAttributes,
			&bondAttributes.TCRuleDeletion,
			&bondAttributes.TCRuleSwap,
			&tempTimeStamp,
			&tempTestCaseModelElementTypeAsString, //&bondAttributes.TestCaseModelElementType,
			&tempTestCaseModelElementTypeGrpcMappingId,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "5dbe5e92-eefb-44bb-825f-3f6013902cf0",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert TimeStamp into proto-format for TimeStamp
		bondAttributes.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Convert 'tempTestCaseModelElementTypeAsString' into gRPC-type
		bondAttributes.TestCaseModelElementType = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[tempTestCaseModelElementTypeAsString])

		// Add BondAttribute to BondsAttributes
		bondsAttributes = append(bondsAttributes, bondAttributes)

	}

	return bondsAttributes, err
}
