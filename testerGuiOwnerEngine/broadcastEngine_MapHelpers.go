package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"fmt"
	"sync"
)

var subscriptionsMapMutex = &sync.RWMutex{}

// Load Subscription from the Subscriptions-Map
func loadFromTestCaseExecutionsSubscriptionsMap(
	testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType) (
	guiExecutionServerResponsibilityStruct *common_config.GuiExecutionServerResponsibilityStruct,
	existInMap bool) {

	// Lock Map for Reading
	subscriptionsMapMutex.RLock()

	// Read Map
	guiExecutionServerResponsibilityStruct, existInMap = testCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey]

	if existInMap == false {
		for x, y := range testCaseExecutionsSubscriptionsMap {
			yy := *y
			fmt.Println(x, yy)
		}
	}

	//UnLock Map
	subscriptionsMapMutex.RUnlock()

	return guiExecutionServerResponsibilityStruct, existInMap
}

// Save Subscription to the Subscriptions-Map
func saveToTestCaseExecutionsSubscriptionsMap(
	testCaseExecutionsSubscriptionsMapKey testCaseExecutionsSubscriptionsMapKeyType,
	guiExecutionServerResponsibilityStruct *common_config.GuiExecutionServerResponsibilityStruct) {

	// Lock Map for Writing
	subscriptionsMapMutex.Lock()

	// Save to Subscription-Map
	testCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey] = guiExecutionServerResponsibilityStruct

	//UnLock Map
	subscriptionsMapMutex.Unlock()

}
