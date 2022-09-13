package main

import (
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// ****************************************************************************************************************
// Load data from CloudDB
//

// Load TestInstructions for Client

/*
func (fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct) loadClientsImmatureTestInstructionsFromCloudDB(userID string) (cloudDBImmatureTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage, err error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"Id": "38fbd4e2-cfe8-405c-84ce-1667c2292c58",
	}).Debug("Entering: loadClientsImmatureTestInstructionsFromCloudDB()")

	defer func() {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id": "6acee974-1280-48f5-9c4f-886aeff58863",
		}).Debug("Exiting: loadClientsImmatureTestInstructionsFromCloudDB()")
	}()

	var (
		basicTestInstructionInformation            fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage
		basicTestInstructionInformationSQLCount    int64
		immatureTestInstructionInformation         fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage
		immatureTestInstructionInformationSQLCount int64
		//immatureSubTestCaseModel                   fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
		//immatureSubTestCaseModelSQLCount           int64
	)

	ImmatureTestInstructionMessageMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage)

	/* Example
	   "DomainUuid"                   uuid      not null,
	   "DomainName"                   varchar   not null,
	   "TestInstructionUuid"          uuid      not null (Key)
	   "TestInstructionName"          varchar   not null,
	   "TestInstructionTypeUuid"      uuid      not null,
	   "TestInstructionTypeName"      varchar   not null,
	   "TestInstructionDescription"   varchar   not null,
	   "TestInstructionMouseOverText" varchar   not null,
	   "Deprecated"                   boolean   not null,
	   "Enabled"                      boolean   not null,
	   "MajorVersionNumber"           integer   not null,
	   "MinorVersionNumber"           integer   not null,
	   "UpdatedTimeStamp"             timestamp not null



	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	// **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT BTI.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionInformation\" BTI "
	sqlToExecute = sqlToExecute + "ORDER BY BTI.\"DomainUuid\" ASC,  BTI.\"TestInstructionTypeUuid\" ASC, BTI.\"TestInstructionUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "2f130d7e-f8aa-466f-b29d-0fb63608c1a6",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Variables to used when extract data from result set
	//var basicTestInstructionInformation fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var tempTimeStamp time.Time

	// Get number of rows for 'basicTestInstructionInformation'
	basicTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()
	var (
		nonEditableInformation      fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage
		editableInformation         fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage
		invisibleBasicInformation   fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage
		immatureElementModelMessage fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		nonEditableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage{}
		editableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage{}
		invisibleBasicInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage{}

		err := rows.Scan(
			// NonEditableInformation
			&nonEditableInformation.DomainUuid,
			&nonEditableInformation.DomainName,
			&nonEditableInformation.TestInstructionUuid,
			&nonEditableInformation.TestInstructionName,
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
		)

		if err != nil {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":           "e7925b78-327c-40ad-9144-ae4a8a6f35f5",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert TimeStamp into proto-format for TimeStamp
		nonEditableInformation.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Add 'basicTestInstructionInformation' to map
		testInstructionUuid := nonEditableInformation.TestInstructionUuid
		immatureTestInstructionMessage, existsInMap := ImmatureTestInstructionMessageMap[testInstructionUuid]
		// testInstructionUuid shouldn't exist in map. If so then there is a problem
		if existsInMap == true {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":                  "4713a8c8-c9d0-4315-9341-27365d64cdc8",
				"testInstructionUuid": testInstructionUuid,
				"sqlToExecute":        sqlToExecute,
			}).Fatal("TestInstructionUuid shouldn't exist in map. If so then there is a problem")

		}

		// Create 'basicTestInstructionInformation' of the parts
		basicTestInstructionInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage{
			NonEditableInformation:    &nonEditableInformation,
			EditableInformation:       &editableInformation,
			InvisibleBasicInformation: &invisibleBasicInformation,
		}

		// Create 'immatureTestInstructionMessage' and add 'BasicTestInstructionInformation' and a small part of 'ImmatureSubTestCaseModel'
		immatureTestInstructionMessage = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage{
			BasicTestInstructionInformation: &basicTestInstructionInformation,
			ImmatureSubTestCaseModel:        &immatureElementModelMessage}

		// Save immatureTestInstructionMessage in map
		ImmatureTestInstructionMessageMap[testInstructionUuid] = immatureTestInstructionMessage

	}

	// **** immatureTestInstructionInformation **** **** immatureTestInstructionInformation **** **** immatureTestInstructionInformation ****
	sqlToExecute = ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"ImmatureTestInstructionInformation\" ITII "
	sqlToExecute = sqlToExecute + "ORDER BY ITII.\"DomainUuid\" ASC, ITII.\"TestInstructionUuid\" ASC,  ITII.\"DropZoneUuid\" ASC, ITII.\"TestInstructionAttributeGuid\" ASC; "

	// Query DB
	rows, err = fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b3ef4fec-9097-46c4-8ff6-85a758967e46",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Get number of rows for 'immatureTestInstructionInformation'
	immatureTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionInformationMessages
	immatureTestInstructionInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage)

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
		availableDropZones                           []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	)

	var (
		dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
		dropZonePreSetTestInstructionAttributes                                                []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	var firstImmatureElementUuid string

	var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionUuid = ""

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		availableDropZone = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
		dropZonePreSetTestInstructionAttribute = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

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
		)

		if err != nil {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":           "9f0618f2-ca04-41e9-aeef-60cd1874f6b7",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert 'tempTestInstructionAttributeType' into gRPC-type
		dropZonePreSetTestInstructionAttribute.TestInstructionAttributeType = fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum_value[tempTestInstructionAttributeType])

		// Handle the correct order of building together the full object
		dataStateChange = 0

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				domainUuid != previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid != previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionUuid == previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid == previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionUuid == previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 4
		}

		// Only TestInstructionUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 5
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		case 1:

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		// Only TestInstructionUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		case 2, 5:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &previousDropZonePreSetTestInstructionAttribute)

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, &previousAvailableDropZone)

			// Add the availableDropZones to the ImmatureTestInstructionInformationMessage-map
			immatureTestInstructionInformation.AvailableDropZones = availableDropZones
			immatureTestInstructionInformationMessagesMap[previousTestInstructionUuid] = immatureTestInstructionInformation

			// Create fresh versions of variables
			immatureTestInstructionInformation = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage{}
			availableDropZones = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
			dropZonePreSetTestInstructionAttributes = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		case 3:
			// Add the DropZone attribute to the array for attributes
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &previousDropZonePreSetTestInstructionAttribute)

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		case 4:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &previousDropZonePreSetTestInstructionAttribute)

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, &previousAvailableDropZone)

			// Create fresh versions of variables
			dropZonePreSetTestInstructionAttributes = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":                                     "ca46deb0-a788-4c68-aefb-27bb7ccaad0d",
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

		// Set to be first row
		firstRowInSQLRespons = false

	}

	// Handle last row from database
	// Add the previous DropZone-attributes to the DropZone-array
	dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &previousDropZonePreSetTestInstructionAttribute)

	// Add attributes to previousDropZone
	previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

	// Add previousAvailableDropZone to array of DropZone
	availableDropZones = append(availableDropZones, &previousAvailableDropZone)

	// Add the availableDropZones to the ImmatureTestInstructionInformationMessage-map
	immatureTestInstructionInformation.AvailableDropZones = availableDropZones
	immatureTestInstructionInformationMessagesMap[previousTestInstructionUuid] = immatureTestInstructionInformation

	// Add 'basicTestInstructionInformation' to map
	immatureTestInstructionMessage, existsInMap := ImmatureTestInstructionMessageMap[testInstructionUuid]
	// testInstructionUuid shouldn't exist in map. If so then there is a problem
	if existsInMap == false {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":                  "4894cf70-08fd-401b-b076-d643ea721abb",
			"testInstructionUuid": testInstructionUuid,
		}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
	}

	// Immature part to 'immatureTestInstructionMessage'
	immatureTestInstructionMessage.ImmatureTestInstructionInformation = &immatureTestInstructionInformation

	// Add 'firstImmatureElementUuid'
	immatureTestInstructionMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = firstImmatureElementUuid

	// Store the information back in the map
	ImmatureTestInstructionMessageMap[testInstructionUuid] = immatureTestInstructionMessage

	// ***************************************************************************************************

	// **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage ****

	sqlToExecute = ""
	sqlToExecute = sqlToExecute + "SELECT IEM.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionInformation\" BTI, "
	sqlToExecute = sqlToExecute + "\"" + usedDBSchema + "\".\"ImmatureElementModelMessage\" IEM "
	sqlToExecute = sqlToExecute + "WHERE BTI.\"TestInstructionUuid\" = IEM.\"ImmatureElementUuid\" "
	sqlToExecute = sqlToExecute + "ORDER BY BTI.\"DomainUuid\" ASC, BTI.\"TestInstructionUuid\" ASC; "

	// Query DB
	rows, err = fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b3ef4fec-9097-46c4-8ff6-85a758967e46",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Get number of rows for 'immatureTestInstructionInformation'
	immatureTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionInformationMessages
	//immatureTestInstructionInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage)

	// Temp variables used when extracting data
	var tempImmatureElementModelDomainUuid, previousTempImmatureDomainUuid string
	var tempImmatureElementModelDomainName string
	var tempTestCaseModelElementTypeAsString string
	//var previousOriginalElementUuid string
	//var testInstructionUuid, previousTestInstructionUuid string
	//var testInstructionName string
	//var tempTestInstructionAttributeType string
	// First Row in TestData
	//var firstRowInSQLRespons bool
	firstRowInSQLRespons = true

	var (
	//availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	//availableDropZones                           []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	)

	var (
	//dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	//dropZonePreSetTestInstructionAttributes                                                []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	//var immatureElementModelMessage fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	var immatureElementModelElement, previousImmatureElementModelElement fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage
	var immatureElementModelElements []*fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage

	//var firstImmatureElementUuid string

	//var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionUuid = ""

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		immatureElementModelElement = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{}

		err = rows.Scan(

			// temp-data which is not stored in object
			&tempImmatureElementModelDomainUuid,
			&tempImmatureElementModelDomainName,

			// ImmatureElementModel

			&immatureElementModelElement.OriginalElementUuid,
			&immatureElementModelElement.OriginalElementName,
			&immatureElementModelElement.PreviousElementUuid,
			&immatureElementModelElement.NextElementUuid,
			&immatureElementModelElement.FirstChildElementUuid,
			&immatureElementModelElement.ParentElementUuid,
			&tempTestCaseModelElementTypeAsString,
			&immatureElementModelElement.CurrentElementModelElement,
		)

		if err != nil {
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":           "808377e2-70ec-4894-bb17-7d92321caaa2",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert 'tempTestCaseModelElementTypeAsString' into gRPC-type
		immatureElementModelElement.TestCaseModelElementType = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[tempTestCaseModelElementTypeAsString])

		// Handle the correct order of building together the full object
		dataStateChange = 0

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				tempImmatureElementModelDomainUuid != previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				tempImmatureElementModelDomainUuid != previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only immatureElementModelElement.OriginalElementUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				tempImmatureElementModelDomainUuid == previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		case 1:

			// All UUIDs are changed and this is not the first row [dataStateChange=2]
			// Only immatureElementModelElement.OriginalElementUuid is changed and this is not the first row [dataStateChange=3]
		case 2, 5:
			// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
			immatureElementModelElements = append(immatureElementModelElements, &previousImmatureElementModelElement)

			// Add immatureElementModelElements to 'immatureTestInstructionMessage.ImmatureSubTestCaseModel' which can be found in map
			immatureTestInstructionMessage, existsInMap = ImmatureTestInstructionMessageMap[previousImmatureElementModelElement.OriginalElementUuid]
			// testInstructionUuid shouldn exist in map. If not so then there is a problem
			if existsInMap == false {
				fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
					"Id": "c2d2448e-9e86-4947-a5f5-7787a72e7ef8",
					"previousImmatureElementModelElement.OriginalElementUuid": previousImmatureElementModelElement.OriginalElementUuid,
				}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
			}

			immatureElementModelMessage.TestCaseModelElements = immatureElementModelElements
			immatureTestInstructionMessage.ImmatureSubTestCaseModel = &immatureElementModelMessage

			// Create fresh versions of variables
			immatureElementModelMessage = fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{}

		// Only immatureElementModelElement.OriginalElementUuid is changed and this is not the first row [dataStateChange=3]
		case 3:
			// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
			immatureElementModelElements = append(immatureElementModelElements, &previousImmatureElementModelElement)

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
				"Id":                          "e464567a-79ab-49e7-9519-ce187b34458d",
				"domainUuid":                  domainUuid,
				"previousDomainUuid":          previousDomainUuid,
				"testInstructionUuid":         testInstructionUuid,
				"previousTestInstructionUuid": previousTestInstructionUuid,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

			// Move actual values into previous-variables
			previousImmatureElementModelElement = immatureElementModelElement

		}

		// Move previous values to current
		previousImmatureElementModelElement = immatureElementModelElement
		previousTempImmatureDomainUuid = tempImmatureElementModelDomainUuid

		// Set to not be first row
		firstRowInSQLRespons = false

	}
	// Handle last row from database

	// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
	immatureElementModelElements = append(immatureElementModelElements, &immatureElementModelElement)

	// Add immatureElementModelElements to 'immatureTestInstructionMessage.ImmatureSubTestCaseModel' which can be found in map
	immatureTestInstructionMessage, existsInMap = ImmatureTestInstructionMessageMap[immatureElementModelElement.OriginalElementUuid]
	// testInstructionUuid shouldn exist in map. If not so then there is a problem
	if existsInMap == false {
		fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
			"Id": "c2d2448e-9e86-4947-a5f5-7787a72e7ef8",
			"immatureElementModelElement.OriginalElementUuid": immatureElementModelElement.OriginalElementUuid,
		}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
	}

	immatureElementModelMessage.TestCaseModelElements = immatureElementModelElements
	immatureTestInstructionMessage.ImmatureSubTestCaseModel = &immatureElementModelMessage
	ImmatureTestInstructionMessageMap[immatureElementModelElement.OriginalElementUuid] = immatureTestInstructionMessage

	// Loop all ImmatureTestInstructionMessage and create gRPC-response
	var allImmatureTestInstructionMessage []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage

	for _, value := range ImmatureTestInstructionMessageMap { // Order not specified
		allImmatureTestInstructionMessage = append(allImmatureTestInstructionMessage, &value)
	}

	cloudDBImmatureTestInstructionItems = allImmatureTestInstructionMessage

	fmt.Println(basicTestInstructionInformationSQLCount)
	fmt.Println(immatureTestInstructionInformationSQLCount)

	// No errors occurred
	return cloudDBImmatureTestInstructionItems, nil

}

*/
// Load TestInstructions for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadClientsImmatureTestInstructionsFromCloudDB(userID string) (cloudDBImmatureTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage, err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "273dceef-7982-4e7d-98db-c132342e530b",
	}).Debug("Entering: loadClientsImmatureTestInstructionsFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "0894a272-3e91-407a-b5a4-1b70f8e00e6b",
		}).Debug("Exiting: loadClientsImmatureTestInstructionsFromCloudDB()")
	}()

	immatureTestInstructionMessageMap := make(map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage)

	err = fenixGuiTestCaseBuilderServerObject.processTestInstructionsBasicTestInstructionInformation(immatureTestInstructionMessageMap)
	if err != nil {
		return nil, err
	}

	err = fenixGuiTestCaseBuilderServerObject.processTestInstructionsImmatureTestInstructionInformation(immatureTestInstructionMessageMap)
	if err != nil {
		return nil, err
	}

	err = fenixGuiTestCaseBuilderServerObject.processTestInstructionsImmatureElementModel(immatureTestInstructionMessageMap)
	if err != nil {
		return nil, err
	}

	// Loop all ImmatureTestInstructionMessage and create gRPC-response
	var allImmatureTestInstructionMessage []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage

	for _, value := range immatureTestInstructionMessageMap { // Order not specified
		allImmatureTestInstructionMessage = append(allImmatureTestInstructionMessage, value)
	}

	cloudDBImmatureTestInstructionItems = allImmatureTestInstructionMessage

	// No errors occurred
	return cloudDBImmatureTestInstructionItems, nil

}

