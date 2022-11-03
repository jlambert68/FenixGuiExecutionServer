package common_config

// ***********************************************************************************************************
// The following variables receives their values from environment variables

// ExecutionLocationTypeType
// Definitions for where GuiExecutionServer and ExecutionServer is running
type ExecutionLocationTypeType int

// Constants used for where stuff is running
const (
	LocalhostNoDocker ExecutionLocationTypeType = iota
	LocalhostDocker
	GCP
)

// FenixGuiBuilderServer
var LocationForFenixGuiBuilderServerTypeMapping = map[ExecutionLocationTypeType]string{
	LocalhostNoDocker: "LOCALHOST_NODOCKER",
	LocalhostDocker:   "LOCALHOST_DOCKER",
	GCP:               "GCP",
}

// Environmentvaribales extracted when program starts up
var ExecutionLocationForFenixGuiExecutionServer ExecutionLocationTypeType
var FenixGuiExecutionServerAddress string
var FenixGuiExecutionServerPort int

var ExecutionLocationForFenixExecutionServer ExecutionLocationTypeType
var FenixExecutionServerAddress string
var FenixExecutionServerPort int
