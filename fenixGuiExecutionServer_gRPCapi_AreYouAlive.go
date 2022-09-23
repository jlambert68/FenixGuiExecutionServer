package main

import (
	fenixExecutionServerGuiGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGuiGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// AreYouAlive - *********************************************************************
//Anyone can check if Fenix TestCase Builder server is alive with this service
func (s *fenixExecutionServerGrpcServicesServer) AreYouAlive(ctx context.Context, emptyParameter *fenixExecutionServerGuiGrpcApi.EmptyParameter) (*fenixExecutionServerGuiGrpcApi.AckNackResponse, error) {

	fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "1ff67695-9a8b-4821-811d-0ab8d33c4d8b",
	}).Debug("Incoming 'gRPC - AreYouAlive'")

	defer fenixGuiExecutionServerObject.logger.WithFields(logrus.Fields{
		"id": "9c7f0c3d-7e9f-4c91-934e-8d7a22926d84",
	}).Debug("Outgoing 'gRPC - AreYouAlive'")

	return &fenixExecutionServerGuiGrpcApi.AckNackResponse{AckNack: true, Comments: "I'am alive."}, nil
}
