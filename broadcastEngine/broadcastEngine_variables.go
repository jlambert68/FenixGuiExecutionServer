package broadcastEngine

import fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"

// *******************************************************************************************
// Channel used for forwarding MessagesToTestGui to stream-server which then forwards it to the TesterGui
var MessageToTesterGuiForwardChannel MessageToTesterGuiForwardChannelType

type MessageToTesterGuiForwardChannelType chan MessageToTestGuiForwardChannelStruct

const MessageToTesterGuiForwardChannelMaxSize int32 = 100

type MessageToTestGuiForwardChannelStruct struct {
	SubscribeToMessagesStreamResponse *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
	IsKeepAliveMessage                bool
}
