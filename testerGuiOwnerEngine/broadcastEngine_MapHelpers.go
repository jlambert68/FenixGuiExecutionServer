package testerGuiOwnerEngine

import (
	"fmt"
	"sync"
)

var subscriptionsMapMutex = &sync.RWMutex{}

// Load Subscription from the Subscriptions-Map
func loadFromTestCaseExecutionsSubscriptionsMap(
	testCaseExecutionsSubscriptionsMapKey TestCaseExecutionsSubscriptionsMapKeyType) (
	applicationsRunTimeUuid *[]ApplicationRunTimeUuidType,
	existInMap bool) {

	// Lock Map for Reading
	subscriptionsMapMutex.RLock()

	// Read Map
	applicationsRunTimeUuid, existInMap = TestCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey]

	if existInMap == false {
		for x, y := range TestCaseExecutionsSubscriptionsMap {
			yy := *y
			fmt.Println(x, yy)
		}
	}

	//UnLock Map
	subscriptionsMapMutex.RUnlock()

	return applicationsRunTimeUuid, existInMap
}

// Save Subscription to the Subscriptions-Map
func saveToTestCaseExecutionsSubscriptionsMap(
	testCaseExecutionsSubscriptionsMapKey TestCaseExecutionsSubscriptionsMapKeyType,
	applicationRunTimeUuidSliceReference *[]ApplicationRunTimeUuidType) {

	// Lock Map for Writing
	subscriptionsMapMutex.Lock()

	// Save to Subscription-Map
	TestCaseExecutionsSubscriptionsMap[testCaseExecutionsSubscriptionsMapKey] = applicationRunTimeUuidSliceReference

	//UnLock Map
	subscriptionsMapMutex.Unlock()

}
