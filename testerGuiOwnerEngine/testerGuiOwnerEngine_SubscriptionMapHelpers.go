package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"sync"
)

// Used to lock map when reading and writing the map
var subscriptionsMapLoadAndSaveMutex = &sync.RWMutex{}

// Load Subscription from the Subscriptions-Map
func loadFromTestCaseExecutionsSubscriptionFromMap(
	testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType) (
	guiExecutionServerResponsibility *common_config.GuiExecutionServerResponsibilityStruct,
	existInMap bool) {

	// Lock Map for Reading
	subscriptionsMapLoadAndSaveMutex.RLock()

	// Read Map
	guiExecutionServerResponsibility, existInMap = testCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey]

	/*
		if existInMap == false {
			for x, y := range testCaseExecutionsSubscriptionsMap {
				yy := *y
				fmt.Println(x, yy)
			}
		}
	*/

	//UnLock Map
	subscriptionsMapLoadAndSaveMutex.RUnlock()

	return guiExecutionServerResponsibility, existInMap
}

// Save Subscription to the Subscriptions-Map
func saveToTestCaseExecutionsSubscriptionToMap(
	testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType,
	guiExecutionServerResponsibility *common_config.GuiExecutionServerResponsibilityStruct) {

	// Lock Map for Writing
	subscriptionsMapLoadAndSaveMutex.Lock()

	// Save to Subscription-Map
	testCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey] = guiExecutionServerResponsibility

	//UnLock Map
	subscriptionsMapLoadAndSaveMutex.Unlock()

}

// Delete Subscription from the Subscriptions-Map
func deleteTestCaseExecutionsSubscriptionFromMap(
	testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType) {

	// Lock Map for Deleting
	subscriptionsMapLoadAndSaveMutex.Lock()

	// Delete from Subscription-Map
	delete(testCaseExecutionsSubscriptionsMap, testCaseExecutionsSubscriptionsMapKey)

	//UnLock Map
	subscriptionsMapLoadAndSaveMutex.Unlock()

}

// Delete All Subscription, for specific TesterGui, from the Subscriptions-Map
func deleteTesterGuiFromTestCaseExecutionsSubscriptionFromMap(testerGuiApplicationUuid string) {

	// Lock Map for Deleting
	subscriptionsMapLoadAndSaveMutex.Lock()

	var testCaseExecutionsSubscriptionsMapKeys []testCaseExecutionsSubscriptionsMapKeyType

	// List all 'testCaseExecutionsSubscriptionsMapKeys' that match TesterGuiApplicationId
	for testCaseExecutionsSubscriptionsMapKey, guiExecutionServerResponsibility := range testCaseExecutionsSubscriptionsMap {

		if guiExecutionServerResponsibility.TesterGuiApplicationId == testerGuiApplicationUuid {
			testCaseExecutionsSubscriptionsMapKeys = append(testCaseExecutionsSubscriptionsMapKeys, testCaseExecutionsSubscriptionsMapKey)
		}

	}

	// Loop the 'testCaseExecutionsSubscriptionsMapKeys' and delete from Subscription-Map
	for _, testCaseExecutionsSubscriptionsMapKey := range testCaseExecutionsSubscriptionsMapKeys {
		delete(testCaseExecutionsSubscriptionsMap, testCaseExecutionsSubscriptionsMapKey)
	}

	//UnLock Map
	subscriptionsMapLoadAndSaveMutex.Unlock()

}

// List all Subscriptions from the Subscriptions-Map
func listAllTestCaseExecutionsSubscriptionsFromMap() (
	guiExecutionServerResponsibilities []common_config.GuiExecutionServerResponsibilityStruct) {

	// Lock Map for Reading
	subscriptionsMapLoadAndSaveMutex.RLock()

	// Get all values from map
	for _, guiExecutionServerResponsibility := range testCaseExecutionsSubscriptionsMap {

		guiExecutionServerResponsibilities = append(guiExecutionServerResponsibilities, *guiExecutionServerResponsibility)
	}
	//UnLock Map
	subscriptionsMapLoadAndSaveMutex.RUnlock()

	return guiExecutionServerResponsibilities
}

// Save Subscription to the Subscriptions-Map, if it is not already in the map
func saveToTestCaseExecutionsSubscriptionToMapIfMissingInMap(
	testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType,
	guiExecutionServerResponsibility *common_config.GuiExecutionServerResponsibilityStruct) {

	var existInMap bool

	// Lock Map for Reading
	subscriptionsMapLoadAndSaveMutex.Lock()

	_, existInMap = loadFromTestCaseExecutionsSubscriptionFromMap(testCaseExecutionsSubscriptionsMapKey)

	// If missing in map then save to map
	if existInMap == false {
		saveToTestCaseExecutionsSubscriptionToMap(
			testCaseExecutionsSubscriptionsMapKey,
			guiExecutionServerResponsibility)
	}

	//UnLock Map
	subscriptionsMapLoadAndSaveMutex.Unlock()

}
