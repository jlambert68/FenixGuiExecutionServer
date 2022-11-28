package common_config

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	fenixTestDataSyncServerGrpcApi "github.com/jlambert68/FenixGrpcApi/Fenix/fenixTestDataSyncServerGrpcApi/go_grpc_api"
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

// ********************************************************************************************************************
// Check if Calling Client is using correct proto-file version
func IsClientUsingCorrectTestDataProtoFileVersion(callingClientUuid string, usedProtoFileVersion fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum) (returnMessage *fenixExecutionServerGuiGrpcApi.AckNackResponse) {

	var clientUseCorrectProtoFileVersion bool
	var protoFileExpected fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum
	var protoFileUsed fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum

	protoFileUsed = usedProtoFileVersion
	protoFileExpected = fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum(GetHighestFenixGuiExecutionServerProtoFileVersion())

	// Check if correct proto files is used
	if protoFileExpected == protoFileUsed {
		clientUseCorrectProtoFileVersion = true
	} else {
		clientUseCorrectProtoFileVersion = false
	}

	// Check if Client is using correct proto files version
	if clientUseCorrectProtoFileVersion == false {
		// Not correct proto-file version is used

		// Set Error codes to return message
		var errorCodes []fenixExecutionServerGuiGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionServerGuiGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionServerGuiGrpcApi.ErrorCodesEnum_ERROR_WRONG_PROTO_FILE_VERSION
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixExecutionServerGuiGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Wrong proto file used. Expected: '" + protoFileExpected.String() + "', but got: '" + protoFileUsed.String() + "'",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: protoFileExpected,
		}

		Logger.WithFields(logrus.Fields{
			"id": "513dd8fb-a0bb-4738-9a0b-b7eaf7bb8adb",
		}).Debug("Wrong proto file used. Expected: '" + protoFileExpected.String() + "', but got: '" + protoFileUsed.String() + "' for Client: " + callingClientUuid)

		return returnMessage

	} else {
		return nil
	}

}

// ********************************************************************************************************************
// Get the highest FenixGuiExecutionServerProtoFileVersionEnumeration
func GetHighestFenixGuiExecutionServerProtoFileVersion() int32 {

	// Check if there already is a 'highestFenixGuiExecutionServerProtoFileVersion' saved, if so use that one
	if highestFenixGuiExecutionServerProtoFileVersion != -1 {
		return highestFenixGuiExecutionServerProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixExecutionServerGuiGrpcApi.CurrentFenixExecutionGuiProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestFenixGuiExecutionServerProtoFileVersion = maxValue

	return highestFenixGuiExecutionServerProtoFileVersion
}

// Exctract Values, and create, for TestDataHeaderItemMessageHash
func CreateTestDataHeaderItemMessageHash(testDataHeaderItemMessage *fenixTestDataSyncServerGrpcApi.TestDataHeaderItemMessage) (testDataHeaderItemMessageHash string) {

	var valuesToHash []string
	var valueToHash string

	// Extract and add values to array
	// HeaderLabel
	valueToHash = testDataHeaderItemMessage.HeaderLabel
	valuesToHash = append(valuesToHash, valueToHash)

	// HeaderShouldBeUsedForTestDataFilter as 'true' or 'false'
	if testDataHeaderItemMessage.HeaderShouldBeUsedForTestDataFilter == false {
		valuesToHash = append(valuesToHash, "false")
	} else {
		valuesToHash = append(valuesToHash, "true")
	}

	// HeaderIsMandatoryInTestDataFilter as 'true' or 'false'
	if testDataHeaderItemMessage.HeaderIsMandatoryInTestDataFilter == false {
		valuesToHash = append(valuesToHash, "false")
	} else {
		valuesToHash = append(valuesToHash, "true")
	}

	// HeaderSelectionType
	valueToHash = testDataHeaderItemMessage.HeaderSelectionType.String()
	valuesToHash = append(valuesToHash, valueToHash)

	// HeaderFilterValues - An array thar is added
	for _, headerFilterValue := range testDataHeaderItemMessage.HeaderFilterValues {
		headerFilterValueToAdd := headerFilterValue.String()
		valuesToHash = append(valuesToHash, headerFilterValueToAdd)
	}

	// Hash all values in the array
	testDataHeaderItemMessageHash = HashValues(valuesToHash, true)

	return testDataHeaderItemMessageHash
}

// Hash a single value
func HashSingleValue(valueToHash string) (hashValue string) {

	hash := sha256.New()
	hash.Write([]byte(valueToHash))
	hashValue = hex.EncodeToString(hash.Sum(nil))

	return hashValue

}

// GenerateDatetimeTimeStampForDB
// Generate DataBaseTimeStamp, eg '2022-02-08 17:35:04.000000'
func GenerateDatetimeTimeStampForDB() (currentTimeStampAsString string) {

	timeStampLayOut := "2006-01-02 15:04:05.000000" //milliseconds
	currentTimeStamp := time.Now()
	currentTimeStampAsString = currentTimeStamp.Format(timeStampLayOut)

	return currentTimeStampAsString
}

// ConvertGrpcTimeStampToStringForDB
// Convert a gRPC-timestamp into a string that can be used to store in the database
func ConvertGrpcTimeStampToStringForDB(grpcTimeStamp *timestamppb.Timestamp) (grpcTimeStampAsTimeStampAsString string) {
	grpcTimeStampAsTimeStamp := grpcTimeStamp.AsTime()

	timeStampLayOut := "2006-01-02 15:04:05.000000" //milliseconds

	grpcTimeStampAsTimeStampAsString = grpcTimeStampAsTimeStamp.Format(timeStampLayOut)

	return grpcTimeStampAsTimeStampAsString
}

// Extracts 'ParserLayout' from the TimeStamp(as string)
func GenerateTimeStampParserLayout(timeStampAsString string) (parserLayout string, err error) {
	// "2006-01-02 15:04:05.999999999 -0700 MST"

	var timeStampParts []string
	var timeParts []string
	var numberOfDecimals int

	// Split TimeStamp into separate parts
	timeStampParts = strings.Split(timeStampAsString, " ")

	// Validate that first part is a date with the following form '2006-01-02'
	if len(timeStampParts[0]) != 10 {

		Logger.WithFields(logrus.Fields{
			"Id":                "ffbf0682-ebc7-4e27-8ad1-0e5005fbc364",
			"timeStampAsString": timeStampAsString,
			"timeStampParts[0]": timeStampParts[0],
		}).Error("Date part has not the correct form, '2006-01-02'")

		err = errors.New(fmt.Sprintf("Date part, '%s' has not the correct form, '2006-01-02'", timeStampParts[0]))

		return "", err

	}

	// Add Date to Parser Layout
	parserLayout = "2006-01-02"

	// Add Time to Parser Layout
	parserLayout = parserLayout + " 15:04:05."

	// Split time into time and decimals
	timeParts = strings.Split(timeStampParts[1], ".")

	// Get number of decimals
	numberOfDecimals = len(timeParts[1])

	// Add Decimals to Parser Layout
	parserLayout = parserLayout + strings.Repeat("9", numberOfDecimals)

	// Add time zone if that information exists
	if len(timeStampParts) > 3 {
		parserLayout = parserLayout + " -0700 MST"
	}

	return parserLayout, err
}