// Load TestInstructionContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadClientsImmatureTestInstructionContainersFromCloudDB(userID string) (cloudDBImmatureTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage, err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "68b965ea-234c-425b-b525-1f8b7154850b",
	}).Debug("Entering: loadClientsImmatureTestInstructionContainersFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "12021bfa-154f-48f2-bd8c-0809e1877fd4",
		}).Debug("Exiting: loadClientsImmatureTestInstructionContainersFromCloudDB()")
	}()

	immatureTestInstructionContainerMessageMap := make(map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage)

	err = fenixGuiTestCaseBuilderServerObject.processTestInstructionContainersBasicTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap)
	if err != nil {
		return nil, err
	}

	err = fenixGuiTestCaseBuilderServerObject.processTestInstructionContainersImmatureTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap)
	if err != nil {
		return nil, err
	}

	err = fenixGuiTestCaseBuilderServerObject.processTestInstructionContainersImmatureElementModel(immatureTestInstructionContainerMessageMap)
	if err != nil {
		return nil, err
	}

	// Loop all ImmatureTestInstructionContainerMessage and create gRPC-response
	var allImmatureTestInstructionContainerMessage []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage

	for _, value := range immatureTestInstructionContainerMessageMap { // Order not specified
		allImmatureTestInstructionContainerMessage = append(allImmatureTestInstructionContainerMessage, value)
	}

	cloudDBImmatureTestInstructionContainerItems = allImmatureTestInstructionContainerMessage

	// No errors occurred
	return cloudDBImmatureTestInstructionContainerItems, nil

}

