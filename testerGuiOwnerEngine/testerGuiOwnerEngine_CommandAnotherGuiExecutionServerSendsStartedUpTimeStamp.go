package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"github.com/sirupsen/logrus"
)

// Process the actual command 'ChannelCommand_AnotherGuiExecutionServerSendsStartedUpTimeStamp'
func commandAnotherGuiExecutionServerSendsStartedUpTimeStamp(
	tempGuiExecutionServerStartedUpTimeStampRefresher *common_config.GuiExecutionServerStartedUpTimeStampRefresherStruct) {

	// Verify that it is not this GuiExecutionServer in the message, if so then just exit
	if tempGuiExecutionServerStartedUpTimeStampRefresher.GuiExecutionServerApplicationId == common_config.ApplicationRunTimeUuid {

		return
	}

	// Verify that 'GuiExecutionServer' exists in the time ordered slice in the correct position
	// Try to insert GuiExecutionServer-information into slice with GuiExecutionServers
	// Logic if it can be inserted is handled by the function itself
	var guiExecutionServerToBeVerifiedForExistence *guiExecutionServerStartUpOrderStruct
	guiExecutionServerToBeVerifiedForExistence = &guiExecutionServerStartUpOrderStruct{
		applicationRunTimeUuid:        tempGuiExecutionServerStartedUpTimeStampRefresher.GuiExecutionServerApplicationId,
		applicationRunTimeStartUpTime: tempGuiExecutionServerStartedUpTimeStampRefresher.MessageTimeStamp,
	}

	var existsInTimeOrderedSliceInCorrectPosition bool
	existsInTimeOrderedSliceInCorrectPosition = verifyThatGuiExecutionServerExistsInTimeOrderedSlice(guiExecutionServerToBeVerifiedForExistence)

	if existsInTimeOrderedSliceInCorrectPosition == false {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                             "75ce2c6a-c714-41fc-884e-53d93fa0fb97",
			"guiExecutionServerStartUpOrder": guiExecutionServerStartUpOrder,
		}).Error("The GuiExecutionServer did not already exist in the time ordered slice. This shouldn't happen")
	}

}
