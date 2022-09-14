package main

import (
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) processTestInstructionContainersBasicTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage) (err error) {

	var (
	//	basicTestInstructionContainerInformation            fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage
	//basicTestInstructionContainerInformationSQLCount    int64
	//immatureTestInstructionContainerInformation fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage
	//immatureTestInstructionContainerInformationSQLCount int64
	//immatureSubTestCaseModel                   fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	//immatureSubTestCaseModelSQLCount           int64
	)

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	// **** BasicTestInstructionContainerInformation **** **** BasicTestInstructionContainerInformation **** **** BasicTestInstructionContainerInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT BTIC.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionContainerInformation\" BTIC "
	sqlToExecute = sqlToExecute + "ORDER BY BTIC.\"DomainUuid\" ASC,  BTIC.\"TestInstructionContainerTypeUuid\" ASC, BTIC.\"TestInstructionContainerUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b944c506-4ded-4f5e-98c4-06f272d16e1a",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	//var basicTestInstructionContainerInformation fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var tempTimeStamp time.Time
	var tempTestInstructionContainerExecutionType string

	// Get number of rows for 'basicTestInstructionContainerInformation'
	//basicTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()
	var (
		nonEditableInformation                    fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage
		editableInformation                       fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage
		invisibleBasicInformation                 fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage
		editableTestInstructionContainerAttribute fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage
		//immatureElementModelMessage                        fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
		//immatureTestInstructionContainerInformationMessage fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		nonEditableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage{}
		editableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage{}
		invisibleBasicInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage{}
		editableTestInstructionContainerAttribute = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage{}

		err := rows.Scan(
			// NonEditableInformation
			&nonEditableInformation.DomainUuid,
			&nonEditableInformation.DomainName,
			&nonEditableInformation.TestInstructionContainerUuid,
			&nonEditableInformation.TestInstructionContainerName,
			&nonEditableInformation.TestInstructionContainerTypeUuid,
			&nonEditableInformation.TestInstructionContainerTypeName,
			&nonEditableInformation.Deprecated,
			&nonEditableInformation.MajorVersionNumber,
			&nonEditableInformation.MinorVersionNumber,
			&tempTimeStamp,
			&nonEditableInformation.TestInstructionContainerColor,
			&nonEditableInformation.TCRuleDeletion,
			&nonEditableInformation.TCRuleSwap,

			// EditableInformation
			&editableInformation.TestInstructionContainerDescription,
			&editableInformation.TestInstructionContainerMouseOverText,

			// InvisibleBasicInformation
			&invisibleBasicInformation.Enabled,

			// EditableTestInstructionContainerAttribute
			&tempTestInstructionContainerExecutionType,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "7d082f7c-f987-44e7-97b7-c3c1652955c3",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert TimeStamp into proto-format for TimeStamp
		nonEditableInformation.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Convert 'tempTestInstructionContainerExecutionType' gRPC-type
		editableTestInstructionContainerAttribute.TestInstructionContainerExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_value[tempTestInstructionContainerExecutionType])

		// Add 'basicTestInstructionContainerInformation' to map
		testInstructionContainerUuid := nonEditableInformation.TestInstructionContainerUuid

		_, existsInMap := immatureTestInstructionContainerMessageMap[testInstructionContainerUuid]
		// testInstructionContainerUuid shouldn't exist in map. If so then there is a problem
		if existsInMap == true {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                           "58cd4928-e4b5-4faf-9724-047c1cbc82a1",
				"testInstructionContainerUuid": testInstructionContainerUuid,
				"sqlToExecute":                 sqlToExecute,
			}).Fatal("TestInstructionContainerUuid shouldn't exist in map. If so then there is a problem")

		}

		// Create 'basicTestInstructionContainerInformation' of the parts
		basicTestInstructionContainerInformation := fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage{
				DomainUuid:                       nonEditableInformation.DomainUuid,
				DomainName:                       nonEditableInformation.DomainName,
				TestInstructionContainerUuid:     nonEditableInformation.TestInstructionContainerUuid,
				TestInstructionContainerName:     nonEditableInformation.TestInstructionContainerName,
				TestInstructionContainerTypeUuid: nonEditableInformation.TestInstructionContainerTypeUuid,
				TestInstructionContainerTypeName: nonEditableInformation.TestInstructionContainerTypeName,
				Deprecated:                       nonEditableInformation.Deprecated,
				MajorVersionNumber:               nonEditableInformation.MajorVersionNumber,
				MinorVersionNumber:               nonEditableInformation.MinorVersionNumber,
				UpdatedTimeStamp:                 nonEditableInformation.UpdatedTimeStamp,
				TestInstructionContainerColor:    nonEditableInformation.TestInstructionContainerColor,
				TCRuleDeletion:                   nonEditableInformation.TCRuleDeletion,
				TCRuleSwap:                       nonEditableInformation.TCRuleSwap,
			},
			EditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage{
				TestInstructionContainerDescription:   editableInformation.TestInstructionContainerDescription,
				TestInstructionContainerMouseOverText: editableInformation.TestInstructionContainerMouseOverText,
			},
			InvisibleBasicInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage{
				Enabled: invisibleBasicInformation.Enabled},
			EditableTestInstructionContainerAttributes: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage{
				TestInstructionContainerExecutionType: editableTestInstructionContainerAttribute.TestInstructionContainerExecutionType},
		}

		immatureTestInstructionContainerInformationMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage{}
		immatureElementModelMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{}

		// Create 'immatureTestInstructionContainerMessage' and add 'BasicTestInstructionInformation' and a small part of 'ImmatureSubTestCaseModel'
		newImmatureTestInstructionContainerMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage{
			BasicTestInstructionContainerInformation:    &basicTestInstructionContainerInformation,
			ImmatureTestInstructionContainerInformation: &immatureTestInstructionContainerInformationMessage,
			ImmatureSubTestCaseModel:                    &immatureElementModelMessage}

		// Save immatureTestInstructionContainerMessage in map
		immatureTestInstructionContainerMessageMap[testInstructionContainerUuid] = &newImmatureTestInstructionContainerMessage

	}
	return nil
}

