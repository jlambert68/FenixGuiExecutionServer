package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"github.com/sirupsen/logrus"
	"sort"
	"sync"
)

// Used to lock slice when reading and writing the slice
var guiExecutonServerSliceLoadAndSaveMutex = &sync.RWMutex{}

// Insert GuiExecutionServer into slice with all known GuiExecutionsServers that have StartUp-time
// after this GuiExecutionServer's StartUp-time.
// GuiExecutionServers are ordered in StartUp-time-order, ascending, with this GuiExecutionServer as
// the first item
func insertGuiExecutionServerIntoTimeOrderedSlice(
	elementToInsert *guiExecutionServerStartUpOrderStruct) {

	// Lock Slice for update
	guiExecutonServerSliceLoadAndSaveMutex.Lock()

	// If the slice is empty, then do simple insert
	if len(guiExecutionServerStartUpOrder) == 0 {
		guiExecutionServerStartUpOrder = append(guiExecutionServerStartUpOrder, elementToInsert)

		return
	}

	// Do not insert elements with a TimeStamp  that is before to current GuiExecutionsServers StartUp-TimeStamp
	// Current GuiExecutionsServer is always store first in the slice
	if guiExecutionServerStartUpOrder[0].applicationRunTimeStartUpTime.
		After(elementToInsert.applicationRunTimeStartUpTime) {

		return
	}

	// If the slice only has one element, then do simple insert
	if len(guiExecutionServerStartUpOrder) == 1 {
		guiExecutionServerStartUpOrder = append(guiExecutionServerStartUpOrder, elementToInsert)

		return
	}

	// Find the index were item should be inserted
	var index int
	index = sort.Search(len(guiExecutionServerStartUpOrder), func(i int) bool {

		return guiExecutionServerStartUpOrder[i].applicationRunTimeStartUpTime.After(elementToInsert.applicationRunTimeStartUpTime)
	})

	// Insert Item at index
	guiExecutionServerStartUpOrder = insertAtIndex(guiExecutionServerStartUpOrder, index, elementToInsert)

	//UnLock Slice
	guiExecutonServerSliceLoadAndSaveMutex.Unlock()

}

// Helper function that inserts element into slice at index and returns a new slice.
func insertAtIndex(
	oldSlice []*guiExecutionServerStartUpOrderStruct,
	index int,
	elementToInsert *guiExecutionServerStartUpOrderStruct) (
	newSlice []*guiExecutionServerStartUpOrderStruct) {

	if index == len(oldSlice) {
		// Insert at end is the easy case.
		newSlice = append(oldSlice, elementToInsert)
		return newSlice
	}

	// Make space for the inserted element by shifting
	// values at the insertion index up one index. The call
	// to append does not allocate memory when cap(data) is
	// greater than len(data).
	newSlice = append(oldSlice[:index+1], oldSlice[index:]...)

	// Insert the new element.
	newSlice[index] = elementToInsert

	// Return the updated slice.
	return newSlice
}

// Verify that 'GuiExecutionServer' exists in the time ordered slice in the correct position
func verifyThatGuiExecutionServerExistsInTimeOrderedSlice(
	guiExecutionServerToBeVerifiedForExistence *guiExecutionServerStartUpOrderStruct) (
	existsInTimeOrderedSliceInCorrectPosition bool) {

	// Lock Slice for update
	guiExecutonServerSliceLoadAndSaveMutex.Lock()

	// Always unlock slice when leaving
	defer func() {
		//UnLock Slice
		guiExecutonServerSliceLoadAndSaveMutex.Unlock()
	}()

	// When Sending GuiExecutionServer's StartupTimeStamp is before current GuiExecutionServer's StartupTimeStamp then just exit
	if guiExecutionServerToBeVerifiedForExistence.applicationRunTimeStartUpTime.
		Before(common_config.ApplicationRunTimeStartUpTime) == true {
		return true
	}

	// Number of elements in slice
	var numberOfElementsInSlice int
	numberOfElementsInSlice = len(guiExecutionServerStartUpOrder)

	existsInTimeOrderedSliceInCorrectPosition = false
	// Loop the slice and verify existence and position
	for sliceIndex, tempGuiExecutionServerStartUpOrderElement := range guiExecutionServerStartUpOrder {

		// First element in slice is always this GuiExecutionServer
		if sliceIndex == 0 {
			continue
		}

		// If this is the one we are looking for
		if tempGuiExecutionServerStartUpOrderElement.applicationRunTimeUuid ==
			guiExecutionServerToBeVerifiedForExistence.applicationRunTimeUuid {

			// Verify that it's StartUpTime is not after the previous one's StartUpTime
			if tempGuiExecutionServerStartUpOrderElement.applicationRunTimeStartUpTime.
				After(guiExecutionServerStartUpOrder[sliceIndex-1].applicationRunTimeStartUpTime) == true {

				// This shouldn't be like this
				common_config.Logger.WithFields(logrus.Fields{
					"id":                             "84288fef-62be-4355-b8fc-36337e9e41a9",
					"guiExecutionServerStartUpOrder": guiExecutionServerStartUpOrder,
					"sliceIndex":                     sliceIndex,
				}).Error("The StartUpTime is not after the previous one's StartUpTime. This shouldn't happen")

				return false
			}

			// Check if this is the last position in slice
			if sliceIndex != numberOfElementsInSlice-1 {

				return true
			} else {
				// Verify that it's StartUpTime is before the next one's StartUpTime
				if tempGuiExecutionServerStartUpOrderElement.applicationRunTimeStartUpTime.
					Before(guiExecutionServerStartUpOrder[sliceIndex+1].applicationRunTimeStartUpTime) == true {

					return true
				} else {
					// This shouldn't be like this
					common_config.Logger.WithFields(logrus.Fields{
						"id":                             "736f397e-41a3-454b-af03-ae7480c41527",
						"guiExecutionServerStartUpOrder": guiExecutionServerStartUpOrder,
						"sliceIndex":                     sliceIndex,
					}).Error("The StartUpTime is not before the next one's StartUpTime. This shouldn't happen")

					return false
				}
			}
		}

		break
	}

	return existsInTimeOrderedSliceInCorrectPosition
}

// Removes the specified GuiExecutionServer from time sorted slice of all GuiExecutionsServers
func removeGuiExecutionServerFromSlice(applicationRunTimeUuid string) {

	// Lock Slice for update
	guiExecutonServerSliceLoadAndSaveMutex.Lock()

	// Find the index of the item to be removed
	var index int
	index = sort.Search(len(guiExecutionServerStartUpOrder), func(i int) bool {

		return guiExecutionServerStartUpOrder[i].applicationRunTimeUuid == applicationRunTimeUuid
	})

	// Remove the item at index
	guiExecutionServerStartUpOrder = removeAtIndex(guiExecutionServerStartUpOrder, index)

	//UnLock Slice
	guiExecutonServerSliceLoadAndSaveMutex.Unlock()

}

// Helper function which removes the element, at index position, from slice
func removeAtIndex(
	oldSlice []*guiExecutionServerStartUpOrderStruct,
	index int) (
	newSlice []*guiExecutionServerStartUpOrderStruct) {

	// Remove the element at index from slice.
	copy(oldSlice[index:], oldSlice[index+1:]) // Shift a[i+1:] left one index.

	oldSlice[len(oldSlice)-1] = nil       // Erase last element (write zero value).
	newSlice = oldSlice[:len(oldSlice)-1] // Truncate slice.

	return newSlice

}
