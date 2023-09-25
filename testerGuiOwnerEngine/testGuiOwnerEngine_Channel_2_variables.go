package testerGuiOwnerEngine

import "time"

// PostgresChannel2MessageMessageTypeType
// The type for the constants used within the message sent in the 'PostgresChannel2Message'
type PostgresChannel2MessageMessageTypeType uint8

// The specified of application that is closing down
const (
	ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationMessage PostgresChannel2MessageMessageTypeType = iota
	UserUnsubscribesToUserAndTestCaseExecutionCombinationMessage
)

// // The following message is sent over Postgres Broadcast system, TesterGuiOwner-Channel 2
type BroadcastMesageForPostgresChannel2MessageStruct struct {
	PostgresChannel2MessageMessageType                                 PostgresChannel2MessageMessageTypeType                                   `json:"postgreschannel2messagemessagetype"`
	ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombination ThisGuiExecutionServerTakesThisUserAndTestCaseExecutionCombinationStruct `json:"thisguiexecutionservertakesthisuserandtestcaseexecutioncombination"`
	UserUnsubscribesToUserAndTestCaseExecutionCombination              UserUnsubscribesToUserAndTestCaseExecutionCombinationStruct              `json:"userunsubscribestouserandtestcaseexecutioncombination"`
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