// **** immatureTestInstructionContainerInformation **** **** immatureTestInstructionContainerInformation **** **** immatureTestInstructionContainerInformation ****
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) processTestInstructionContainersImmatureTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage) (err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT ITICI.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"ImmatureTestInstructionContainerMessage\" ITICI "
	sqlToExecute = sqlToExecute + "ORDER BY ITICI.\"DomainUuid\" ASC, ITICI.\"TestInstructionContainerUuid\" ASC,  ITICI.\"DropZoneUuid\" ASC, ITICI.\"TestInstructionAttributeUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "aa4b0e8e-3644-491d-be99-8c87ea9b9c23",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Get number of rows for 'immatureTestInstructionContainerInformation'
	//immatureTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionContainerInformationMessages
	//immatureTestInstructionContainerInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage)

	// Temp variables used when extracting data
	var domainUuid, previousDomainUuid string
	var domainName string
	var testInstructionContainerUuid, previousTestInstructionContainerUuid string
	var testInstructionContainerName string
	var tempTestInstructionAttributeType string
	// First Row in TestData
	var firstRowInSQLRespons bool
	firstRowInSQLRespons = true

	var (
		availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
		availableDropZones                           []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
	)

	var (
		dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
		dropZonePreSetTestInstructionAttributes                                                []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	var firstImmatureElementUuid string

	var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionContainerUuid = ""

	// Initiate a new variable to store the data
	newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
	availableDropZone = newAvailableDropZone

	newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
	dropZonePreSetTestInstructionAttribute = newDropZonePreSetTestInstructionAttribute

	// Extract data from DB result set
	for rows.Next() {

		err := rows.Scan(

			// temp-data which is not stored in object
			&domainUuid,
			&domainName,
			&testInstructionContainerUuid,
			&testInstructionContainerName,

			// DropZone-data
			&availableDropZone.DropZoneUuid,
			&availableDropZone.DropZoneName,
			&availableDropZone.DropZoneDescription,
			&availableDropZone.DropZoneMouseOver,
			&availableDropZone.DropZoneColor,

			// DropZoneAttributes-data
			&tempTestInstructionAttributeType,
			&dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
			&dropZonePreSetTestInstructionAttribute.TestInstructionAttributeName,
			&dropZonePreSetTestInstructionAttribute.AttributeValueAsString,
			&dropZonePreSetTestInstructionAttribute.AttributeValueUuid,

			// Reference to first element in element-model
			&firstImmatureElementUuid,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "525079b7-8484-4e61-a811-fa863a41ee2f",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert 'tempTestInstructionAttributeType' into gRPC-type
		dropZonePreSetTestInstructionAttribute.TestInstructionAttributeType = fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum_value[tempTestInstructionAttributeType])

		// Handle the correct order of building together the full object
		dataStateChange = 0

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				domainUuid != previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid != previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionContainerUuid == previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid == previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionContainerUuid == previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 4
		}

		// Only TestInstructionContainerUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 5
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		case 1:
			newDropZonePreSetTestInstructionAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
			dropZonePreSetTestInstructionAttributes = newDropZonePreSetTestInstructionAttributes

			newAvailableDropZones := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
			availableDropZones = newAvailableDropZones

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		// Only TestInstructionContainerUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		case 2, 5:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, previousDropZonePreSetTestInstructionAttribute)

			// Convert to pointer object instead before storing in map
			var dropZonePreSetTestInstructionAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
			for _, tempDropZonePreSetTestInstructionAttributeToStore := range dropZonePreSetTestInstructionAttributes {
				newAdropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
				newAdropZonePreSetTestInstructionAttribute = tempDropZonePreSetTestInstructionAttributeToStore
				dropZonePreSetTestInstructionAttributesToStore = append(dropZonePreSetTestInstructionAttributesToStore, &newAdropZonePreSetTestInstructionAttribute)
			}

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributesToStore

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, previousAvailableDropZone)

			// Add the availableDropZones to the ImmatureTestInstructionInformationMessage-map
			immatureTestInstructionContainerMessage, existsInMap := immatureTestInstructionContainerMessageMap[previousTestInstructionContainerUuid]
			// testInstructionContainerUuid shouldn't exist in map. If so then there is a problem
			if existsInMap == false {
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"Id":                           "9fd1b07e-c87a-4583-869b-b3ed28b44616",
					"testInstructionContainerUuid": testInstructionContainerUuid,
					"sqlToExecute":                 sqlToExecute,
				}).Fatal("TestInstructionContainerUuid should exist in map. If not so then there is a problem")
			}

			// Convert to pointer object instead before storing in map
			var availableDropZoneMessageToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
			for _, tempAvailableDropZones := range availableDropZones {
				newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
				newAvailableDropZone = tempAvailableDropZones
				availableDropZoneMessageToStore = append(availableDropZoneMessageToStore, &newAvailableDropZone)
			}

			immatureTestInstructionContainerMessage.ImmatureTestInstructionContainerInformation.AvailableDropZones = availableDropZoneMessageToStore
			immatureTestInstructionContainerMessageMap[previousTestInstructionContainerUuid] = immatureTestInstructionContainerMessage

			// Create fresh versions of variables
			newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
			availableDropZone = newAvailableDropZone

			newAailableDropZones := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
			availableDropZones = newAailableDropZones

			newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
			dropZonePreSetTestInstructionAttribute = newDropZonePreSetTestInstructionAttribute

			newDropZonePreSetTestInstructionAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
			dropZonePreSetTestInstructionAttributes = newDropZonePreSetTestInstructionAttributes

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		case 3:
			// Add the DropZone attribute to the array for attributes
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, previousDropZonePreSetTestInstructionAttribute)

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		case 4:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, previousDropZonePreSetTestInstructionAttribute)

			// Convert to pointer object instead before storing in map
			var dropZonePreSetTestInstructionAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
			for _, tempDropZonePreSetTestInstructionAttributeToStore := range dropZonePreSetTestInstructionAttributes {
				newAdropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
				newAdropZonePreSetTestInstructionAttribute = tempDropZonePreSetTestInstructionAttributeToStore
				dropZonePreSetTestInstructionAttributesToStore = append(dropZonePreSetTestInstructionAttributesToStore, &newAdropZonePreSetTestInstructionAttribute)
			}

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributesToStore

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, previousAvailableDropZone)

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                                     "0779886a-8280-42b6-9434-46ec1afd1d7f",
				"domainUuid":                             domainUuid,
				"previousDomainUuid":                     previousDomainUuid,
				"testInstructionContainerUuid":           testInstructionContainerUuid,
				"previousTestInstructionContainerUuid":   previousTestInstructionContainerUuid,
				"availableDropZone.DropZoneUuid":         availableDropZone.DropZoneUuid,
				"previousAvailableDropZone.DropZoneUuid": previousAvailableDropZone.DropZoneUuid,
				"dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid":         dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
				"previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid": previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

		}

		// Move actual values into previous-variables
		previousDomainUuid = domainUuid
		previousTestInstructionContainerUuid = testInstructionContainerUuid
		previousAvailableDropZone = availableDropZone
		previousDropZonePreSetTestInstructionAttribute = dropZonePreSetTestInstructionAttribute

		// Create fresh versions of variables
		newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
		availableDropZone = newAvailableDropZone

		newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
		dropZonePreSetTestInstructionAttribute = newDropZonePreSetTestInstructionAttribute

		newDropZonePreSetTestInstructionAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
		dropZonePreSetTestInstructionAttributes = newDropZonePreSetTestInstructionAttributes

		// Set to not be the first row
		firstRowInSQLRespons = false

	}

	// Handle last row from database
	// Add the previous DropZone-attributes to the DropZone-array
	dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, previousDropZonePreSetTestInstructionAttribute)

	// Convert to pointer object instead before storing in map
	var dropZonePreSetTestInstructionAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	for _, tempDropZonePreSetTestInstructionAttributeToStore := range dropZonePreSetTestInstructionAttributes {
		newAdropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
		newAdropZonePreSetTestInstructionAttribute = tempDropZonePreSetTestInstructionAttributeToStore
		dropZonePreSetTestInstructionAttributesToStore = append(dropZonePreSetTestInstructionAttributesToStore, &newAdropZonePreSetTestInstructionAttribute)
	}

	// Add attributes to previousDropZone
	previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributesToStore

	// Add previousAvailableDropZone to array of DropZone
	availableDropZones = append(availableDropZones, previousAvailableDropZone)

	// Add 'basicTestInstructionContainerInformation' to map
	immatureTestInstructionContainerMessage, existsInMap := immatureTestInstructionContainerMessageMap[testInstructionContainerUuid]
	if existsInMap == false {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":                           "8630d2e6-261b-4dab-a499-71463346c5a3",
			"testInstructionContainerUuid": testInstructionContainerUuid,
		}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
	}

	// Convert to pointer object instead before storing in map
	var availableDropZoneMessageToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
	for _, tempAvailableDropZones := range availableDropZones {
		newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
		newAvailableDropZone = tempAvailableDropZones
		availableDropZoneMessageToStore = append(availableDropZoneMessageToStore, &newAvailableDropZone)
	}

	// Store the result back in the map
	immatureTestInstructionContainerMessage.ImmatureTestInstructionContainerInformation.AvailableDropZones = availableDropZoneMessageToStore
	immatureTestInstructionContainerMessageMap[testInstructionContainerUuid] = immatureTestInstructionContainerMessage

	return err
}

