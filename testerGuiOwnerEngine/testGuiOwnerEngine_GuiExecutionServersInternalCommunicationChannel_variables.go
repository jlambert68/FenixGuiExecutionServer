package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
)

// GuiExecutionServersInternalCommunicationChannelTypeType
// The type for the constants used within the message sent in the 'BroadcastMessageForSomeoneIsClosingDownMessage'
type GuiExecutionServersInternalCommunicationChannelTypeType uint8

// The specified of application that is closing down
const (
	TesterGuiIsClosingDownMessage GuiExecutionServersInternalCommunicationChannelTypeType = iota
	GuiExecutionServerIsClosingDownMessage
	ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationMessage
	UserUnsubscribesToUserAndTestCaseExecutionCombinationMessage
	GuiExecutionServerIsStartingUpMessage
)

// BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct
// The following message is sent over Postgres Broadcast system, 'Channel 1'
// Used to specify that a TesterGui or GuiExecutionServer is Closing Down
type BroadcastMessageForGuiExecutionServersInternalCommunicationChannelStruct struct {
	GuiExecutionServersInternalCommunicationChannelType                GuiExecutionServersInternalCommunicationChannelTypeType                                `json:"guiexecutionserversinternalcommunicationchanneltype"`
	TesterGuiIsClosingDown                                             common_config.TesterGuiIsClosingDownStruct                                             `json:"testerguiisclosingdown"`
	GuiExecutionServerIsClosingDown                                    common_config.GuiExecutionServerIsClosingDownStruct                                    `json:"guiexecutionserverisclosingdown"`
	ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination common_config.ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct `json:"thisguiexecutionservertakesthisuserandtestcaseexecutioncombination"`
	UserUnsubscribesToUserAndTestCaseExecutionCombination              common_config.UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct              `json:"userunsubscribestouserandtestcaseexecutioncombination"`
	GuiExecutionServerIsStartingUp                                     common_config.GuiExecutionServerIsStartingUpStruct                                     `json:"guiexecutionserverisstartingup"`
}

/*
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
	GuiExecutionServerApplicationId    string                                   `json:"guiexecutionserverapplicationid"`
	MessageTimeStamp                   time.Time                                `json:"messagetimestamp"`
	GuiExecutionServerResponsibilities []GuiExecutionServerResponsibilityStruct `json:"guiexecutionserverresponsibilities"`
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


*/