// Load Pinned TestInstructions for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadClientsPinnedTestInstructionsFromCloudDB(userID string) (availablePinnedPreCreatedTestInstructionContainerMessage []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage, err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "9901525f-a271-4f4f-a798-fea7fdf29dfb",
	}).Debug("Entering: loadClientsPinnedTestInstructionsFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "0f7be73c-4065-4d4a-ae02-40a4d93fc2a3",
		}).Debug("Exiting: loadClientsPinnedTestInstructionsFromCloudDB()")
	}()

	/*
		SELECT PTITIC.*
		FROM "FenixGuiBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers" PTITIC
		WHERE PTITIC."PinnedType" = 1;
	*/

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT PTITIC.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"PinnedTestInstructionsAndPreCreatedTestInstructionContainers\" PTITIC "
	sqlToExecute = sqlToExecute + "WHERE PTITIC.\"PinnedType\" = 1 AND " // 1 = TestInstructions
	sqlToExecute = sqlToExecute + "PTITIC.\"UserId\" = '" + userID + "' "
	sqlToExecute = sqlToExecute + "ORDER BY PTITIC.\"PinnedUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "6be84b22-613f-4d93-afe8-e8ee22826e7b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var tempUserId string
	var tempPinnedType int
	var tempTimeStamp time.Time
	var availablePinnedTestInstructionMessages []fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage

	// Extract data from DB result set
	for rows.Next() {

		// Initiate new fresh variable
		newAvailablePinnedTestInstructionMessage := fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage{}

		err = rows.Scan(

			&tempUserId,
			&newAvailablePinnedTestInstructionMessage.TestInstructionUuid,
			&newAvailablePinnedTestInstructionMessage.TestInstructionName,
			&tempPinnedType,
			&tempTimeStamp,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "e1d695b7-ec8a-4692-9e9a-416869923e82",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// append PinnedTestInstruction to array
		availablePinnedTestInstructionMessages = append(availablePinnedTestInstructionMessages, newAvailablePinnedTestInstructionMessage)

	}

	// Convert to pointer-array that fits gRPC api
	var availablePinnedTestInstructionToSendOvergRPC []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage
	for _, tempAvailablePinnedTestInstructionMessage := range availablePinnedTestInstructionMessages {
		newAvailablePinnedTestInstructionMessage := fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage{}
		newAvailablePinnedTestInstructionMessage = tempAvailablePinnedTestInstructionMessage
		availablePinnedTestInstructionToSendOvergRPC = append(availablePinnedTestInstructionToSendOvergRPC, &newAvailablePinnedTestInstructionMessage)
	}

	return availablePinnedTestInstructionToSendOvergRPC, err
}

