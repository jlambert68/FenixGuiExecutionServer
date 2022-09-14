package main

import (
	"context"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadSingleTestCaseExecutionSummaryFromCloudDB(userID string) (singleTestCaseExecutionSummaryMessages []*fenixExecutionServerGuiGrpcApi.SingleTestCaseExecutionSummaryMessage, err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	// **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TIATTR.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructionAttributes\" TIATTR "
	sqlToExecute = sqlToExecute + "ORDER BY TIATTR.\"DomainUuid\" ASC, TIATTR.\"TestInstructionUuid\" ASC, TIATTR.\"TestInstructionAttributeUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "5f769af2-f75a-4ea6-8c3d-2108c9dfb9b7",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Variables to used when extract data from result set
	var tempTestInstructionAttributeInputMask string

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		singleTestCaseExecutionSummaryMessage := fenixExecutionServerGuiGrpcApi.SingleTestCaseExecutionSummaryMessage{}

		err := rows.Scan(
			&singleTestCaseExecutionSummaryMessage.DomainUuid,
			&singleTestCaseExecutionSummaryMessage.DomainName,
			&singleTestCaseExecutionSummaryMessage.TestInstructionUuid,
			&singleTestCaseExecutionSummaryMessage.TestInstructionName,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeUuid,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeName,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeDescription,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeMouseOver,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeTypeUuid,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeTypeName,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeValueAsString,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeValueUuid,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeVisible,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeEnable,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeMandatory,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeVisibleInTestCaseArea,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeIsDeprecated,
			&tempTestInstructionAttributeInputMask,
			&singleTestCaseExecutionSummaryMessage.TestInstructionAttributeUIType,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "7cd322cb-2219-4c4d-a8c8-2770a42b0c23",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Add BondAttribute to BondsAttributes
		singleTestCaseExecutionSummaryMessages = append(singleTestCaseExecutionSummaryMessages, &singleTestCaseExecutionSummaryMessage)

	}

	return singleTestCaseExecutionSummaryMessages, err
}
