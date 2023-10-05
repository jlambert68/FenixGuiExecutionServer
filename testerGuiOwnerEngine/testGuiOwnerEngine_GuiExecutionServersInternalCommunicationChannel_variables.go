package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"time"
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

// testCaseExecutionsSubscriptionsMap
// Map that holds information about all TestCaseExecutions that different TesterGui:s are subscribing to from this GuiExecutionServer
// map['TestCaseExecutionKey']*common_config.GuiExecutionServerResponsibilityStruct
var testCaseExecutionsSubscriptionsMap map[testCaseExecutionsSubscriptionsMapKeyType]*common_config.GuiExecutionServerResponsibilityStruct

// testCaseExecutionsSubscriptionsMapKeyType
// the Key to 'testCaseExecutionsSubscriptionsMap'. Is a concatenation of 'TestCaseExecutionUuid' and 'TestCaseExecutionUuidVersion'
type testCaseExecutionsSubscriptionsMapKeyType string

// guiExecutionServerStartUpOrderStruct
// Structure holding one applicationRunTimeUuid and StartUpTime for a GuiExecutionServer
type guiExecutionServerStartUpOrderStruct struct {
	applicationRunTimeUuid        string
	applicationRunTimeStartUpTime time.Time
}

// guiExecutionServerStartUpOrder
// Slice containing all GuiExecutionServers broadcasted starting order. GuiExecutionServers are stored in StartUpTimeOrder
// When the length == 1 then this GuiExecutionServer takes over all responsibility from other closing GuiExecutionServer
var guiExecutionServerStartUpOrder []*guiExecutionServerStartUpOrderStruct
