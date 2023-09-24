package common_config

import "time"

// TestInstructionTesterGuiOwnerEngineEngineObjectStruct
// The struct for the object that hold all functions together within the TesterGuiOwnerEngineEngine
type TestInstructionTesterGuiOwnerEngineEngineObjectStruct struct {
}

// TestInstructionExecutionTesterGuiOwnerEngineEngineObject
// The object that hold all functions together within the TesterGuiOwnerEngineEngine
var TestInstructionExecutionTesterGuiOwnerEngineEngineObject TestInstructionTesterGuiOwnerEngineEngineObjectStruct

// TesterGuiOwnerEngineChannelEngineCommandChannelReferenceSlice
// A slice with references to  channels for the TestInstructionExecutionEngine
// Each position in slice represents one execution track
var TesterGuiOwnerEngineChannelEngineCommandChannel TesterGuiOwnerEngineChannelEngineType

// TesterGuiOwnerEngineChannelEngineType
// The channel type
type TesterGuiOwnerEngineChannelEngineType chan *TesterGuiOwnerEngineChannelCommandStruct

// TesterGuiOwnerEngineChannelCommandType
// The type for the constants used within the message sent in the TesterGuiOwnerEngineChannel
type TesterGuiOwnerEngineChannelCommandType uint8

const (
	ChannelCommand_ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination TesterGuiOwnerEngineChannelCommandType = iota
	ChannelCommand_AnotherGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination
	ChannelCommand_ThisGuiExecutionServerIsClosingDown
	ChannelCommand_AnotherGuiExecutionServerIsClosingDown
	ChannelCommand_UserSubscribesToUserAndTestCaseExecutionCombination
	ChannelCommand_UserUnsubscribesToUserAndTestCaseExecutionCombination
	ChannelCommand_UserIsClosingDown
)

// TesterGuiOwnerEngineChannelCommandStruct
// The struct for the message that are sent over the channel to the TesterGuiOwnerEngineEngine
type TesterGuiOwnerEngineChannelCommandStruct struct {
	TesterGuiOwnerEngineChannelCommand                                 TesterGuiOwnerEngineChannelCommandType
	SomeoneIsClosingDown                                               *SomeoneIsClosingDownStruct
	ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination *ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct
	UserUnsubscribesToUserAndTestCaseExecutionCombination              *UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
}

// SomeoneIsClosingDownStruct
// Used to specify an GuiExecutionServer or a TesterGui that is Closing Down
type SomeoneIsClosingDownStruct struct {
	WhoISClosingDown WhoISClosingDownType
	ApplicationId    string
	UserId           string
	MessageTimeStamp time.Time
}

// WhoISClosingDownType
// The type for the constants used within the message sent in the SomeoneIsClosingDownStruct
type WhoISClosingDownType uint8

// The specified of application that is closing down
const (
	GuiExecutionServer WhoISClosingDownType = iota
	TesterGui
)

// ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct
// Used to specify that a specified GuiExecutionServer takes over status-sending-control for a TesterGui for a specific TestCaseExecutionUuid
type ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct struct {
	TesterGuiApplicationId          string
	UserId                          string
	GuiExecutionServerApplicationId string
	TestCaseExecutionUuid           string
	TestCaseExecutionVersion        string
	MessageTimeStamp                time.Time
}

// UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
// Used to specify that a specified TesterGui unsubscribes to a  for a specific TestCaseExecutionUuid
type UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct struct {
	TesterGuiApplicationId          string
	UserId                          string
	GuiExecutionServerApplicationId string
	TestCaseExecutionUuid           string
	TestCaseExecutionVersion        string
	MessageTimeStamp                time.Time
}

// TesterGuiOwnerEngineChannelSize
// The size of the channel
const TesterGuiOwnerEngineChannelSize = 100

// TesterGuiOwnerEngineChannelWarningLevel
// The size of warning level for the channel
const TesterGuiOwnerEngineChannelWarningLevel = 90
