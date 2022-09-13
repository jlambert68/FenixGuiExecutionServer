package main

import (
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) processTestInstructionsBasicTestInstructionInformation(immatureTestInstructionMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage) (err error) {

	var (
	//	basicTestInstructionInformation            fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage
	//basicTestInstructionInformationSQLCount    int64
	//immatureTestInstructionInformation fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage
	//immatureTestInstructionInformationSQLCount int64
	//immatureSubTestCaseModel                   fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	//immatureSubTestCaseModelSQLCount           int64
	)

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	// **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT BTI.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionInformation\" BTI "
	sqlToExecute = sqlToExecute + "ORDER BY BTI.\"DomainUuid\" ASC,  BTI.\"TestInstructionTypeUuid\" ASC, BTI.\"TestInstructionUuid\" ASC; "

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
	//var basicTestInstructionInformation fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var tempTimeStamp time.Time
	//var tempTestInstructionExecutionType string

	// Get number of rows for 'basicTestInstructionInformation'
	//basicTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()
	var (
		nonEditableInformation    fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage
		editableInformation       fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage
		invisibleBasicInformation fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage
		//editableTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableTestInstructionAttributesMessage
		//immatureElementModelMessage                        fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
		//immatureTestInstructionInformationMessage fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		nonEditableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage{}
		editableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage{}
		invisibleBasicInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage{}
		//editableTestInstructionAttribute = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableTestInstructionAttributesMessage{}

		err := rows.Scan(
			// NonEditableInformation
			&nonEditableInformation.DomainUuid,
			&nonEditableInformation.DomainName,
			&nonEditableInformation.TestInstructionOrignalUuid,
			&nonEditableInformation.TestInstructionOriginalName,
			&nonEditableInformation.TestInstructionTypeUuid,
			&nonEditableInformation.TestInstructionTypeName,
			&nonEditableInformation.Deprecated,
			&nonEditableInformation.MajorVersionNumber,
			&nonEditableInformation.MinorVersionNumber,
			&tempTimeStamp,
			&nonEditableInformation.TestInstructionColor,
			&nonEditableInformation.TCRuleDeletion,
			&nonEditableInformation.TCRuleSwap,

			// EditableInformation
			&editableInformation.TestInstructionDescription,
			&editableInformation.TestInstructionMouseOverText,

			// InvisibleBasicInformation
			&invisibleBasicInformation.Enabled,

			// EditableTestInstructionAttribute
			//&tempTestInstructionExecutionType,
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

		// Convert 'tempTestInstructionExecutionType' gRPC-type
		//editableTestInstructionAttribute.TestInstructionExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionExecutionTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionExecutionTypeEnum_value[tempTestInstructionExecutionType])

		// Add 'basicTestInstructionInformation' to map
		testInstructionUuid := nonEditableInformation.TestInstructionOrignalUuid

		_, existsInMap := immatureTestInstructionMessageMap[testInstructionUuid]
		// testInstructionUuid shouldn't exist in map. If so then there is a problem
		if existsInMap == true {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                  "58cd4928-e4b5-4faf-9724-047c1cbc82a1",
				"testInstructionUuid": testInstructionUuid,
				"sqlToExecute":        sqlToExecute,
			}).Fatal("TestInstructionUuid shouldn't exist in map. If so then there is a problem")

		}

		// Create 'basicTestInstructionInformation' of the parts
		basicTestInstructionInformation := fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage{
				DomainUuid:                  nonEditableInformation.DomainUuid,
				DomainName:                  nonEditableInformation.DomainName,
				TestInstructionOrignalUuid:  nonEditableInformation.TestInstructionOrignalUuid,
				TestInstructionOriginalName: nonEditableInformation.TestInstructionOriginalName,
				TestInstructionTypeUuid:     nonEditableInformation.TestInstructionTypeUuid,
				TestInstructionTypeName:     nonEditableInformation.TestInstructionTypeName,
				Deprecated:                  nonEditableInformation.Deprecated,
				MajorVersionNumber:          nonEditableInformation.MajorVersionNumber,
				MinorVersionNumber:          nonEditableInformation.MinorVersionNumber,
				UpdatedTimeStamp:            nonEditableInformation.UpdatedTimeStamp,
				TestInstructionColor:        nonEditableInformation.TestInstructionColor,
				TCRuleDeletion:              nonEditableInformation.TCRuleDeletion,
				TCRuleSwap:                  nonEditableInformation.TCRuleSwap,
			},
			EditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage{
				TestInstructionDescription:   editableInformation.TestInstructionDescription,
				TestInstructionMouseOverText: editableInformation.TestInstructionMouseOverText,
			},
			InvisibleBasicInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage{
				Enabled: invisibleBasicInformation.Enabled},
			//EditableTestInstructionAttributes: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableTestInstructionAttributesMessage{
			//	TestInstructionExecutionType: editableTestInstructionAttribute.TestInstructionExecutionType},
		}

		immatureTestInstructionInformationMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage{}
		immatureElementModelMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{}

		// Create 'immatureTestInstructionMessage' and add 'BasicTestInstructionInformation' and a small part of 'ImmatureSubTestCaseModel'
		newImmatureTestInstructionMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage{
			BasicTestInstructionInformation:    &basicTestInstructionInformation,
			ImmatureTestInstructionInformation: &immatureTestInstructionInformationMessage,
			ImmatureSubTestCaseModel:           &immatureElementModelMessage}

		// Save immatureTestInstructionMessage in map
		immatureTestInstructionMessageMap[testInstructionUuid] = &newImmatureTestInstructionMessage

	}
	return nil
}

