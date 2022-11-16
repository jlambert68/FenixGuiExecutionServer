package main

import (
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type fenixGuiExecutionServerObjectStruct struct {
	logger *logrus.Logger
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

// Used  by gRPC server that receives Connector-connections to inform gRPC-server that receives ExecutionServer-connections
var TesterGuiHasConnected bool
