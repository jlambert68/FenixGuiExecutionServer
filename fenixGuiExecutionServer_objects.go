package main

import (
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"net"
)

type fenixGuiExecutionServerObjectStruct struct {
	logger         *logrus.Logger
	gcpAccessToken *oauth2.Token
}

// Variable holding everything together
var fenixGuiExecutionServerObject *fenixGuiExecutionServerObjectStruct

// gRPC variables
var (
	registerFenixGuiExecutionServerGrpcServicesServer *grpc.Server // registerFenixGuiExecutionServerGrpcServicesServer *grpc.Server
	lis                                               net.Listener
)

// gRPC Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type fenixGuiExecutionServerGrpcServicesServer struct {
	fenixExecutionServerGuiGrpcApi.UnimplementedFenixExecutionServerGuiGrpcServicesForGuiClientServer
}

//TODO FIXA DENNA PATH, HMMM borde köra i DB framöver
// For now hardcoded MerklePath
//var merkleFilterPath string = //"AccountEnvironment/ClientJuristictionCountryCode/MarketSubType/MarketName/" //SecurityType/"

var highestFenixGuiExecutionServerProtoFileVersion int32 = -1

// Used  by gRPC server that receives Connector-connections to inform gRPC-server that receives ExecutionServer-connections
var TesterGuiHasConnected bool

// *******************************************************************************************
// Channel used for forwarding TestInstructionExecutions to stream-server which then forwards it to the Connector
var messageToTesterGuiForwardChannel messageToTesterGuiForwardChannelType

type messageToTesterGuiForwardChannelType chan messageToTestGuiForwardChannelStruct

const messageToTesterGuiForwardChannelMaxSize int32 = 100

type messageToTestGuiForwardChannelStruct struct {
	subscribeToMessagesStreamResponse *fenixExecutionServerGuiGrpcApi.SubscribeToMessagesStreamResponse
	isKeepAliveMessage                bool
}