// **** immatureTestInstructionInformation **** **** immatureTestInstructionInformation **** **** immatureTestInstructionInformation ****
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) processTestInstructionsImmatureTestInstructionInformation(immatureTestInstructionMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage) (err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT ITII.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"ImmatureTestInstructionInformation\" ITII "
	sqlToExecute = sqlToExecute + "ORDER BY ITII.\"DomainUuid\" ASC, ITII.\"TestInstructionUuid\" ASC,  ITII.\"DropZoneUuid\" ASC, ITII.\"TestInstructionAttributeUuid\" ASC; "

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

	// Get number of rows for 'immatureTestInstructionInformation'
	//immatureTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionInformationMessages
	//immatureTestInstructionInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage)

	// Temp variables used when extracting data
	var domainUuid, previousDomainUuid string
	var domainName string
	var testInstructionUuid, previousTestInstructionUuid string
	var testInstructionName string
	var tempTestInstructionAttributeType string

	// First Row in TestData
	var firstRowInSQLRespons bool
	firstRowInSQLRespons = true

	var (
		availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
		availableDropZones                           []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	)

	var (
		dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
		dropZonePreSetTestInstructionAttributes                                                []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	var firstImmatureElementUuid string

	var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionUuid = ""

	// Initiate a new variable to store the data
	newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
	availableDropZone = newAvailableDropZone

	newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
	dropZonePreSetTestInstructionAttribute = newDropZonePreSetTestInstructionAttribute

	// Extract data from DB result set
	for rows.Next() {

		err := rows.Scan(

			// temp-data which is not stored in object
			&domainUuid,
			&domainName,
			&testInstructionUuid,
			&testInstructionName,

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

			// Attribute Action Command controls have to use the attribute
			&dropZonePreSetTestInstructionAttribute.AttributeActionCommand,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "e514dbca-530d-490e-9fb7-58eaa114a721",
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
				//domainUuid != previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				//domainUuid != previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				//domainUuid == previousDomainUuid &&
				testInstructionUuid == previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid == previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				//domainUuid == previousDomainUuid &&
				testInstructionUuid == previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 4
		}

		// Only TestInstructionUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				//domainUuid == previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 5
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		case 1:
			newDropZonePreSetTestInstructionAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
			dropZonePreSetTestInstructionAttributes = newDropZonePreSetTestInstructionAttributes

			newAvailableDropZones := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
			availableDropZones = newAvailableDropZones

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		// Only TestInstructionUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		case 2, 5:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, previousDropZonePreSetTestInstructionAttribute)

			// Convert to pointer object instead before storing in map
			var dropZonePreSetTestInstructionAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
			for _, tempDropZonePreSetTestInstructionAttributeToStore := range dropZonePreSetTestInstructionAttributes {
				newAdropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
				newAdropZonePreSetTestInstructionAttribute = tempDropZonePreSetTestInstructionAttributeToStore
				dropZonePreSetTestInstructionAttributesToStore = append(dropZonePreSetTestInstructionAttributesToStore, &newAdropZonePreSetTestInstructionAttribute)
			}

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributesToStore

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, previousAvailableDropZone)

			// Add the availableDropZones to the ImmatureTestInstructionInformationMessage-map
			immatureTestInstructionMessage, existsInMap := immatureTestInstructionMessageMap[previousTestInstructionUuid]
			// testInstructionUuid shouldn't exist in map. If so then there is a problem
			if existsInMap == false {
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"Id":                  "9fd1b07e-c87a-4583-869b-b3ed28b44616",
					"testInstructionUuid": testInstructionUuid,
					"sqlToExecute":        sqlToExecute,
				}).Fatal("TestInstructionUuid should exist in map. If not so then there is a problem")
			}

			// Convert to pointer object instead before storing in map
			var availableDropZoneMessageToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
			for _, tempAvailableDropZones := range availableDropZones {
				newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
				newAvailableDropZone = tempAvailableDropZones
				availableDropZoneMessageToStore = append(availableDropZoneMessageToStore, &newAvailableDropZone)
			}

			immatureTestInstructionMessage.ImmatureTestInstructionInformation.AvailableDropZones = availableDropZoneMessageToStore
			immatureTestInstructionMessageMap[previousTestInstructionUuid] = immatureTestInstructionMessage

			// Create fresh versions of variables
			//newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
			//availableDropZone = newAvailableDropZone

			newAailableDropZones := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
			availableDropZones = newAailableDropZones

			//newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
			//dropZonePreSetTestInstructionAttribute = newDropZonePreSetTestInstructionAttribute

			newDropZonePreSetTestInstructionAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
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
			var dropZonePreSetTestInstructionAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
			for _, tempDropZonePreSetTestInstructionAttributeToStore := range dropZonePreSetTestInstructionAttributes {
				newAdropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
				newAdropZonePreSetTestInstructionAttribute = tempDropZonePreSetTestInstructionAttributeToStore
				dropZonePreSetTestInstructionAttributesToStore = append(dropZonePreSetTestInstructionAttributesToStore, &newAdropZonePreSetTestInstructionAttribute)
			}

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributesToStore

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, previousAvailableDropZone)

			// Clear DropZone-attributes-array
			newDropZonePreSetTestInstructionAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
			dropZonePreSetTestInstructionAttributes = newDropZonePreSetTestInstructionAttributes

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                                     "352075ba-32f8-4374-86d9-a936ec91b179",
				"domainUuid":                             domainUuid,
				"previousDomainUuid":                     previousDomainUuid,
				"testInstructionUuid":                    testInstructionUuid,
				"previousTestInstructionUuid":            previousTestInstructionUuid,
				"availableDropZone.DropZoneUuid":         availableDropZone.DropZoneUuid,
				"previousAvailableDropZone.DropZoneUuid": previousAvailableDropZone.DropZoneUuid,
				"dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid":         dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
				"previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid": previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

		}

		// Move actual values into previous-variables
		previousDomainUuid = domainUuid
		previousTestInstructionUuid = testInstructionUuid
		previousAvailableDropZone = availableDropZone
		previousDropZonePreSetTestInstructionAttribute = dropZonePreSetTestInstructionAttribute

		// Create fresh versions of variables
		newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
		availableDropZone = newAvailableDropZone

		newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
		dropZonePreSetTestInstructionAttribute = newDropZonePreSetTestInstructionAttribute

		//newDropZonePreSetTestInstructionAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
		//dropZonePreSetTestInstructionAttributes = newDropZonePreSetTestInstructionAttributes

		// Set to not be the first row
		firstRowInSQLRespons = false

	}

	// Handle last row from database
	// Add the previous DropZone-attributes to the DropZone-array
	dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, previousDropZonePreSetTestInstructionAttribute)

	// Convert to pointer object instead before storing in map
	var dropZonePreSetTestInstructionAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	for _, tempDropZonePreSetTestInstructionAttributeToStore := range dropZonePreSetTestInstructionAttributes {
		newAdropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}
		newAdropZonePreSetTestInstructionAttribute = tempDropZonePreSetTestInstructionAttributeToStore
		dropZonePreSetTestInstructionAttributesToStore = append(dropZonePreSetTestInstructionAttributesToStore, &newAdropZonePreSetTestInstructionAttribute)
	}

	// Add attributes to previousDropZone
	previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributesToStore

	// Add previousAvailableDropZone to array of DropZone
	availableDropZones = append(availableDropZones, previousAvailableDropZone)

	// Add 'basicTestInstructionInformation' to map
	immatureTestInstructionMessage, existsInMap := immatureTestInstructionMessageMap[testInstructionUuid]
	if existsInMap == false {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":                  "0f59327f-84a9-47bd-bfe2-337c3402ab0c",
			"testInstructionUuid": testInstructionUuid,
		}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
	}

	// Convert to pointer object instead before storing in map
	var availableDropZoneMessageToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	for _, tempAvailableDropZones := range availableDropZones {
		newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
		newAvailableDropZone = tempAvailableDropZones
		availableDropZoneMessageToStore = append(availableDropZoneMessageToStore, &newAvailableDropZone)
	}

	// Store the result back in the map
	immatureTestInstructionMessage.ImmatureTestInstructionInformation.AvailableDropZones = availableDropZoneMessageToStore
	immatureTestInstructionMessageMap[testInstructionUuid] = immatureTestInstructionMessage

	return err
}

