package testerGuiOwnerEngine

import (
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
	oldSlice []*guiExecutionServerStartUpOrderStruct,
	v *guiExecutionServerStartUpOrderStruct) (
	newSlice []*guiExecutionServerStartUpOrderStruct) {

	// Lock Slice for update
	guiExecutonServerSliceLoadAndSaveMutex.Lock()

	// Find the index were item should be inserted
	var index int
	index = sort.Search(len(oldSlice), func(i int) bool {

		return oldSlice[i].applicationRunTimeStartUpTime.After(v.applicationRunTimeStartUpTime)
	})

	// Insert Item at index
	newSlice = insertAtIndex(oldSlice, index, v)

	//UnLock Slice
	guiExecutonServerSliceLoadAndSaveMutex.Unlock()

	return newSlice
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
