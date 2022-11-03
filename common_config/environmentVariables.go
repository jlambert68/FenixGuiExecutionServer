package common_config

// ***********************************************************************************************************
// The following variables receives their values from environment variables

// ExecutionLocationForFenixExecutionServer
// Where is the Worker running
var ExecutionLocationForFenixExecutionServer ExecutionLocationTypeType

// ExecutionLocationForFenixGuiExecutionServer
// Where is Fenix Execution Server running
var ExecutionLocationForFenixGuiExecutionServer ExecutionLocationTypeType

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

// Address to Fenix TestData Server & Client, will have their values from Environment variables at startup
var (
	FenixTestDataSyncServerAddress string // TODO remove, but is referenced by code that is not removed yet
	FenixTestDataSyncServerPort    int    // TODO remove,
	FenixGuiServerAddress          string
	FenixExecutionServerPort       int
)

// FenixExecutionWorkerServerPort
// Execution Worker Port to use, will have its value from Environment variables at startup
var FenixExecutionWorkerServerPort int

// Address to use when not run locally, not on GCP/Cloud which gets its address from DB
var FenixExecutionWorkerAddress string

// FenixExecutionExecutionServerPort
// Execution Server Port to use, will have its value from Environment variables at startup
var FenixExecutionExecutionServerPort int
