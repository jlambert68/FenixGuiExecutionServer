package messagesToExecutionServer

import (
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type MessagesToExecutionServerObjectStruct struct {
	Logger         *logrus.Logger
	gcpAccessToken *oauth2.Token
}

var MessagesToExecutionServerObject MessagesToExecutionServerObjectStruct

// Variables used for contacting Fenix Execution Worker Server
/*
var (
	remoteFenixExecutionWorkerServerConnection *grpc.ClientConn
	FenixExecutionServerAddressToDial          string
	fenixExecutionWorkerServerGrpcClient       fenixExecutionServerGrpcApi.FenixExecutionServerGrpcServicesClient
)

*/

// Used for keeping track of the proto file versions for ExecutionServer and this Worker
var highestFenixExecutionServerProtoFileVersion int32 = -1

//var highestExecutionWorkerProtoFileVersion int32 = -1

// Variables used for contacting Fenix ExecutionServer
var (
	RemoteFenixExecutionServerConnection *grpc.ClientConn
	FenixExecutionServerAddressToDial    string
	FenixExecutionServerGrpcClient       fenixExecutionServerGrpcApi.FenixExecutionServerGrpcServicesClient
	FenixExecutionServerAddressToUse     string
)