// **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage ****
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) processTestInstructionContainersImmatureElementModel(immatureTestInstructionContainerMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage) (err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT IEM.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionContainerInformation\" BTICI, "
	sqlToExecute = sqlToExecute + "\"" + usedDBSchema + "\".\"ImmatureElementModelMessage\" IEM "
	sqlToExecute = sqlToExecute + "WHERE BTICI.\"TestInstructionContainerUuid\" = IEM.\"TopImmatureElementUuid\" "
	sqlToExecute = sqlToExecute + "ORDER BY IEM.\"DomainUuid\" ASC, IEM.\"TopImmatureElementUuid\" ASC, IEM.\"IsTopElement\" DESC; " //, IEM.\"CurrentElementModelElement\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "4ef75e5a-8386-4a1d-a04c-4992ee9d7559",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Get number of rows for 'immatureTestInstructionContainerInformation'
	//immatureTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionInformationMessages
	//immatureTestInstructionContainerInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage)

	// Temp variables used when extracting data
	var tempImmatureElementModelDomainUuid string
	var tempImmatureElementModelDomainName string
	var tempTestCaseModelElementTypeAsString string
	var tempIsTopElement bool
	var tempTopElementUuid string
	var previousTempTopElementUuid string

	//var previousOriginalElementUuid string
	//var testInstructionContainerUuid, previousTestInstructionContainerUuid string
	//var testInstructionContainerName string
	//var tempTestInstructionAttributeType string
	// First Row in TestData
	//var firstRowInSQLRespons bool
	firstRowInSQLRespons := true

	var (
	//availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	//availableDropZones                           []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	)

	var (
	//dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	//dropZonePreSetTestInstructionAttributes                                                []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	//var immatureElementModelMessage fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	var immatureElementModelElement fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
	var previousImmatureElementModelElement fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
	var immatureElementModelElements []fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage

	//var firstImmatureElementUuid string

	//var dataStateChange uint8

	// Clear previous variables
	//previousDomainUuid := ""
	//previousTestInstructionContainerUuid := ""

	// Initiate a new variable to store the data
	newImmatureElementModelElement := fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
	immatureElementModelElement = newImmatureElementModelElement

	previousImmatureElementModelElement = fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}

	// Extract data from DB result set
	for rows.Next() {

		// Initiate new fresh variable
		newImmatureElementModelElement := fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
		immatureElementModelElement = newImmatureElementModelElement

		err = rows.Scan(

			// temp-data which is not stored in object
			&tempImmatureElementModelDomainUuid,
			&tempImmatureElementModelDomainName,

			// ImmatureElementModel

			&immatureElementModelElement.ImmatureElementUuid,
			&immatureElementModelElement.OriginalElementName,
			&immatureElementModelElement.PreviousElementUuid,
			&immatureElementModelElement.NextElementUuid,
			&immatureElementModelElement.FirstChildElementUuid,
			&immatureElementModelElement.ParentElementUuid,
			&tempTestCaseModelElementTypeAsString,
			&immatureElementModelElement.OriginalElementUuid,
			&tempTopElementUuid,
			&tempIsTopElement,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "d4dcd3d8-ab65-46d2-b4a5-85d92481718d",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert 'tempTestCaseModelElementTypeAsString' into gRPC-type
		immatureElementModelElement.TestCaseModelElementType = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[tempTestCaseModelElementTypeAsString])

		// Handle the correct order of building together the full object
		dataStateChange := 0

		// This is the first row, and it is flagged as Top-element [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				tempIsTopElement == true

		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// This is not the first row, and it is flagged as Top-element [dataStateChange=2]
		dataStateChangeFound =
			dataStateChange == 0 &&
				firstRowInSQLRespons == false &&
				tempIsTopElement == true

		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		//  This is not the first row, and it is not flagged as Top-element [dataStateChange=3]
		dataStateChangeFound =
			dataStateChange == 0 &&
				firstRowInSQLRespons == false &&
				tempIsTopElement == false

		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// This is the first row, and it is flagged as Top-element [dataStateChange=1]
		case 1:

			newImmatureElementModelElements := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
			immatureElementModelElements = newImmatureElementModelElements

		// This is not the first row, and it is flagged as Top-element [dataStateChange=2]
		case 2:

			immatureElementModelElements = append(immatureElementModelElements, previousImmatureElementModelElement)

			// Add immatureElementModelElements to 'immatureTestInstructionContainerMessage' which can be found in map
			var immatureTestInstructionContainerMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage
			var existsInMap bool
			immatureTestInstructionContainerMessage, existsInMap = immatureTestInstructionContainerMessageMap[previousTempTopElementUuid]
			if existsInMap == false {
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"Id": "c757d974-805d-4d1c-98e9-464868aa273e",
					"previousImmatureElementModelElement.ImmatureElementUuid": previousImmatureElementModelElement.ImmatureElementUuid,
				}).Fatal("ImmatureElementUuid should exist in map. If not then there is a problem")
			}

			// Convert to pointer object instead before storing in map
			var immatureElementModelElementsToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
			for _, tempImmatureElementModelElement := range immatureElementModelElements {
				newImmatureElementModelElement := fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
				newImmatureElementModelElement = tempImmatureElementModelElement
				immatureElementModelElementsToStore = append(immatureElementModelElementsToStore, &newImmatureElementModelElement)
			}

			immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElementsToStore
			immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = previousTempTopElementUuid
			immatureTestInstructionContainerMessageMap[previousTempTopElementUuid] = immatureTestInstructionContainerMessage

			// Create fresh versions of variables
			newIimmatureElementModelElements := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
			immatureElementModelElements = newIimmatureElementModelElements

		//  This is not the first row, and it is not flagged as Top-element [dataStateChange=3]
		case 3:

			immatureElementModelElements = append(immatureElementModelElements, previousImmatureElementModelElement)

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                                  "24be5ad9-09b3-41a2-81e8-b4171dded878",
				"immatureElementModelElement":         immatureElementModelElements,
				"previousImmatureElementModelElement": previousImmatureElementModelElement,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

		}

		// Move previous values to current
		previousImmatureElementModelElement = immatureElementModelElement
		previousTempTopElementUuid = tempTopElementUuid

		// Set to be not the first row
		firstRowInSQLRespons = false

	}
	// Handle last row from database

	// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
	immatureElementModelElements = append(immatureElementModelElements, immatureElementModelElement)

	// Add immatureElementModelElements to 'immatureTestInstructionContainerMessage' which can be found in map
	var immatureTestInstructionContainerMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage
	var existsInMap bool
	immatureTestInstructionContainerMessage, existsInMap = immatureTestInstructionContainerMessageMap[tempTopElementUuid]
	if existsInMap == false {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "a1744497-782f-4e82-bec0-ae0205c6573f",
			"immatureElementModelElement.ImmatureElementUuid": immatureElementModelElement.ImmatureElementUuid,
		}).Fatal("ImmatureElementUuid should exist in map. If not then there is a problem")
	}

	// Convert to pointer object instead before storing in map
	var immatureElementModelElementsToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
	for _, tempImmatureElementModelElement := range immatureElementModelElements {
		newImmatureElementModelElement := fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
		newImmatureElementModelElement = tempImmatureElementModelElement
		immatureElementModelElementsToStore = append(immatureElementModelElementsToStore, &newImmatureElementModelElement)
	}

	immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElementsToStore
	immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = tempTopElementUuid
	immatureTestInstructionContainerMessageMap[tempTopElementUuid] = immatureTestInstructionContainerMessage

	return nil

}
