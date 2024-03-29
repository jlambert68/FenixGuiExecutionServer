package broadcastEngine_ExecutionStatusUpdate

import (
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"time"
)

// *******************************************************************************************
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
// map['applicationRunTimeUuid']*TestCaseExecutionsSubscriptionChannelInformationStruct
var TestCaseExecutionsSubscriptionChannelInformationMap map[ApplicationRunTimeUuidType]*TestCaseExecutionsSubscriptionChannelInformationStruct

// TestCaseExecutionsSubscriptionsMap
// Map that holds information about who is subscribing to a certain TestCaseExecution
// map['TestCaseExecutionUuid']*TestCaseExecutionsSubscriptionsStruct
var TestCaseExecutionsSubscriptionsMap map[TestCaseExecutionsSubscriptionsMapKeyType]*[]ApplicationRunTimeUuidType

// TestCaseExecutionsSubscriptionChannelInformationStruct
// Holds all information needed to be able to send ExecutionStatus-messages back to correct TesterGui
type TestCaseExecutionsSubscriptionChannelInformationStruct struct {
	ApplicationRunTimeUuid           ApplicationRunTimeUuidType
	LastConnectionFromTesterGui      time.Time
	MessageToTesterGuiForwardChannel *MessageToTesterGuiForwardChannelType
}

// ApplicationRunTimeUuidType
// Type used for 'applicationRunTimeUuid'
type ApplicationRunTimeUuidType string

// TestCaseExecutionUuidType
// Type used for 'TestCaseExecutionUuid'
type TestCaseExecutionUuidType string

// TestCaseExecutionsSubscriptionsMapKeyType
// the Key to 'TestCaseExecutionsSubscriptionsMap'. Is a concatenation of 'TestCaseExecutionUuid' and 'TestCaseExecutionUuidVersion'
type TestCaseExecutionsSubscriptionsMapKeyType string

// TestCaseExecutionUuidVersionType
// Type used for 'TestCaseExecutionUuidVersion'
type TestCaseExecutionUuidVersionType int