// Load Pinned TestInstructionContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadClientsPinnedTestInstructionContainersFromCloudDB(userID string) (availablePinnedPreCreatedTestInstructionContainerContainerMessage []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage, err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "c2decb25-9f53-44c0-be49-88ac5c9cde5d",
	}).Debug("Entering: loadClientsPinnedTestInstructionContainersFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "a9863d2a-4f59-4eef-a939-117bcddea3c4",
		}).Debug("Exiting: loadClientsPinnedTestInstructionContainersFromCloudDB()")
	}()

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT PTITIC.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"PinnedTestInstructionsAndPreCreatedTestInstructionContainers\" PTITIC "
	sqlToExecute = sqlToExecute + "WHERE PTITIC.\"PinnedType\" = 2 AND " // 2 = TestInstructionContainers
	sqlToExecute = sqlToExecute + "PTITIC.\"UserId\" = '" + userID + "' "
	sqlToExecute = sqlToExecute + "ORDER BY PTITIC.\"PinnedUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "6be84b22-613f-4d93-afe8-e8ee22826e7b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var tempUserId string
	var tempPinnedType int
	var tempTimeStamp time.Time
	var availablePinnedTestInstructionContainerMessages []fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage

	// Extract data from DB result set
	for rows.Next() {

		// Initiate new fresh variable
		newAvailablePinnedTestInstructionContainerMessage := fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage{}

		err = rows.Scan(

			&tempUserId,
			&newAvailablePinnedTestInstructionContainerMessage.TestInstructionContainerUuid,
			&newAvailablePinnedTestInstructionContainerMessage.TestInstructionContainerName,
			&tempPinnedType,
			&tempTimeStamp,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "e1d695b7-ec8a-4692-9e9a-416869923e82",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// append PinnedTestInstructionContainer to array
		availablePinnedTestInstructionContainerMessages = append(availablePinnedTestInstructionContainerMessages, newAvailablePinnedTestInstructionContainerMessage)

	}

	// Convert to pointer-array that fits gRPC api
	var availablePinnedTestInstructionContainerToSendOvergRPC []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage
	for _, tempAvailablePinnedTestInstructionContainerMessage := range availablePinnedTestInstructionContainerMessages {
		newAvailablePinnedTestInstructionContainerMessage := fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage{}
		newAvailablePinnedTestInstructionContainerMessage = tempAvailablePinnedTestInstructionContainerMessage
		availablePinnedTestInstructionContainerToSendOvergRPC = append(availablePinnedTestInstructionContainerToSendOvergRPC, &newAvailablePinnedTestInstructionContainerMessage)
	}

	return availablePinnedTestInstructionContainerToSendOvergRPC, err
}

