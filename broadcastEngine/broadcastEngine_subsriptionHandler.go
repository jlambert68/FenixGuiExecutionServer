package broadcastEngine

import (
	"FenixGuiExecutionServer/common_config"
	"github.com/sirupsen/logrus"
	"strconv"
)

// InitiateSubscriptionHandler
// Initiate the handler that takes care of the information about which TesterGui:s that subscribe to what TestCaseExecution, regarding status updates
func InitiateSubscriptionHandler() {

	// Initiate map for holding information about the data needed to route ExecutionStatuses to correct TesterGui
	TestCaseExecutionsSubscriptionChannelInformationMap = make(map[ApplicationRunTimeUuidType]*TestCaseExecutionsSubscriptionChannelInformationStruct)

	// Initiate map that holds information about who is subscribing to a certain TestCaseExecution
	TestCaseExecutionsSubscriptionsMap = make(map[TestCaseExecutionsSubscriptionsMapKeyType]*[]ApplicationRunTimeUuidType)
}

// AddSubscriptionForTestCaseExecutionToTesterGui
// Create Subscription on this TestCaseExecution for this TestGui
func AddSubscriptionForTestCaseExecutionToTesterGui(
	applicationRunTimeUuid ApplicationRunTimeUuidType,
	testCaseExecutionUuid TestCaseExecutionUuidType,
	testCaseExecutionUuidVersion TestCaseExecutionUuidVersionType) {

	//var allApplicationRunTimeUuidsReference *[]ApplicationRunTimeUuidType
	var allApplicationRunTimeUuids *[]ApplicationRunTimeUuidType

	// Create Key used for 'TestCaseExecutionsSubscriptionsMap'
	var testCaseExecutionsSubscriptionsMapKey TestCaseExecutionsSubscriptionsMapKeyType
	testCaseExecutionsSubscriptionsMapKey = TestCaseExecutionsSubscriptionsMapKeyType(string(testCaseExecutionUuid) +
		strconv.Itoa(int(testCaseExecutionUuidVersion)))

	// Check if TesterGui already exist in Subscription-map for incoming 'TestCaseExecutionUuid'
	//allApplicationRunTimeUuids, existInMap = TestCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey]
	allApplicationRunTimeUuids, _ = loadFromTestCaseExecutionsSubscriptionsMap(testCaseExecutionsSubscriptionsMapKey)

	// Nothing in subscription-map then initiate slice and store it in Map
	if allApplicationRunTimeUuids == nil {
		var tempAllApplicationRunTimeUuids []ApplicationRunTimeUuidType

		// Add new 'ApplicationRunTimeUuid' to slice
		tempAllApplicationRunTimeUuids = append(tempAllApplicationRunTimeUuids, applicationRunTimeUuid)

		// Add it to map
		//TestCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey] = &tempAllApplicationRunTimeUuids
		saveToTestCaseExecutionsSubscriptionsMap(testCaseExecutionsSubscriptionsMapKey, &tempAllApplicationRunTimeUuids)

		//allApplicationRunTimeUuids = &tempAllApplicationRunTimeUuids

	} else {

		// Loop all 'ApplicationRunTimeUuid' to verify if incoming 'applicationRunTimeUuid' exists in slice
		var foundApplicationRunTimeUuidInSlice bool
		for _, tempApplicationRunTimeUuid := range *allApplicationRunTimeUuids {
			if tempApplicationRunTimeUuid == applicationRunTimeUuid {
				// 'applicationRunTimeUuid' existed in slice
				foundApplicationRunTimeUuidInSlice = true
				break
			}
		}

		// if 'applicationRunTimeUuid' didn't exist in slice then add it to the slice
		if foundApplicationRunTimeUuidInSlice == false {
			*allApplicationRunTimeUuids = append(*allApplicationRunTimeUuids, applicationRunTimeUuid)
		}
	}
}

// Generates a slice with pointers to all 'MessageToTesterGuiForwardChannel' for
// 'TestCaseExecutionUuidTestCaseExecutionVersion' contains  ('TestCaseExecutionUuid' + 'TestCaseExecutionVersion')
func whoIsSubscribingToTestCaseExecution(testCaseExecutionUuidTestCaseExecutionVersion string) (messageToTesterGuiForwardChannels []*MessageToTesterGuiForwardChannelType) {

	var applicationsRunTimeUuidSlice *[]ApplicationRunTimeUuidType
	var existInMap bool

	// Extract slice of Applications that subscribes to combination of ('TestCaseExecutionUuid' + 'TestCaseExecutionVersion')
	//applicationsRunTimeUuidSlice, existInMap = TestCaseExecutionsSubscriptionsMap[TestCaseExecutionsSubscriptionsMapKeyType(testCaseExecutionUuidTestCaseExecutionVersion)]
	applicationsRunTimeUuidSlice, existInMap = loadFromTestCaseExecutionsSubscriptionsMap(TestCaseExecutionsSubscriptionsMapKeyType(testCaseExecutionUuidTestCaseExecutionVersion))

	if existInMap == false {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "0ed78746-1ff8-4261-9657-023048d8db84",
			"testCaseExecutionUuidTestCaseExecutionVersion": testCaseExecutionUuidTestCaseExecutionVersion,
		}).Error("No TesterGui is subscribing to the  combination of ('TestCaseExecutionUuid' + 'TestCaseExecutionVersion') ")

		return messageToTesterGuiForwardChannels
	}

	// Loop Subscribing 'applicationsRunTimeUuidSlice' to get their channel-reference
	var tempApplicationRunTimeUuid ApplicationRunTimeUuidType
	for _, tempApplicationRunTimeUuid = range *applicationsRunTimeUuidSlice {

		// Get Channel-reference based on 'tempApplicationRunTimeUuid'
		var tempTestCaseExecutionsSubscriptionChannelInformation *TestCaseExecutionsSubscriptionChannelInformationStruct
		tempTestCaseExecutionsSubscriptionChannelInformation, existInMap = TestCaseExecutionsSubscriptionChannelInformationMap[tempApplicationRunTimeUuid]

		// If 'tempApplicationRunTimeUuid' doesn't exit the most local reason is that TestGui hasn't yet open up gRPC-stream for messages. But it could be an error
		if existInMap == false {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                         "81fb8977-2ff4-4cfa-84e5-ba9c2f03485e",
				"tempApplicationRunTimeUuid": tempApplicationRunTimeUuid,
			}).Info("Couldn't find Channel data based on 'ApplicationRunTimeUuid'. Could be an error or that TesterGui hasn't yet open up gRPC-stream for Messages")

			return messageToTesterGuiForwardChannels
		}

		// Add Channel-reference to return slice
		messageToTesterGuiForwardChannels = append(messageToTesterGuiForwardChannels, tempTestCaseExecutionsSubscriptionChannelInformation.MessageToTesterGuiForwardChannel)

	}

	return messageToTesterGuiForwardChannels

}
