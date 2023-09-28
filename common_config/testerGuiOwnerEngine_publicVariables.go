package common_config

import "time"

/*
// TestInstructionTesterGuiOwnerEngineEngineObjectStruct
// The struct for the object that hold all functions together within the TesterGuiOwnerEngineEngine
type TestInstructionTesterGuiOwnerEngineEngineObjectStruct struct {
}

// TestInstructionExecutionTesterGuiOwnerEngineEngineObject
// The object that hold all functions together within the TesterGuiOwnerEngineEngine
var TestInstructionExecutionTesterGuiOwnerEngineEngineObject TestInstructionTesterGuiOwnerEngineEngineObjectStruct
*/

// TesterGuiOwnerEngineChannelEngineCommandChannel
// The channels for the TestInstructionExecutionEngine
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
	TesterGuiIsClosingDown                                             *TesterGuiIsClosingDownStruct
	GuiExecutionServerIsClosingDown                                    *GuiExecutionServerIsClosingDownStruct
	ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination *ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct
	UserUnsubscribesToUserAndTestCaseExecutionCombination              *UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
}

// TesterGuiIsClosingDownStruct
// The following message is sent over Postgres Broadcast system and over TesterGuiOwnerEngine-channel
// Used to specify that a TesterGui is Closing Down
type TesterGuiIsClosingDownStruct struct {
	TesterGuiApplicationId          string    `json:"testerguiapplicationid"`
	UserId                          string    `json:"userid"`
	GuiExecutionServerApplicationId string    `json:"guiexecutionserverapplicationid"`
	MessageTimeStamp                time.Time `json:"messagetimestamp"`
}

// GuiExecutionServerIsClosingDownStruct
// The following message is sent over Postgres Broadcast system and over TesterGuiOwnerEngine-channel
// Used to specify that a GuiExecutionServer is Closing Down
type GuiExecutionServerIsClosingDownStruct struct {
	GuiExecutionServerApplicationId                     string                                   `json:"guiexecutionserverapplicationid"`
	MessageTimeStamp                                    time.Time                                `json:"messagetimestamp"`
	GuiExecutionServerResponsibilities                  []GuiExecutionServerResponsibilityStruct `json:"guiexecutionserverresponsibilities"`
	CurrentGuiExecutionServerIsClosingDownReturnChannel *chan bool                               // Should not be converted into json
}

// GuiExecutionServerResponsibilityStruct
// Holds one Responsibility for GuiExecutionServer, regarding sending ExecutionStatusUpdates
type GuiExecutionServerResponsibilityStruct struct {
	TesterGuiApplicationId   string `json:"testguiapplicationid"`
	UserId                   string `json:"userid"`
	TestCaseExecutionUuid    string `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion int    `json:"testcaseexecutionversion"`
}

// ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct
// Used to specify that a specified GuiExecutionServer takes over status-sending-control for a TesterGui for a specific TestCaseExecutionUuid
type ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct struct {
	TesterGuiApplicationId          string    `json:"testerguiapplicationid"`
	UserId                          string    `json:"userid"`
	GuiExecutionServerApplicationId string    `json:"guiexecutionserverapplicationid"`
	TestCaseExecutionUuid           string    `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion        string    `json:"testcaseexecutionversion"`
	MessageTimeStamp                time.Time `json:"messagetimestamp"`
}

// UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
// Used to specify that a specified TesterGui unsubscribes to a  for a specific TestCaseExecutionUuid
type UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct struct {
	TesterGuiApplicationId          string    `json:"testerguiapplicationid"`
	UserId                          string    `json:"userid"`
	GuiExecutionServerApplicationId string    `json:"guiexecutionserverapplicationid"`
	TestCaseExecutionUuid           string    `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion        string    `json:"testcaseexecutionversion"`
	MessageTimeStamp                time.Time `json:"messagetimestamp"`
}

// TesterGuiOwnerEngineChannelSize
// The size of the channel
const TesterGuiOwnerEngineChannelSize = 100

// TesterGuiOwnerEngineChannelWarningLevel
// The size of warning level for the channel
const TesterGuiOwnerEngineChannelWarningLevel = 90
