package main

import (
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/messagesToExecutionServer"
	"github.com/sirupsen/logrus"
	"strconv"

	//"flag"
	"fmt"
	"log"
	"os"
)

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func main() {
	//time.Sleep(15 * time.Second)
	fenixGuiExecutionServerMain()
}

func init() {
	//executionLocationForGuiExecutionServer := flag.String("startupType", "0", "The application should be started with one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
	//flag.Parse()

	var err error

	// Where is GuiExecutionServer started
	var executionLocationForGuiExecutionServer = mustGetenv("ExecutionLocationForFenixGuiExecutionServer")

	switch executionLocationForGuiExecutionServer {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForFenixGuiExecutionServer = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForFenixGuiExecutionServer = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForFenixGuiExecutionServer = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for FenixGuiExecutionServer: " + executionLocationForGuiExecutionServer + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Address to GuiExecutionServer
	common_config.FenixGuiExecutionServerAddress = mustGetenv("FenixGuiExecutionServerAddress")

	// Port for GuiExecutionServer
	common_config.FenixGuiExecutionServerPort, err = strconv.Atoi(mustGetenv("FenixGuiExecutionServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixGuiExecutionServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Where is ExecutionServer started
	var executionLocationForExecutionServer = mustGetenv("ExecutionLocationForFenixExecutionServer")

	switch executionLocationForExecutionServer {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for FenixExecutionServer: " + executionLocationForExecutionServer + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Address to ExecutionServer
	common_config.FenixExecutionServerAddress = mustGetenv("FenixExecutionServerAddress")

	// Port for ExecutionServer
	common_config.FenixExecutionServerPort, err = strconv.Atoi(mustGetenv("FenixExecutionServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixExecutionServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Save the Dial-string to use for connecting to ExecutionServer
	messagesToExecutionServer.FenixExecutionServerAddressToDial = common_config.FenixExecutionServerAddress + ":" + strconv.Itoa(common_config.FenixExecutionServerPort)

	// Save the address to use for getting access token
	messagesToExecutionServer.FenixExecutionServerAddressToUse = common_config.FenixExecutionServerAddress

	// Extract Debug level
	var loggingLevel = mustGetenv("LoggingLevel")

	switch loggingLevel {

	case "DebugLevel":
		common_config.LoggingLevel = logrus.DebugLevel

	case "InfoLevel":
		common_config.LoggingLevel = logrus.InfoLevel

	default:
		fmt.Println("Unknown loggingLevel '" + loggingLevel + "'. Expected one of the following: 'DebugLevel', 'InfoLevel'")
		os.Exit(0)

	}

	fmt.Printf("%s", common_config.LoggingLevel)

}
