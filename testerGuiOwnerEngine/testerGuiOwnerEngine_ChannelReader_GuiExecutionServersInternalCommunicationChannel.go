package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"github.com/sirupsen/logrus"
)

// Channel reader which is used for reading out commands to TesterGuiOwnerEngine
func startTesterGuiOwnerEngineChannelReader() {

	var incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct
	var channelSize int

	// If the channel is not initialized then do that
	if common_config.TesterGuiOwnerEngineChannelEngineCommandChannel == nil {
		common_config.TesterGuiOwnerEngineChannelEngineCommandChannel = make(
			common_config.TesterGuiOwnerEngineChannelEngineType, common_config.TesterGuiOwnerEngineChannelSize)
	}

	// Inform other running GuiExecutionServers that this server is starting up
	informOtherGuiExecutionServersThatThisGuiExecutionServerIsStartingUp()

	for {
		// Wait for incoming command over channel
		incomingTesterGuiOwnerEngineChannelCommand = <-common_config.TesterGuiOwnerEngineChannelEngineCommandChannel

		common_config.Logger.WithFields(logrus.Fields{
			"Id": "a2809c91-87bc-44fc-894b-c8cdd73b521f",
			"incomingTesterGuiOwnerEngineChannelCommand": incomingTesterGuiOwnerEngineChannelCommand,
		}).Debug("Message received on 'TesterGuiOwnerEngineChannel'")

		// If size of Channel > 'TesterGuiOwnerEngineChannelWarningLevel' then log Warning message
		channelSize = len(common_config.TesterGuiOwnerEngineChannelEngineCommandChannel)
		if channelSize > common_config.TesterGuiOwnerEngineChannelWarningLevel {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":          "f36b0cc8-a728-4a9b-a421-86f8e8dd137a",
				"channelSize": channelSize,
				"TesterGuiOwnerEngineChannelWarningLevel": common_config.TesterGuiOwnerEngineChannelWarningLevel,
				"TesterGuiOwnerEngineChannelSize":         common_config.TesterGuiOwnerEngineChannelSize,
			}).Warning("Number of messages on queue for 'TesterGuiOwnerEngineChannel' has reached a critical level")
		} else {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":          "189001d3-890a-4c3e-9396-d665daf11c3f",
				"channelSize": channelSize,
				"TesterGuiOwnerEngineChannelWarningLevel":                                       common_config.TesterGuiOwnerEngineChannelWarningLevel,
				"TesterGuiOwnerEngineChannelSize":                                               common_config.TesterGuiOwnerEngineChannelSize,
				"incomingTesterGuiOwnerEngineChannelCommand.TesterGuiOwnerEngineChannelCommand": incomingTesterGuiOwnerEngineChannelCommand.TesterGuiOwnerEngineChannelCommand,
			}).Info("Incoming TesterGuiOwnerEngineEngine-command")
		}

		switch incomingTesterGuiOwnerEngineChannelCommand.TesterGuiOwnerEngineChannelCommand {

		case common_config.ChannelCommand_ThisGuiExecutionServerIsClosingDown:
			processThisGuiExecutionServerIsClosingDown(
				incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_AnotherGuiExecutionServerIsClosingDown:
			processAnotherGuiExecutionServerIsClosingDown(
				incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_UserSubscribesToUserAndTestCaseExecutionCombination:
			processUserSubscribesToUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination:
			processUserUnsubscribesToUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_UserIsClosingDown:
			processUserIsClosingDown(incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_ThisGuiExecutionServerIsStartingUp:
			processThisGuiExecutionServerIsStartingUp(incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_AnotherGuiExecutionServerIsStartingUp:
			processAnotherGuiExecutionServerIsStartingUp(incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_ThisGuiExecutionServerSendsStartedUpTimeStamp:
			processThisGuiExecutionServerSendsStartedUpTimeStamp(incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_AnotherGuiExecutionServerSendsStartedUpTimeStamp:
			processAnotherGuiExecutionServerSendsStartedUpTimeStamp(incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination:
			processAnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination(incomingTesterGuiOwnerEngineChannelCommand)

		// No other command is supported
		default:
			common_config.Logger.WithFields(logrus.Fields{
				"Id": "8ef55340-bb8c-42cb-bfc-879b7407d64d",
				"incomingTesterGuiOwnerEngineChannelCommand": incomingTesterGuiOwnerEngineChannelCommand,
			}).Fatalln("Unhandled command in TesterGuiOwnerEngineChannel for TesterGuiOwnerEngine")
		}

		// Clear memory for Message
		incomingTesterGuiOwnerEngineChannelCommand = nil
	}
}

// Process channel command 'ChannelCommand_ThisGuiExecutionServerIsClosingDown'
func processThisGuiExecutionServerIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_ThisGuiExecutionServerIsClosingDown'
	commandThisGuiExecutionServerIsClosingDown(
		incomingTesterGuiOwnerEngineChannelCommand.GuiExecutionServerIsClosingDown)

	// Continue process to close down this server
	*incomingTesterGuiOwnerEngineChannelCommand.GuiExecutionServerIsClosingDown.
		CurrentGuiExecutionServerIsClosingDownReturnChannel <- true
}

// Process channel command 'ChannelCommand_AnotherGuiExecutionServerIsClosingDown'
func processAnotherGuiExecutionServerIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_AnotherGuiExecutionServerIsClosingDown'
	commandAnotherGuiExecutionServerIsClosingDown(
		incomingTesterGuiOwnerEngineChannelCommand.GuiExecutionServerIsClosingDown)
}

// Process channel command 'ChannelCommand_UserSubscribesToUserAndTestCaseExecutionCombination
func processUserSubscribesToUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_UserSubscribesToUserAndTestCaseExecutionCombination'
	commandUserSubscribesToUserAndTestCaseExecutionCombination(
		incomingTesterGuiOwnerEngineChannelCommand.UserSubscribesToUserAndTestCaseExecutionCombination)
}

// Process channel command 'ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination'
func processUserUnsubscribesToUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination'
	commandUserUnsubscribesToUserAndTestCaseExecutionCombination(
		incomingTesterGuiOwnerEngineChannelCommand.UserUnsubscribesToUserAndTestCaseExecutionCombination)
}

// Process channel command 'ChannelCommand_UserIsClosingDown'
func processUserIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_UserIsClosingDown'
	commandUserIsClosingDown(
		incomingTesterGuiOwnerEngineChannelCommand.TesterGuiIsClosingDown)
}