// Load TestInstructions for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadAvailableBondsFromCloudDB() (cloudDBAvailableBondsItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage_ImmatureBondMessage, err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "4b7058fe-c46d-4ab8-8612-895c8e1102a1",
	}).Debug("Entering: loadAvailableBondsFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "665cecf6-a2cc-4c3a-bcb7-9ac1170bd8d3",
		}).Debug("Exiting: loadAvailableBondsFromCloudDB()")
	}()

	//availableBondsAttributes := []fenixTestCaseBuilderServerGrpcApi.BasicBondInformationMessage_VisibleBondAttributesMessage

	availableBondsAttributes, err := fenixGuiTestCaseBuilderServerObject.processVisibleBondAttributesInformation()
	if err != nil {
		return nil, err
	}

	// Loop all Bonds-messages and create gRPC-response
	for _, visibleBondAttributesMessage := range availableBondsAttributes {

		// Deep copy of values
		tempBondAttributesMessage := fenixTestCaseBuilderServerGrpcApi.BasicBondInformationMessage_VisibleBondAttributesMessage{
			BondUuid:                 visibleBondAttributesMessage.BondUuid,
			BondName:                 visibleBondAttributesMessage.BondName,
			BondDescription:          visibleBondAttributesMessage.BondDescription,
			BondMouseOverText:        visibleBondAttributesMessage.BondMouseOverText,
			Deprecated:               visibleBondAttributesMessage.Deprecated,
			Enabled:                  visibleBondAttributesMessage.Enabled,
			Visible:                  visibleBondAttributesMessage.Visible,
			BondColor:                visibleBondAttributesMessage.BondColor,
			CanBeDeleted:             visibleBondAttributesMessage.CanBeDeleted,
			CanBeSwappedOut:          visibleBondAttributesMessage.CanBeSwappedOut,
			UpdatedTimeStamp:         visibleBondAttributesMessage.UpdatedTimeStamp,
			TestCaseModelElementType: visibleBondAttributesMessage.TestCaseModelElementType,
			ShowBondAttributes:       visibleBondAttributesMessage.ShowBondAttributes,
			TCRuleDeletion:           visibleBondAttributesMessage.TCRuleDeletion,
			TCRuleSwap:               visibleBondAttributesMessage.TCRuleSwap,
		}

		basicBondInformationMessage := fenixTestCaseBuilderServerGrpcApi.BasicBondInformationMessage{
			VisibleBondAttributes: &tempBondAttributesMessage}

		immatureBondsMessage_ImmatureBondMessage_B0_BOND := fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage_ImmatureBondMessage{
			BasicBondInformation: &basicBondInformationMessage}

		cloudDBAvailableBondsItems = append(cloudDBAvailableBondsItems, &immatureBondsMessage_ImmatureBondMessage_B0_BOND)
	}

	// No errors occurred
	return cloudDBAvailableBondsItems, nil

}
