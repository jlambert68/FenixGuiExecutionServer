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
	fenixExecutionServerGuiGrpcApi.UnimplementedFenixExecutionServerGuiGrpcServicesServer
}

//TODO FIXA DENNA PATH, HMMM borde köra i DB framöver
// For now hardcoded MerklePath
//var merkleFilterPath string = //"AccountEnvironment/ClientJuristictionCountryCode/MarketSubType/MarketName/" //SecurityType/"

var highestFenixGuiExecutionServerProtoFileVersion int32 = -1