// Process channel command 'ChannelCommand_ThisGuiExecutionServerIsStartingUp'
func processThisGuiExecutionServerIsStartingUp(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_ThisGuiExecutionServerIsStartingUp'
	commandThisGuiExecutionServerIsStartingUp(
		incomingTesterGuiOwnerEngineChannelCommand.GuiExecutionServerIsStartingUp)
}

// Process channel command 'ChannelCommand_AnotherGuiExecutionServerIsStartingUp'
func processAnotherGuiExecutionServerIsStartingUp(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_AnotherGuiExecutionServerIsStartingUp'
	commandAnotherGuiExecutionServerIsStartingUp(
		incomingTesterGuiOwnerEngineChannelCommand.GuiExecutionServerIsStartingUp)
}

// Process channel command 'ChannelCommand_ThisGuiExecutionServerSendsStartedUpTimeStamp'
func processThisGuiExecutionServerSendsStartedUpTimeStamp(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_ThisGuiExecutionServerSendsStartedUpTimeStamp'
	commandThisGuiExecutionServerSendsStartedUpTimeStamp(
		incomingTesterGuiOwnerEngineChannelCommand.GuiExecutionServerStartedUpTimeStampRefresher)
}

// Process channel command 'ChannelCommand_AnotherGuiExecutionServerSendsStartedUpTimeStamp'
func processAnotherGuiExecutionServerSendsStartedUpTimeStamp(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_AnotherGuiExecutionServerSendsStartedUpTimeStamp'
	commandAnotherGuiExecutionServerSendsStartedUpTimeStamp(
		incomingTesterGuiOwnerEngineChannelCommand.GuiExecutionServerStartedUpTimeStampRefresher)
}

// Process channel command 'ChannelCommand_AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination'
func processAnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

	// Process the actual command 'ChannelCommand_AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination'
	commandAnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination(
		incomingTesterGuiOwnerEngineChannelCommand.AnotherGuiExecutionServerOvertakesThisTestCaseExecutionCombination)
}
