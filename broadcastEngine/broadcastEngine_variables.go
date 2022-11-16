package broadcastEngine

import (
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"time"
)

// *******************************************************************************************
// Channel used for forwarding MessagesToTestGui to stream-server which then forwards it to the TesterGui
var MessageToTesterGuiForwardChannel MessageToTesterGuiForwardChannelType

// Channel type for 'MessageToTesterGuiForwardChannel'
type MessageToTesterGuiForwardChannelType chan MessageToTestGuiForwardChannelStruct

// MessageToTesterGuiForwardChannelMaxSize
// Maximum size of 'MessageToTesterGuiForwardChannel'
const MessageToTesterGuiForwardChannelMaxSize int32 = 100

// MessageToTestGuiForwardChannelStruct
// Messages on MessageChannel has the follwoing information
type MessageToTestGuiForwardChannelStruct struct {
	SubscribeToMessagesStreamResponse *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
	IsKeepAliveMessage                bool
}

// TestCaseExecutionsSubscriptionChannelInformationMap
// Map for holding information about the data needed to route ExecutionStatuses to correct TesterGui
// map['ApplicationRunTimeUuid']*TestCaseExecutionsSubscriptionChannelInformationStruct
var TestCaseExecutionsSubscriptionChannelInformationMap map[ApplicationRunTimeUuidType]*TestCaseExecutionsSubscriptionChannelInformationStruct

// TestCaseExecutionsSubscriptionsMap
// Map that holds information about who is subscribing to a certain TestCaseExecution
// map['TestCaseExecutionUuid']*TestCaseExecutionsSubscriptionsStruct
var TestCaseExecutionsSubscriptionsMap map[TestCaseExecutionUuidType]*[]ApplicationRunTimeUuidType

// TestCaseExecutionsSubscriptionChannelInformationStruct
// Holds all information needed to be able to send ExecutionStatus-messages back to correct TesterGui
type TestCaseExecutionsSubscriptionChannelInformationStruct struct {
	ApplicationRunTimeUuid           ApplicationRunTimeUuidType
	LastConnectionFromTesterGui      time.Time
	MessageToTesterGuiForwardChannel MessageToTesterGuiForwardChannelType
}

// ApplicationRunTimeUuidType
// Type used for 'ApplicationRunTimeUuid'
type ApplicationRunTimeUuidType string

// TestCaseExecutionUuidType
// Type used for 'TestCaseExecutionUuid'
type TestCaseExecutionUuidType string