// **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage ****
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) processTestInstructionsImmatureElementModel(immatureTestInstructionMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage) (err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT IEM.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionInformation\" BTII, "
	sqlToExecute = sqlToExecute + "\"" + usedDBSchema + "\".\"ImmatureElementModelMessage\" IEM "
	sqlToExecute = sqlToExecute + "WHERE BTII.\"TestInstructionUuid\" = IEM.\"TopImmatureElementUuid\" "
	sqlToExecute = sqlToExecute + "ORDER BY IEM.\"DomainUuid\" ASC, IEM.\"TopImmatureElementUuid\" ASC, IEM.\"IsTopElement\" DESC; " //, IEM.\"CurrentElementModelElement\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "c98209fd-150c-4e4c-bcce-303d66523213",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Get number of rows for 'immatureTestInstructionInformation'
	//immatureTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionInformationMessages
	//immatureTestInstructionInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage)

	// Temp variables used when extracting data
	var tempImmatureElementModelDomainUuid string
	var tempImmatureElementModelDomainName string
	var tempTestCaseModelElementTypeAsString string
	var tempIsTopElement bool
	var tempTopElementUuid string
	var previousTempTopElementUuid string

	//var previousOriginalElementUuid string
	//var testInstructionUuid, previousTestInstructionUuid string
	//var testInstructionName string
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
	//previousTestInstructionUuid := ""

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
				"Id":           "7a937579-bb0a-44d4-850f-4cbdd5fff3a5",
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

		// All UUIDs are changed and this is the first row [dataStateChange=1]

		case 1:

			newImmatureElementModelElements := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
			immatureElementModelElements = newImmatureElementModelElements

			// All UUIDs are changed and this is not the first row [dataStateChange=2]
			// A new Element model Element and this is not the first row [dataStateChange=4]

		case 2, 4:

			immatureElementModelElements = append(immatureElementModelElements, previousImmatureElementModelElement)

			// Add immatureElementModelElements to 'immatureTestInstructionMessage' which can be found in map
			var immatureTestInstructionMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage
			var existsInMap bool
			immatureTestInstructionMessage, existsInMap = immatureTestInstructionMessageMap[previousTempTopElementUuid]
			if existsInMap == false {
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"Id": "ef98b5ca-17d5-4bf8-8af4-a1a954736a47",
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

			immatureTestInstructionMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElementsToStore
			immatureTestInstructionMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = previousTempTopElementUuid
			immatureTestInstructionMessageMap[previousTempTopElementUuid] = immatureTestInstructionMessage

			// Create fresh versions of variables
			newIimmatureElementModelElements := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
			immatureElementModelElements = newIimmatureElementModelElements

			// A new Element model Element , but it belongs to same 'ImmatureElementUuid' as previous Element, and this is not the first row [dataStateChange=3]
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

	// Add immatureElementModelElements to 'immatureTestInstructionMessage' which can be found in map
	var immatureTestInstructionMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage
	var existsInMap bool
	immatureTestInstructionMessage, existsInMap = immatureTestInstructionMessageMap[tempTopElementUuid]
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

	immatureTestInstructionMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElementsToStore
	immatureTestInstructionMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = tempTopElementUuid
	immatureTestInstructionMessageMap[tempTopElementUuid] = immatureTestInstructionMessage

	return nil

}
