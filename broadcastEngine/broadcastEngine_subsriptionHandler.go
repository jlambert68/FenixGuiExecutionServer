package broadcastEngine

import "strconv"

// InitiateSubscriptionHandler
// Initiate the handler that takes care of the information about which TesterGui:s that subscribe to what TestCaseExecution, regarding status updates
func InitiateSubscriptionHandler() {

	// Initiate map for holding information about the data needed to route ExecutionStatuses to correct TesterGui
	TestCaseExecutionsSubscriptionChannelInformationMap = make(map[ApplicationRunTimeUuidType]*TestCaseExecutionsSubscriptionChannelInformationStruct)

	// Initiate map that holds information about who is subscribing to a certain TestCaseExecution
	TestCaseExecutionsSubscriptionsMap = make(map[TestCaseExecutionsSubscriptionsMapKeyType]*[]ApplicationRunTimeUuidType)
}

// AddSubscriptionForTestCaseExecutionToTesterGui
// Create Subscription on this TestCaseExecution and this TestGui
func AddSubscriptionForTestCaseExecutionToTesterGui(applicationRunTimeUuid ApplicationRunTimeUuidType, testCaseExecutionUuid TestCaseExecutionUuidType, testCaseExecutionUuidVersion TestCaseExecutionUuidVersionType) {

	var allApplicationRunTimeUuids *[]ApplicationRunTimeUuidType
	var existInMap bool

	// Create Key used for 'TestCaseExecutionsSubscriptionsMap'
	var testCaseExecutionsSubscriptionsMapKey TestCaseExecutionsSubscriptionsMapKeyType
	testCaseExecutionsSubscriptionsMapKey = TestCaseExecutionsSubscriptionsMapKeyType(string(testCaseExecutionUuid) + strconv.Itoa(int(testCaseExecutionUuidVersion)))

	// Check if TesterGui already exist in Subscription-map for incoming 'TestCaseExecutionUuid'
	allApplicationRunTimeUuids, existInMap = TestCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey]

	// TestCaseExecution doesn't have any subscriptions yet, so just add it
	if existInMap == false {

		*allApplicationRunTimeUuids = append(*allApplicationRunTimeUuids, applicationRunTimeUuid)
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

func WhoIsSubscribingToTestCaseExecution() (messageToTesterGuiForwardChannels []*MessageToTesterGuiForwardChannelType) {

}
