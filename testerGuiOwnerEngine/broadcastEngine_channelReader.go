package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"github.com/sirupsen/logrus"
)

// Channel reader which is used for reading out commands to TesterGuiOwnerEngine
func (testInstructionExecutionTesterGuiOwnerEngineEngineObject *TestInstructionTesterGuiOwnerEngineEngineObjectStruct) startTesterGuiOwnerEngineChannelReader() {

	var incomingTesterGuiOwnerEngineChannelCommand *TesterGuiOwnerEngineChannelCommandStruct
	var channelSize int

	for {
		// Wait for incoming command over channel
		incomingTesterGuiOwnerEngineChannelCommand = <-TesterGuiOwnerEngineChannelEngineCommandChannel

		common_config.Logger.WithFields(logrus.Fields{
			"Id": "a2809c91-87bc-44fc-894b-c8cdd73b521f",
			"incomingTesterGuiOwnerEngineChannelCommand": incomingTesterGuiOwnerEngineChannelCommand,
		}).Debug("Message received on 'TesterGuiOwnerEngineChannel'")

		// If size of Channel > 'TesterGuiOwnerEngineChannelWarningLevel' then log Warning message
		channelSize = len(TesterGuiOwnerEngineChannelEngineCommandChannel)
		if channelSize > TesterGuiOwnerEngineChannelWarningLevel {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":          "f36b0cc8-a728-4a9b-a421-86f8e8dd137a",
				"channelSize": channelSize,
				"TesterGuiOwnerEngineChannelWarningLevel": TesterGuiOwnerEngineChannelWarningLevel,
				"TesterGuiOwnerEngineChannelSize":         TesterGuiOwnerEngineChannelSize,
			}).Warning("Number of messages on queue for 'TesterGuiOwnerEngineChannel' has reached a critical level")
		} else {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":          "189001d3-890a-4c3e-9396-d665daf11c3f",
				"channelSize": channelSize,
				"TesterGuiOwnerEngineChannelWarningLevel":                                       TesterGuiOwnerEngineChannelWarningLevel,
				"TesterGuiOwnerEngineChannelSize":                                               TesterGuiOwnerEngineChannelSize,
				"incomingTesterGuiOwnerEngineChannelCommand.TesterGuiOwnerEngineChannelCommand": incomingTesterGuiOwnerEngineChannelCommand.TesterGuiOwnerEngineChannelCommand,
			}).Info("Incoming TesterGuiOwnerEngineEngine-command")
		}

		switch incomingTesterGuiOwnerEngineChannelCommand.TesterGuiOwnerEngineChannelCommand {

		case ChannelCommand_ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination:
			processThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

		case ChannelCommand_AnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination:
			processAnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

		case ChannelCommand_ThisGuiExecutionServerIsClosingDown:
			processThisGuiExecutionServerIsClosingDown(
				incomingTesterGuiOwnerEngineChannelCommand)

		case ChannelCommand_AnotherGuiExecutionServerIsClosingDown:
			processAnotherGuiExecutionServerIsClosingDown(
				incomingTesterGuiOwnerEngineChannelCommand)

		case ChannelCommand_UserSubscribesToUserAndTestCaseExecutionCombination:
			processUserSubscribesToUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

		case ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination:
			proceessUserUnsubscribesToUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

		case ChannelCommand_UserIsClosingDown:
			processUserIsClosingDown(incomingTesterGuiOwnerEngineChannelCommand)

		// No other command is supported
		default:
			common_config.Logger.WithFields(logrus.Fields{
				"Id": "8ef55340-bb8c-42cb-bfc-879b7407d64d",
				"incomingTesterGuiOwnerEngineChannelCommand": incomingTesterGuiOwnerEngineChannelCommand,
			}).Fatalln("Unknown command in TesterGuiOwnerEngineChannel for TesterGuiOwnerEngine")
		}
	}
}

// Process channel command 'ChannelCommand_ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination'
func processThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_AnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination'
func processAnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_ThisGuiExecutionServerIsClosingDown'
func processThisGuiExecutionServerIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_AnotherGuiExecutionServerIsClosingDown'
func processAnotherGuiExecutionServerIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_UserSubscribesToUserAndTestCaseExecutionCombination
func processUserSubscribesToUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination'
func proceessUserUnsubscribesToUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_UserIsClosingDown'
func processUserIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *TesterGuiOwnerEngineChannelCommandStruct) {

}
