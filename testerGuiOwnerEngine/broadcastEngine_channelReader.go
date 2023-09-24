package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"github.com/sirupsen/logrus"
)

// Channel reader which is used for reading out commands to TesterGuiOwnerEngine
func (testInstructionExecutionTesterGuiOwnerEngineEngineObject *common_config.TestInstructionTesterGuiOwnerEngineEngineObjectStruct) startTesterGuiOwnerEngineChannelReader() {

	var incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct
	var channelSize int

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

		case common_config.ChannelCommand_ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination:
			processThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_AnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination:
			processAnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

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
			proceessUserUnsubscribesToUserAndTestCaseExecutionCombination(
				incomingTesterGuiOwnerEngineChannelCommand)

		case common_config.ChannelCommand_UserIsClosingDown:
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
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_AnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination'
func processAnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_ThisGuiExecutionServerIsClosingDown'
func processThisGuiExecutionServerIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_AnotherGuiExecutionServerIsClosingDown'
func processAnotherGuiExecutionServerIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_UserSubscribesToUserAndTestCaseExecutionCombination
func processUserSubscribesToUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination'
func proceessUserUnsubscribesToUserAndTestCaseExecutionCombination(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

}

// Process channel command 'ChannelCommand_UserIsClosingDown'
func processUserIsClosingDown(
	incomingTesterGuiOwnerEngineChannelCommand *common_config.TesterGuiOwnerEngineChannelCommandStruct) {

}
