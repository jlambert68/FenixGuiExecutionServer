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
	ChannelCommand_ThisGuiExecutionServerIsClosingDown TesterGuiOwnerEngineChannelCommandType = iota
	ChannelCommand_AnotherGuiExecutionServerIsClosingDown
	ChannelCommand_ThisGuiExecutionServersUserSubscribesToUserAndTestCaseExecutionCombination
	ChannelCommand_ThisGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination
	ChannelCommand_ThisGuiExecutionServersTesterGuiIsClosingDown
	ChannelCommand_ThisGuiExecutionServerIsStartingUp
	ChannelCommand_AnotherGuiExecutionServerIsStartingUp
	ChannelCommand_ThisGuiExecutionServerSendsStartedUpTimeStamp
	ChannelCommand_AnotherGuiExecutionServerSendsStartedUpTimeStamp
	ChannelCommand_AnotherGuiExecutionServersTesterGuiIsClosingDown
	ChannelCommand_AnotherGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination
	ChannelCommand_AnotherGuiExecutionServersUserSubscribesToUserAndTestCaseExecutionCombination
)

var ChannelCommand_Descriptions = map[TesterGuiOwnerEngineChannelCommandType]string{
	0:  "ChannelCommand_ThisGuiExecutionServerIsClosingDown",
	1:  "ChannelCommand_AnotherGuiExecutionServerIsClosingDown",
	2:  "ChannelCommand_ThisGuiExecutionServersUserSubscribesToUserAndTestCaseExecutionCombination",
	3:  "ChannelCommand_ThisGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination",
	4:  "ChannelCommand_ThisGuiExecutionServersTesterGuiIsClosingDown",
	5:  "ChannelCommand_ThisGuiExecutionServerIsStartingUp",
	6:  "ChannelCommand_AnotherGuiExecutionServerIsStartingUp",
	7:  "ChannelCommand_ThisGuiExecutionServerSendsStartedUpTimeStamp",
	8:  "ChannelCommand_AnotherGuiExecutionServerSendsStartedUpTimeStamp",
	10: "ChannelCommand_AnotherGuiExecutionServersTesterGuiIsClosingDown",
	11: "ChannelCommand_AnotherGuiExecutionServersUserUnsubscribesToUserAndTestCaseExecutionCombination",
	12: "ChannelCommand_AnotherGuiExecutionServersUserSubscribesToUserAndTestCaseExecutionCombination",
}

// TesterGuiOwnerEngineChannelCommandStruct
// The struct for the message that are sent over the channel to the TesterGuiOwnerEngineEngine
type TesterGuiOwnerEngineChannelCommandStruct struct {
	TesterGuiOwnerEngineChannelCommand                    TesterGuiOwnerEngineChannelCommandType
	TesterGuiIsClosingDown                                *TesterGuiIsClosingDownStruct
	GuiExecutionServerIsClosingDown                       *GuiExecutionServerIsClosingDownStruct
	UserUnsubscribesToUserAndTestCaseExecutionCombination *UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
	GuiExecutionServerIsStartingUp                        *GuiExecutionServerIsStartingUpStruct
	GuiExecutionServerStartedUpTimeStampRefresher         *GuiExecutionServerStartedUpTimeStampRefresherStruct
	UserSubscribesToUserAndTestCaseExecutionCombination   *UserSubscribesToUserAndTestCaseExecutionCombinationStruct
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

// GuiExecutionServerIsStartingUpStruct
// The following message is sent over Postgres Broadcast system and over TesterGuiOwnerEngine-channel
// Used to specify that a GuiExecutionServer is Starting Up
type GuiExecutionServerIsStartingUpStruct struct {
	GuiExecutionServerApplicationId string    `json:"guiexecutionserverapplicationid"`
	MessageTimeStamp                time.Time `json:"messagetimestamp"`
}

// GuiExecutionServerStartedUpTimeStampRefresherStruct
// The following message is sent over Postgres Broadcast system and over TesterGuiOwnerEngine-channel
// Message is Broadcasted by GuiExecutionServer to other GuiExecutionServers to refresh and sync StartUp-TimeStampss
type GuiExecutionServerStartedUpTimeStampRefresherStruct struct {
	GuiExecutionServerApplicationId string    `json:"guiexecutionserverapplicationid"`
	MessageTimeStamp                time.Time `json:"messagetimestamp"`
}

// GuiExecutionServerResponsibilityStruct
// Holds one Responsibility for GuiExecutionServer, regarding sending ExecutionStatusUpdates
type GuiExecutionServerResponsibilityStruct struct {
	TesterGuiApplicationId   string `json:"testguiapplicationid"`
	UserId                   string `json:"userid"`
	TestCaseExecutionUuid    string `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion int32  `json:"testcaseexecutionversion"`
}

// UserSubscribesToUserAndTestCaseExecutionCombinationStruct
// Used to specify that a specified TesterGui subscribes to a specific TestCaseExecutionUuid
type UserSubscribesToUserAndTestCaseExecutionCombinationStruct struct {
	TesterGuiApplicationId          string    `json:"testerguiapplicationid"`
	UserId                          string    `json:"userid"`
	GuiExecutionServerApplicationId string    `json:"guiexecutionserverapplicationid"`
	TestCaseExecutionUuid           string    `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion        int32     `json:"testcaseexecutionversion"`
	MessageTimeStamp                time.Time `json:"messagetimestamp"`
}

// UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct
// Used to specify that a specified TesterGui unsubscribes to a specific TestCaseExecutionUuid
type UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct struct {
	TesterGuiApplicationId          string    `json:"testerguiapplicationid"`
	UserId                          string    `json:"userid"`
	GuiExecutionServerApplicationId string    `json:"guiexecutionserverapplicationid"`
	TestCaseExecutionUuid           string    `json:"testcaseexecutionuuid"`
	TestCaseExecutionVersion        int32     `json:"testcaseexecutionversion"`
	MessageTimeStamp                time.Time `json:"messagetimestamp"`
}

// TesterGuiOwnerEngineChannelSize
// The size of the channel
const TesterGuiOwnerEngineChannelSize = 100

// TesterGuiOwnerEngineChannelWarningLevel
// The size of warning level for the channel
const TesterGuiOwnerEngineChannelWarningLevel = 90
